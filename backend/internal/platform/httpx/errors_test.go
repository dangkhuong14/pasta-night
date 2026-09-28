package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pasta_night/be/internal/domain"
)

func TestWriteErrorMapsEveryCode(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{"validation", domain.ErrValidation, http.StatusBadRequest, "VALIDATION_ERROR"},
		{"unauthorized", ErrUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED"},
		{"route not found", ErrNotFound, http.StatusNotFound, "NOT_FOUND"},
		{"option not found", domain.ErrOptionNotFound, http.StatusNotFound, "OPTION_NOT_FOUND"},
		{"movie not found", domain.ErrMovieNotFound, http.StatusNotFound, "MOVIE_NOT_FOUND"},
		{"method not allowed", ErrMethodNotAllowed, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED"},
		{"cache not ready", domain.ErrCacheNotReady, http.StatusServiceUnavailable, "CACHE_NOT_READY"},
		{"wrapped sentinel", fmt.Errorf("get movie 7: %w", domain.ErrMovieNotFound), http.StatusNotFound, "MOVIE_NOT_FOUND"},
		{"unmapped error", errors.New("disk on fire"), http.StatusInternalServerError, "INTERNAL_ERROR"},
	}
	tested := make(map[string]bool)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			WriteError(rec, httptest.NewRequest(http.MethodGet, "/x", nil), tt.err)
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q", ct)
			}
			body := decodeError(t, rec.Body.Bytes())
			if body.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", body.Code, tt.wantCode)
			}
			if body.Details != nil {
				t.Errorf("details = %v, want null", body.Details)
			}
			if !strings.Contains(rec.Body.String(), `"details":null`) {
				t.Errorf("details must be present as null: %s", rec.Body)
			}
		})
		tested[tt.wantCode] = true
	}
	for _, api := range errorMap {
		if !tested[api.code] {
			t.Errorf("errorMap entry %s has no test row", api.code)
		}
	}
}

func TestWriteErrorWithDetails(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "custom message and details",
			err: WithDetails(domain.ErrOptionNotFound, `option "family" does not exist`,
				map[string]any{"option_id": "family"}),
			want: `{"code":"OPTION_NOT_FOUND","message":"option \"family\" does not exist","details":{"option_id":"family"}}`,
		},
		{
			name: "empty message keeps the default",
			err:  WithDetails(fmt.Errorf("get movie 42: %w", domain.ErrMovieNotFound), "", map[string]any{"movie_id": 42}),
			want: `{"code":"MOVIE_NOT_FOUND","message":"movie is not in any current recommendation list","details":{"movie_id":42}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			WriteError(rec, httptest.NewRequest(http.MethodGet, "/x", nil), tt.err)
			assertJSON(t, rec.Body.Bytes(), tt.want)
		})
	}
}

func TestWriteErrorHidesInternalCause(t *testing.T) {
	rec := httptest.NewRecorder()
	err := WithDetails(errors.New("open /srv/secret: permission denied"), "leaky message", map[string]any{"path": "/srv/secret"})
	WriteError(rec, httptest.NewRequest(http.MethodGet, "/x", nil), err)

	assertJSON(t, rec.Body.Bytes(), `{"code":"INTERNAL_ERROR","message":"internal server error","details":null}`)
}

func TestWriteErrorRetryAfter(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{domain.ErrCacheNotReady, "30"},
		{domain.ErrMovieNotFound, ""},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		WriteError(rec, httptest.NewRequest(http.MethodGet, "/x", nil), tt.err)
		if got := rec.Header().Get("Retry-After"); got != tt.want {
			t.Errorf("%v: Retry-After = %q, want %q", tt.err, got, tt.want)
		}
	}
}
