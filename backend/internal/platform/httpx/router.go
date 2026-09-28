package httpx

import "net/http"

// Fallback serves mux and replaces ServeMux's plain-text 404 and 405 replies
// with the JSON error envelope (API_SPEC.md §6). The 405 reply keeps the
// Allow header ServeMux sets.
func Fallback(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pattern := mux.Handler(r); pattern != "" {
			mux.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(&fallbackWriter{ResponseWriter: w, r: r}, r)
	})
}

// fallbackWriter intercepts the status ServeMux writes for unmatched requests.
type fallbackWriter struct {
	http.ResponseWriter
	r             *http.Request
	isIntercepted bool
}

func (f *fallbackWriter) WriteHeader(code int) {
	switch code {
	case http.StatusNotFound:
		f.isIntercepted = true
		WriteError(f.ResponseWriter, f.r, ErrNotFound)
	case http.StatusMethodNotAllowed:
		f.isIntercepted = true
		WriteError(f.ResponseWriter, f.r, ErrMethodNotAllowed)
	default:
		f.ResponseWriter.WriteHeader(code)
	}
}

func (f *fallbackWriter) Write(b []byte) (int, error) {
	if f.isIntercepted {
		return len(b), nil // drop ServeMux's plain-text body
	}
	return f.ResponseWriter.Write(b)
}
