package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteData(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		status    int
		meta      any
		wantCache string
		wantBody  string
	}{
		{"GET success is cacheable, nil meta is {}", http.MethodGet, http.StatusOK, nil,
			"public, max-age=300", `{"data":{"title":"Netflix & Chill"},"meta":{}}`},
		{"meta passes through", http.MethodGet, http.StatusOK, map[string]int{"total": 1},
			"public, max-age=300", `{"data":{"title":"Netflix & Chill"},"meta":{"total":1}}`},
		{"POST is not cacheable", http.MethodPost, http.StatusAccepted, nil,
			"", `{"data":{"title":"Netflix & Chill"},"meta":{}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			data := map[string]string{"title": "Netflix & Chill"}
			WriteData(rec, httptest.NewRequest(tt.method, "/x", nil), tt.status, data, tt.meta)

			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}
			if got := rec.Header().Get("Cache-Control"); got != tt.wantCache {
				t.Errorf("Cache-Control = %q, want %q", got, tt.wantCache)
			}
			if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf("X-Content-Type-Options = %q", got)
			}
			assertJSON(t, rec.Body.Bytes(), tt.wantBody)
		})
	}
}

func TestWriteJSONKeepsHTMLCharacters(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, httptest.NewRequest(http.MethodGet, "/x", nil), http.StatusOK, map[string]string{"label": "Netflix & Chill"})
	if got, want := rec.Body.String(), "{\"label\":\"Netflix & Chill\"}\n"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestWriteJSONEncodeFailure(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, httptest.NewRequest(http.MethodGet, "/x", nil), http.StatusOK, map[string]any{"bad": make(chan int)})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	assertJSON(t, rec.Body.Bytes(), `{"code":"INTERNAL_ERROR","message":"internal server error","details":null}`)
}

func TestHealthz(t *testing.T) {
	for _, isReady := range []bool{true, false} {
		rec := httptest.NewRecorder()
		Healthz(func() bool { return isReady }).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, HealthzPath, nil))

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", rec.Code)
		}
		if got := rec.Header().Get("Cache-Control"); got != "no-store" {
			t.Errorf("Cache-Control = %q, want no-store", got)
		}
		want := `{"status":"ok","cache_ready":false}`
		if isReady {
			want = `{"status":"ok","cache_ready":true}`
		}
		assertJSON(t, rec.Body.Bytes(), want)
	}
}
