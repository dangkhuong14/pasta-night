package option

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"pasta_night/be/internal/domain"
)

func TestListOptions(t *testing.T) {
	options := []domain.Option{
		{ID: "netflix-chill", Label: "Netflix & Chill", Description: "Cho hai người", Icon: "heart",
			Discover: domain.DiscoverParams{GenreIDs: []int{10749}, WatchProviderIDs: []int{8}}},
		{ID: "solo", Label: "Một mình", Description: "Thời gian cho riêng bạn", Icon: "user"},
	}
	tests := []struct {
		name    string
		options []domain.Option
		want    string
	}{
		{
			name:    "display order, no discover params",
			options: options,
			want: `{"data":[
				{"id":"netflix-chill","label":"Netflix & Chill","description":"Cho hai người","icon":"heart"},
				{"id":"solo","label":"Một mình","description":"Thời gian cho riêng bạn","icon":"user"}
			],"meta":{}}`,
		},
		{name: "no options is an empty array", options: nil, want: `{"data":[],"meta":{}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			NewHandler(NewService(tt.options)).Register(mux)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/options", nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}
			if got := rec.Header().Get("Cache-Control"); got != "public, max-age=300" {
				t.Errorf("Cache-Control = %q", got)
			}
			var got, want any
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("body: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.want), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("body =\n%s\nwant\n%s", rec.Body, tt.want)
			}
		})
	}
}
