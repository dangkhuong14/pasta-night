package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFallback(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/options", func(w http.ResponseWriter, r *http.Request) {
		WriteData(w, r, http.StatusOK, []string{}, nil)
	})
	mux.HandleFunc("POST /api/v1/admin/refresh", func(w http.ResponseWriter, r *http.Request) {
		WriteData(w, r, http.StatusAccepted, map[string]any{}, nil)
	})
	h := Fallback(mux)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantCode   string // empty for success
		wantAllow  string
	}{
		{"matched route", http.MethodGet, "/api/v1/options", http.StatusOK, "", ""},
		{"HEAD matches a GET route", http.MethodHead, "/api/v1/options", http.StatusOK, "", ""},
		{"unknown path", http.MethodGet, "/api/v1/nope", http.StatusNotFound, "NOT_FOUND", ""},
		{"trailing slash is a different path", http.MethodGet, "/api/v1/options/", http.StatusNotFound, "NOT_FOUND", ""},
		{"wrong method", http.MethodDelete, "/api/v1/options", http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GET"},
		{"GET on a POST route", http.MethodGet, "/api/v1/admin/refresh", http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "POST"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body %s", rec.Code, tt.wantStatus, rec.Body)
			}
			if tt.wantCode == "" {
				return
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q, want JSON", ct)
			}
			if code := decodeError(t, rec.Body.Bytes()).Code; code != tt.wantCode {
				t.Errorf("code = %q, want %q", code, tt.wantCode)
			}
			if !strings.Contains(rec.Header().Get("Allow"), tt.wantAllow) {
				t.Errorf("Allow = %q, want it to contain %q", rec.Header().Get("Allow"), tt.wantAllow)
			}
		})
	}
}
