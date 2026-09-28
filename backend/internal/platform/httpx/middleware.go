package httpx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
)

// Header names used by the API (API_SPEC.md §2, §3).
const (
	RequestIDHeader  = "X-Request-ID"
	AdminTokenHeader = "X-Admin-Token"
)

type ctxKey int

const (
	requestIDKey ctxKey = iota
	requestStateKey
)

// requestState collects what handlers report for the request log line.
type requestState struct {
	err error
}

// recordError stores err for the request log line written by LogRequests.
func recordError(r *http.Request, err error) {
	if st, ok := r.Context().Value(requestStateKey).(*requestState); ok {
		st.err = err
	}
}

// Chain wraps h with middleware; the first middleware is the outermost.
func Chain(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// RequestID gives every request an ID and echoes it in the X-Request-ID
// response header. A well-formed incoming ID (e.g. from a proxy) is reused.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if !requestIDPattern.MatchString(id) {
			id = newRequestID()
		}
		w.Header().Set(RequestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))
	})
}

func newRequestID() string {
	var b [8]byte
	_, _ = rand.Read(b[:]) // a failure leaves zeros, which only weakens log correlation
	return hex.EncodeToString(b[:])
}

// RequestIDFrom returns the request ID set by RequestID, or "".
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

// LogRequests logs one line per request. 5xx responses are logged at Error
// with the cause recorded by WriteError; /healthz is logged at Debug to keep
// orchestrator probes out of the logs. Headers are never logged.
func LogRequests(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			st := &requestState{}
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r.WithContext(context.WithValue(r.Context(), requestStateKey, st)))

			attrs := []any{
				"request_id", RequestIDFrom(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
			}
			switch {
			case rec.status >= http.StatusInternalServerError:
				log.Error("request failed", append(attrs, "err", st.err)...)
			case r.URL.Path == HealthzPath:
				log.Debug("request", attrs...)
			default:
				log.Info("request", attrs...)
			}
		})
	}
}

// statusRecorder remembers the status code written by the handler.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status = code
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	r.wroteHeader = true
	return r.ResponseWriter.Write(b)
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

// Recover turns a panic into 500 INTERNAL_ERROR so one bad request cannot
// crash the server. It is a last-resort net, not control flow. Place it
// inside LogRequests so the panic is logged with the request.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}
			if err, ok := v.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				panic(v) // net/http aborts the connection without logging
			}
			err := fmt.Errorf("panic: %v\n%s", v, debug.Stack())
			if rec, ok := w.(*statusRecorder); ok && rec.wroteHeader {
				recordError(r, err) // too late for an error body
				return
			}
			WriteError(w, r, err)
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS lets browsers on the allowed origins make GET requests. Only GET is
// allowed, so browsers can never call the admin endpoints. With no allowed
// origins it adds nothing and cross-origin browser requests fail.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	isAllowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		isAllowed[strings.ToLower(o)] = true
	}
	return func(next http.Handler) http.Handler {
		if len(isAllowed) == 0 {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Add("Vary", "Origin") // responses are publicly cacheable
			origin := r.Header.Get("Origin")
			if origin == "" || !isAllowed[strings.ToLower(origin)] {
				next.ServeHTTP(w, r)
				return
			}
			h.Set("Access-Control-Allow-Origin", origin)
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				h.Set("Access-Control-Allow-Methods", http.MethodGet)
				h.Set("Access-Control-Allow-Headers", "Accept, Content-Type")
				h.Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h.Set("Access-Control-Expose-Headers", RequestIDHeader+", Retry-After")
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdminToken rejects requests whose X-Admin-Token header does not
// match token with 401 UNAUTHORIZED. Both values are hashed first so the
// constant-time comparison does not leak the token length. An empty token
// rejects everything.
func RequireAdminToken(token string) func(http.Handler) http.Handler {
	want := sha256.Sum256([]byte(token))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := sha256.Sum256([]byte(r.Header.Get(AdminTokenHeader)))
			if token == "" || subtle.ConstantTimeCompare(got[:], want[:]) != 1 {
				WriteError(w, r, ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
