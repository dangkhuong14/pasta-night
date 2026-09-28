package refresh

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"pasta_night/be/internal/platform/httpx"
)

func TestTriggerRefresh(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		wantStatus  int
		wantBody    string
		wantPending []string
	}{
		{
			name:        "all options",
			query:       "",
			wantStatus:  http.StatusAccepted,
			wantBody:    `{"data":{"accepted":["netflix-chill","solo","friends"]},"meta":{}}`,
			wantPending: []string{"netflix-chill", "solo", "friends"},
		},
		{
			name:        "empty option_id means all options",
			query:       "?option_id=",
			wantStatus:  http.StatusAccepted,
			wantBody:    `{"data":{"accepted":["netflix-chill","solo","friends"]},"meta":{}}`,
			wantPending: []string{"netflix-chill", "solo", "friends"},
		},
		{
			name:        "one option",
			query:       "?option_id=friends",
			wantStatus:  http.StatusAccepted,
			wantBody:    `{"data":{"accepted":["friends"]},"meta":{}}`,
			wantPending: []string{"friends"},
		},
		{
			name:       "unknown option",
			query:      "?option_id=family",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":"OPTION_NOT_FOUND","message":"option \"family\" does not exist","details":{"option_id":"family"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, _, _ := newTestJob(t, newFakeTMDB())
			mux := http.NewServeMux()
			NewHandler(j, testOptions).Register(mux, httpx.RequireAdminToken("s3cret-token"))

			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/refresh"+tt.query, nil)
			req.Header.Set(httpx.AdminTokenHeader, "s3cret-token")
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			var got, want any
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("body: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.wantBody), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("body =\n%s\nwant\n%s", rec.Body, tt.wantBody)
			}
			if pending := j.takePending(); !reflect.DeepEqual(pending, tt.wantPending) {
				t.Errorf("triggered %v, want %v", pending, tt.wantPending)
			}
		})
	}
}

func TestTriggerRefreshRequiresToken(t *testing.T) {
	j, _, _ := newTestJob(t, newFakeTMDB())
	mux := http.NewServeMux()
	NewHandler(j, testOptions).Register(mux, httpx.RequireAdminToken("s3cret-token"))

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/admin/refresh", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
	if pending := j.takePending(); len(pending) != 0 {
		t.Errorf("an unauthorized request triggered %v", pending)
	}
}
