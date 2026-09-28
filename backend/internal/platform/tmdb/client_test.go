package tmdb

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

const testToken = "test-read-token"

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return data
}

// newTestClient starts a server with handler and returns a client for it.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewClient(srv.URL, testToken, "vi-VN")
}

// captured is the part of a request the server saw.
type captured struct {
	path   string
	query  url.Values
	header http.Header
}

func serveFixture(t *testing.T, name string, got *captured) http.HandlerFunc {
	body := fixture(t, name)
	return func(w http.ResponseWriter, r *http.Request) {
		*got = captured{path: r.URL.Path, query: r.URL.Query(), header: r.Header.Clone()}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

func assertQuery(t *testing.T, got url.Values, want map[string]string, absent ...string) {
	t.Helper()
	for k, v := range want {
		if got.Get(k) != v {
			t.Errorf("query %s = %q, want %q", k, got.Get(k), v)
		}
	}
	for _, k := range absent {
		if got.Has(k) {
			t.Errorf("query %s = %q, want absent", k, got.Get(k))
		}
	}
}

func TestDiscover(t *testing.T) {
	tests := []struct {
		name   string
		query  DiscoverQuery
		want   map[string]string
		absent []string
	}{
		{
			name:  "genres only",
			query: DiscoverQuery{Page: 1, WithGenres: "35|28|27", VoteAverageGTE: 6.5, VoteCountGTE: 200, SortBy: "popularity.desc"},
			want: map[string]string{
				"language": "vi-VN", "include_adult": "false", "include_video": "false",
				"with_genres": "35|28|27", "vote_average.gte": "6.5", "vote_count.gte": "200",
				"sort_by": "popularity.desc", "page": "1",
			},
			absent: []string{"with_watch_providers", "watch_region", "with_watch_monetization_types"},
		},
		{
			name: "with watch providers",
			query: DiscoverQuery{Page: 2, WithGenres: "10749,35", VoteAverageGTE: 7, VoteCountGTE: 200,
				SortBy: "popularity.desc", WithWatchProviders: "8", WatchRegion: "VN"},
			want: map[string]string{
				"with_genres": "10749,35", "vote_average.gte": "7", "page": "2",
				"with_watch_providers": "8", "watch_region": "VN", "with_watch_monetization_types": "flatrate",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got captured
			c := newTestClient(t, serveFixture(t, "discover_page1.json", &got))
			resp, err := c.Discover(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("Discover: %v", err)
			}
			if got.path != "/discover/movie" {
				t.Errorf("path = %q", got.path)
			}
			if got.header.Get("Authorization") != "Bearer "+testToken || got.header.Get("Accept") != "application/json" {
				t.Errorf("headers = %v", got.header)
			}
			assertQuery(t, got.query, tt.want, tt.absent...)
			if resp.Page != 1 || resp.TotalPages != 12 || len(resp.Results) != 4 {
				t.Errorf("response = %+v", resp)
			}
			if first := resp.Results[0]; first.ID != 603 || first.Adult {
				t.Errorf("first result = %+v", first)
			}
			if !resp.Results[2].Adult {
				t.Error("adult flag not decoded")
			}
		})
	}
}

func TestMovieDetails(t *testing.T) {
	var got captured
	c := newTestClient(t, serveFixture(t, "movie_603_full.json", &got))
	m, err := c.MovieDetails(context.Background(), 603)
	if err != nil {
		t.Fatalf("MovieDetails: %v", err)
	}
	if got.path != "/movie/603" {
		t.Errorf("path = %q", got.path)
	}
	assertQuery(t, got.query, map[string]string{
		"language":               "vi-VN",
		"append_to_response":     "credits,videos,watch/providers",
		"include_video_language": "vi,en",
	})
	if m.ID != 603 || m.Title != "Ma Trận" || m.Runtime != 136 || m.ReleaseDate != "1999-03-30" {
		t.Errorf("movie = %+v", m)
	}
	if len(m.Credits.Cast) != 6 || len(m.Credits.Crew) != 5 || len(m.Videos.Results) != 5 {
		t.Errorf("appended objects not decoded: %d cast, %d crew, %d videos",
			len(m.Credits.Cast), len(m.Credits.Crew), len(m.Videos.Results))
	}
	if vn := m.WatchProviders.Results["VN"]; len(vn.Flatrate) != 1 || vn.Flatrate[0].ProviderName != "Netflix" {
		t.Errorf("watch/providers VN = %+v", vn)
	}
}

func TestErrorSentinels(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		want     error
		wantCode int
	}{
		{"401 invalid key", http.StatusUnauthorized, string(fixture(t, "error_401.json")), ErrUnauthorized, 7},
		{"404 movie removed", http.StatusNotFound, `{"success":false,"status_code":34,"status_message":"The resource you requested could not be found."}`, ErrNotFound, 34},
		{"400 bad params", http.StatusBadRequest, `{"success":false,"status_code":22,"status_message":"Invalid page."}`, ErrBadRequest, 22},
		{"422 bad params", http.StatusUnprocessableEntity, `{"success":false,"status_code":18,"status_message":"Validation failed."}`, ErrBadRequest, 18},
		{"500", http.StatusInternalServerError, `{"success":false,"status_code":11,"status_message":"Internal error."}`, ErrUnavailable, 11},
		{"503 with html body", http.StatusServiceUnavailable, "<html>down</html>", ErrUnavailable, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			})
			_, err := c.MovieDetails(context.Background(), 603)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("err = %v, want *APIError", err)
			}
			if apiErr.HTTPStatus != tt.status || apiErr.Code != tt.wantCode || apiErr.Message == "" {
				t.Errorf("APIError = %+v", apiErr)
			}
		})
	}
}

func TestRateLimitRetry(t *testing.T) {
	tests := []struct {
		name        string
		failures    int32 // 429 responses before success
		wantCalls   int32
		wantLimited bool
	}{
		{"succeeds after two retries", 2, 3, false},
		{"gives up after two retries", 10, 3, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls atomic.Int32
			body429, body200 := fixture(t, "error_429.json"), fixture(t, "movie_603_full.json")
			c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				if calls.Add(1) <= tt.failures {
					w.Header().Set("Retry-After", "0")
					w.WriteHeader(http.StatusTooManyRequests)
					_, _ = w.Write(body429)
					return
				}
				_, _ = w.Write(body200)
			})
			_, err := c.MovieDetails(context.Background(), 603)
			if isLimited := errors.Is(err, ErrRateLimited); isLimited != tt.wantLimited {
				t.Errorf("err = %v, want rate limited: %v", err, tt.wantLimited)
			}
			if !tt.wantLimited && err != nil {
				t.Errorf("err = %v, want success", err)
			}
			if calls.Load() != tt.wantCalls {
				t.Errorf("calls = %d, want %d", calls.Load(), tt.wantCalls)
			}
		})
	}
}

func TestNoRetryOnOtherErrors(t *testing.T) {
	var calls atomic.Int32
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	if _, err := c.MovieDetails(context.Background(), 603); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("err = %v, want ErrUnavailable", err)
	}
	if calls.Load() != 1 {
		t.Errorf("calls = %d, want 1 (only 429 is retried)", calls.Load())
	}
}

func TestMalformedBody(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	})
	if _, err := c.MovieDetails(context.Background(), 603); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("err = %v, want ErrUnavailable", err)
	}
}

func TestNetworkError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close() // nothing listens on this address anymore
	c := NewClient(srv.URL, testToken, "vi-VN")
	if _, err := c.MovieDetails(context.Background(), 603); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("err = %v, want ErrUnavailable", err)
	}
}

func TestCancelledContextIsNotAnOutage(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(fixture(t, "movie_603_full.json"))
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.MovieDetails(ctx, 603)
	if !errors.Is(err, context.Canceled) || errors.Is(err, ErrUnavailable) {
		t.Fatalf("err = %v, want context.Canceled only", err)
	}
}

func TestRetryWaitRespectsContext(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := c.MovieDetails(ctx, 603)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("waited %v; the retry wait must stop when the context ends", elapsed)
	}
}

func TestParseRetryAfter(t *testing.T) {
	tests := []struct {
		in   string
		want time.Duration
	}{
		{"0", 0},
		{"3", 3 * time.Second},
		{" 7 ", 7 * time.Second},
		{"", defaultRetryAfter},
		{"soon", defaultRetryAfter},
		{"-1", defaultRetryAfter},
		{"Wed, 21 Oct 2015 07:28:00 GMT", 0}, // a date in the past
	}
	for _, tt := range tests {
		if got := parseRetryAfter(tt.in); got != tt.want {
			t.Errorf("parseRetryAfter(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
