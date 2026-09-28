package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestRequireAdminToken(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		header     *string
		wantStatus int
	}{
		{"correct token", "s3cret-token", ptr("s3cret-token"), http.StatusAccepted},
		{"missing header", "s3cret-token", nil, http.StatusUnauthorized},
		{"empty header", "s3cret-token", ptr(""), http.StatusUnauthorized},
		{"wrong token", "s3cret-token", ptr("guess"), http.StatusUnauthorized},
		{"prefix of the token", "s3cret-token", ptr("s3cret"), http.StatusUnauthorized},
		{"empty configured token rejects everything", "", ptr(""), http.StatusUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				isCalled = true
				w.WriteHeader(http.StatusAccepted)
			})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/refresh", nil)
			if tt.header != nil {
				req.Header.Set(AdminTokenHeader, *tt.header)
			}
			rec := httptest.NewRecorder()
			RequireAdminToken(tt.token)(next).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if wantCalled := tt.wantStatus == http.StatusAccepted; isCalled != wantCalled {
				t.Errorf("next called = %v, want %v", isCalled, wantCalled)
			}
			if tt.wantStatus == http.StatusUnauthorized {
				if code := decodeError(t, rec.Body.Bytes()).Code; code != "UNAUTHORIZED" {
					t.Errorf("code = %q, want UNAUTHORIZED", code)
				}
			}
		})
	}
}

func TestRequestID(t *testing.T) {
	generated := regexp.MustCompile(`^[0-9a-f]{16}$`)
	tests := []struct {
		name     string
		incoming string
		want     *regexp.Regexp
	}{
		{"generated when absent", "", generated},
		{"well-formed incoming ID is reused", "edge-42.abc_Z", regexp.MustCompile(`^edge-42\.abc_Z$`)},
		{"malformed incoming ID is replaced", "bad id\"", generated},
		{"oversized incoming ID is replaced", strings.Repeat("a", 65), generated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctxID string
			next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				ctxID = RequestIDFrom(r.Context())
			})
			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			if tt.incoming != "" {
				req.Header.Set(RequestIDHeader, tt.incoming)
			}
			rec := httptest.NewRecorder()
			RequestID(next).ServeHTTP(rec, req)

			got := rec.Header().Get(RequestIDHeader)
			if !tt.want.MatchString(got) {
				t.Errorf("X-Request-ID = %q, want match %s", got, tt.want)
			}
			if ctxID != got {
				t.Errorf("context ID = %q, header ID = %q", ctxID, got)
			}
		})
	}
}

func TestCORS(t *testing.T) {
	allowed := []string{"http://localhost:3000"}
	tests := []struct {
		name          string
		origins       []string
		method        string
		origin        string
		isPreflight   bool
		wantAllowed   string
		wantVary      bool
		wantStatus    int
		wantNextCalls bool
	}{
		{"allowed origin", allowed, http.MethodGet, "http://localhost:3000", false, "http://localhost:3000", true, http.StatusOK, true},
		{"origin match ignores case", allowed, http.MethodGet, "http://LOCALHOST:3000", false, "http://LOCALHOST:3000", true, http.StatusOK, true},
		{"other origin", allowed, http.MethodGet, "https://evil.example", false, "", true, http.StatusOK, true},
		{"no origin header", allowed, http.MethodGet, "", false, "", true, http.StatusOK, true},
		{"preflight from allowed origin", allowed, http.MethodOptions, "http://localhost:3000", true, "http://localhost:3000", true, http.StatusNoContent, false},
		{"preflight from other origin", allowed, http.MethodOptions, "https://evil.example", true, "", true, http.StatusOK, true},
		{"no allowed origins configured", nil, http.MethodGet, "http://localhost:3000", false, "", false, http.StatusOK, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				isCalled = true
				w.WriteHeader(http.StatusOK)
			})
			req := httptest.NewRequest(tt.method, "/api/v1/options", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.isPreflight {
				req.Header.Set("Access-Control-Request-Method", http.MethodGet)
			}
			rec := httptest.NewRecorder()
			CORS(tt.origins)(next).ServeHTTP(rec, req)

			h := rec.Header()
			if got := h.Get("Access-Control-Allow-Origin"); got != tt.wantAllowed {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantAllowed)
			}
			if hasVary := h.Get("Vary") == "Origin"; hasVary != tt.wantVary {
				t.Errorf("Vary = %q, want Origin: %v", h.Get("Vary"), tt.wantVary)
			}
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if isCalled != tt.wantNextCalls {
				t.Errorf("next called = %v, want %v", isCalled, tt.wantNextCalls)
			}
			if tt.isPreflight && tt.wantAllowed != "" && h.Get("Access-Control-Allow-Methods") != http.MethodGet {
				t.Errorf("Access-Control-Allow-Methods = %q, want GET only", h.Get("Access-Control-Allow-Methods"))
			}
			if !tt.isPreflight && tt.wantAllowed != "" && !strings.Contains(h.Get("Access-Control-Expose-Headers"), RequestIDHeader) {
				t.Errorf("Access-Control-Expose-Headers = %q, want %s", h.Get("Access-Control-Expose-Headers"), RequestIDHeader)
			}
		})
	}
}

func TestLogRequests(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		handler   http.Handler
		wantLevel string
		wantErr   string
	}{
		{
			name: "client error at info",
			path: "/api/v1/movies/abc",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				WriteError(w, r, ErrNotFound)
			}),
			wantLevel: "INFO",
		},
		{
			name: "server error at error with cause",
			path: "/api/v1/options",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				WriteError(w, r, errors.New("snapshot exploded"))
			}),
			wantLevel: "ERROR",
			wantErr:   "snapshot exploded",
		},
		{
			name:      "healthz at debug",
			path:      HealthzPath,
			handler:   okHandler(http.StatusOK),
			wantLevel: "DEBUG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, buf := bufferLogger()
			h := Chain(tt.handler, RequestID, LogRequests(logger))
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set(AdminTokenHeader, "super-secret-token")
			h.ServeHTTP(httptest.NewRecorder(), req)

			var line map[string]any
			if err := json.Unmarshal(buf.Bytes(), &line); err != nil {
				t.Fatalf("log is not one JSON line: %v\n%s", err, buf)
			}
			for _, key := range []string{"request_id", "method", "path", "status", "duration_ms"} {
				if _, ok := line[key]; !ok {
					t.Errorf("log line misses %q: %s", key, buf)
				}
			}
			if line["level"] != tt.wantLevel {
				t.Errorf("level = %v, want %s", line["level"], tt.wantLevel)
			}
			if tt.wantErr != "" && !strings.Contains(buf.String(), tt.wantErr) {
				t.Errorf("log line misses the cause %q: %s", tt.wantErr, buf)
			}
			if strings.Contains(buf.String(), "super-secret-token") {
				t.Errorf("log line leaks the admin token: %s", buf)
			}
		})
	}
}

func TestRecover(t *testing.T) {
	logger, buf := bufferLogger()
	panicky := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") })
	h := Chain(panicky, RequestID, LogRequests(logger), Recover)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/options", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	assertJSON(t, rec.Body.Bytes(), `{"code":"INTERNAL_ERROR","message":"internal server error","details":null}`)
	if !strings.Contains(buf.String(), "panic: boom") || !strings.Contains(buf.String(), `"level":"ERROR"`) {
		t.Errorf("panic was not logged at error: %s", buf)
	}
}

func TestChainOrder(t *testing.T) {
	var order []string
	mark := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}
	Chain(okHandler(http.StatusOK), mark("outer"), mark("inner")).
		ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if strings.Join(order, ",") != "outer,inner" {
		t.Errorf("order = %v, want outer then inner", order)
	}
}

func ptr[T any](v T) *T { return &v }
