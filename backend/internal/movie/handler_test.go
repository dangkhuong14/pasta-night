package movie

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"pasta_night/be/internal/cache"
	"pasta_night/be/internal/domain"
)

var testFetchedAt = time.Date(2026, 9, 27, 10, 0, 0, 0, time.UTC)

func ptr[T any](v T) *T { return &v }

// matrix is a fully populated movie.
var matrix = domain.MovieDetail{
	ID:             603,
	Title:          "Ma Trận",
	OriginalTitle:  "The Matrix",
	Overview:       "Một hacker phát hiện thế giới anh đang sống chỉ là một mô phỏng.",
	Tagline:        "",
	PosterURL:      ptr("https://image.tmdb.org/t/p/w500/p.jpg"),
	BackdropURL:    nil,
	ReleaseYear:    ptr(1999),
	Rating:         8.2,
	VoteCount:      26000,
	RuntimeMinutes: ptr(136),
	Genres:         []string{"Phim Hành Động"},
	Directors:      []string{"Lana Wachowski", "Lilly Wachowski"},
	Cast:           []string{"Keanu Reeves"},
	TrailerURL:     ptr("https://www.youtube.com/watch?v=k"),
	Providers: []domain.Provider{
		{ID: 8, Name: "Netflix", LogoURL: ptr("https://image.tmdb.org/t/p/w92/n.jpg"), Type: domain.ProviderFlatrate},
	},
	FetchedAt: testFetchedAt.Add(5 * time.Second),
}

// sparse is a movie with every nullable field null and empty arrays.
var sparse = domain.MovieDetail{
	ID:            680,
	Title:         "Pulp Fiction",
	OriginalTitle: "Pulp Fiction",
	FetchedAt:     testFetchedAt,
}

// newTestMux serves the movie routes over a store where "friends" has a
// list, "solo" has none yet, and "family" is not an option.
func newTestMux(t *testing.T) *http.ServeMux {
	t.Helper()
	store := cache.NewStore(t.TempDir(), []string{"friends", "solo"}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	list := domain.MovieList{OptionID: "friends", FetchedAt: testFetchedAt, MovieIDs: []int{603, 680}}
	if err := store.PublishList(context.Background(), list, []domain.MovieDetail{matrix, sparse}); err != nil {
		t.Fatal(err)
	}
	options := []domain.Option{{ID: "friends"}, {ID: "solo"}}
	svc := NewService(NewRepository(store), 169*time.Hour)
	svc.now = func() time.Time { return testFetchedAt.Add(time.Hour) }
	mux := http.NewServeMux()
	NewHandler(svc, options).Register(mux)
	return mux
}

func TestHandler(t *testing.T) {
	mux := newTestMux(t)
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "recommendations in rank order",
			path:       "/api/v1/options/friends/recommendations",
			wantStatus: http.StatusOK,
			wantBody: `{"data":[
				{"id":603,"title":"Ma Trận","overview":"Một hacker phát hiện thế giới anh đang sống chỉ là một mô phỏng.",
				 "poster_url":"https://image.tmdb.org/t/p/w500/p.jpg","release_year":1999,"rating":8.2,"runtime_minutes":136,
				 "genres":["Phim Hành Động"],
				 "providers":[{"id":8,"name":"Netflix","logo_url":"https://image.tmdb.org/t/p/w92/n.jpg","type":"flatrate"}]},
				{"id":680,"title":"Pulp Fiction","overview":"","poster_url":null,"release_year":null,"rating":0,
				 "runtime_minutes":null,"genres":[],"providers":[]}
			],"meta":{"option_id":"friends","total":2,"fetched_at":"2026-09-27T10:00:00Z","stale":false}}`,
		},
		{
			name:       "unknown option",
			path:       "/api/v1/options/family/recommendations",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":"OPTION_NOT_FOUND","message":"option \"family\" does not exist","details":{"option_id":"family"}}`,
		},
		{
			name:       "option without a list yet",
			path:       "/api/v1/options/solo/recommendations",
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   `{"code":"CACHE_NOT_READY","message":"recommendations are being prepared, retry later","details":{"option_id":"solo"}}`,
		},
		{
			name:       "movie detail",
			path:       "/api/v1/movies/603",
			wantStatus: http.StatusOK,
			wantBody: `{"data":{
				"id":603,"title":"Ma Trận","original_title":"The Matrix",
				"overview":"Một hacker phát hiện thế giới anh đang sống chỉ là một mô phỏng.","tagline":"",
				"poster_url":"https://image.tmdb.org/t/p/w500/p.jpg","backdrop_url":null,"release_year":1999,
				"rating":8.2,"vote_count":26000,"runtime_minutes":136,"genres":["Phim Hành Động"],
				"directors":["Lana Wachowski","Lilly Wachowski"],"cast":["Keanu Reeves"],
				"trailer_url":"https://www.youtube.com/watch?v=k",
				"providers":[{"id":8,"name":"Netflix","logo_url":"https://image.tmdb.org/t/p/w92/n.jpg","type":"flatrate"}]
			},"meta":{"fetched_at":"2026-09-27T10:00:05Z"}}`,
		},
		{
			name:       "sparse movie keeps nulls and empty arrays",
			path:       "/api/v1/movies/680",
			wantStatus: http.StatusOK,
			wantBody: `{"data":{
				"id":680,"title":"Pulp Fiction","original_title":"Pulp Fiction","overview":"","tagline":"",
				"poster_url":null,"backdrop_url":null,"release_year":null,"rating":0,"vote_count":0,
				"runtime_minutes":null,"genres":[],"directors":[],"cast":[],"trailer_url":null,"providers":[]
			},"meta":{"fetched_at":"2026-09-27T10:00:00Z"}}`,
		},
		{
			name:       "movie not in any list",
			path:       "/api/v1/movies/42",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":"MOVIE_NOT_FOUND","message":"movie is not in any current recommendation list","details":{"movie_id":42}}`,
		},
		{
			name:       "movie_id not a number",
			path:       "/api/v1/movies/abc",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"VALIDATION_ERROR","message":"movie_id must be a positive integer","details":{"movie_id":"abc"}}`,
		},
		{
			name:       "movie_id zero",
			path:       "/api/v1/movies/0",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"VALIDATION_ERROR","message":"movie_id must be a positive integer","details":{"movie_id":"0"}}`,
		},
		{
			name:       "movie_id negative",
			path:       "/api/v1/movies/-5",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"VALIDATION_ERROR","message":"movie_id must be a positive integer","details":{"movie_id":"-5"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			var got, want any
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("body is not JSON: %v\n%s", err, rec.Body)
			}
			if err := json.Unmarshal([]byte(tt.wantBody), &want); err != nil {
				t.Fatalf("bad expected JSON: %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("body =\n%s\nwant\n%s", rec.Body, tt.wantBody)
			}
		})
	}
}

func TestRecommendationsStaleFlag(t *testing.T) {
	store := cache.NewStore(t.TempDir(), []string{"friends"}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	list := domain.MovieList{OptionID: "friends", FetchedAt: testFetchedAt, MovieIDs: []int{680}}
	if err := store.PublishList(context.Background(), list, []domain.MovieDetail{sparse}); err != nil {
		t.Fatal(err)
	}
	svc := NewService(NewRepository(store), 169*time.Hour)
	svc.now = func() time.Time { return testFetchedAt.Add(200 * time.Hour) }
	mux := http.NewServeMux()
	NewHandler(svc, []domain.Option{{ID: "friends"}}).Register(mux)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/options/friends/recommendations", nil))
	var body struct {
		Meta struct {
			IsStale bool `json:"stale"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK || !body.Meta.IsStale {
		t.Errorf("status %d, stale %v: stale data must still be served, flagged", rec.Code, body.Meta.IsStale)
	}
}
