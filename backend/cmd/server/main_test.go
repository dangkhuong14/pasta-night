package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/config"
	"pasta_night/be/internal/platform/httpx"
	"pasta_night/be/internal/platform/tmdb"
)

const (
	testTMDBToken  = "test-tmdb-token"
	testAdminToken = "test-admin-token"
)

// fakeTMDB serves two discover pages (movies 1-40) and a detail for every movie.
type fakeTMDB struct {
	*httptest.Server
	discoverCalls atomic.Int64
	detailCalls   atomic.Int64
	isDown        atomic.Bool
}

func newFakeTMDB(t *testing.T) *fakeTMDB {
	t.Helper()
	f := &fakeTMDB{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /discover/movie", func(w http.ResponseWriter, r *http.Request) {
		if !f.authorize(w, r) {
			return
		}
		f.discoverCalls.Add(1)
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		results := make([]map[string]any, 0, 20)
		for id := (page-1)*20 + 1; id <= page*20; id++ {
			results = append(results, map[string]any{"id": id, "adult": false})
		}
		writeJSON(w, map[string]any{"page": page, "total_pages": 2, "results": results})
	})
	mux.HandleFunc("GET /movie/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !f.authorize(w, r) {
			return
		}
		f.detailCalls.Add(1)
		id, _ := strconv.Atoi(r.PathValue("id"))
		writeJSON(w, map[string]any{
			"id": id, "title": fmt.Sprintf("Phim %d", id), "original_title": fmt.Sprintf("Movie %d", id),
			"poster_path": fmt.Sprintf("/p%d.jpg", id), "release_date": "2020-01-01",
			"vote_average": 7.46, "vote_count": 1000, "runtime": 100,
			"genres":  []map[string]any{{"id": 35, "name": "Phim Hài"}},
			"credits": map[string]any{"cast": []map[string]any{{"name": "Diễn viên", "order": 0}}, "crew": []map[string]any{}},
			"videos":  map[string]any{"results": []map[string]any{}},
			"watch/providers": map[string]any{"results": map[string]any{
				"VN": map[string]any{"flatrate": []map[string]any{
					{"provider_id": 8, "provider_name": "Netflix", "logo_path": "/n.jpg", "display_priority": 1},
				}},
			}},
		})
	})
	f.Server = httptest.NewServer(mux)
	t.Cleanup(f.Close)
	return f
}

func (f *fakeTMDB) authorize(w http.ResponseWriter, r *http.Request) bool {
	switch {
	case f.isDown.Load():
		w.WriteHeader(http.StatusServiceUnavailable)
		return false
	case r.Header.Get("Authorization") != "Bearer "+testTMDBToken:
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, map[string]any{"success": false, "status_code": 7, "status_message": "Invalid API key"})
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func loadShippedOptions(t *testing.T) []domain.Option {
	t.Helper()
	options, err := config.LoadOptions(filepath.Join("..", "..", optionsPath))
	if err != nil {
		t.Fatalf("load options: %v", err)
	}
	return options
}

func testConfig(cacheDir, adminToken string) config.Config {
	return config.Config{
		TMDBReadToken:        testTMDBToken,
		TMDBLanguage:         "vi-VN",
		CacheDir:             cacheDir,
		ListTTL:              168 * time.Hour,
		DetailTTL:            336 * time.Hour,
		RefreshCheckInterval: time.Hour,
		DiscoverPages:        2,
		TMDBConcurrency:      5,
		CORSAllowedOrigins:   []string{"http://localhost:3000"},
		AdminToken:           adminToken,
	}
}

func newTestApp(t *testing.T, cfg config.Config, fake *fakeTMDB) *app {
	t.Helper()
	client := tmdb.NewClient(fake.URL, cfg.TMDBReadToken, cfg.TMDBLanguage)
	a, err := newApp(context.Background(), cfg, loadShippedOptions(t), client, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}
	return a
}

// startJob runs the refresh job until the test ends.
func startJob(t *testing.T, a *app) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		a.job.Run(ctx)
	}()
	t.Cleanup(func() {
		cancel()
		<-done
	})
}

func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type response struct {
	status int
	header http.Header
	body   map[string]any
}

func call(t *testing.T, a *app, method, target string, header map[string]string) response {
	t.Helper()
	req := httptest.NewRequest(method, target, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	a.handler.ServeHTTP(rec, req)
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("%s %s: body is not a JSON object: %v\n%s", method, target, err, rec.Body)
	}
	if rec.Header().Get(httpx.RequestIDHeader) == "" {
		t.Errorf("%s %s: missing X-Request-ID", method, target)
	}
	return response{status: rec.Code, header: rec.Header(), body: body}
}

func (r response) code() any { return r.body["code"] }

func TestServerEndToEnd(t *testing.T) {
	fake := newFakeTMDB(t)
	a := newTestApp(t, testConfig(t.TempDir(), testAdminToken), fake)

	// Before the first refresh finishes.
	res := call(t, a, http.MethodGet, "/api/v1/options/friends/recommendations", nil)
	if res.status != http.StatusServiceUnavailable || res.code() != "CACHE_NOT_READY" || res.header.Get("Retry-After") != "30" {
		t.Fatalf("before refresh: %d %v Retry-After=%q", res.status, res.body, res.header.Get("Retry-After"))
	}
	if res = call(t, a, http.MethodGet, "/healthz", nil); res.body["cache_ready"] != false {
		t.Errorf("healthz before refresh = %v", res.body)
	}

	startJob(t, a)
	waitFor(t, "cache ready", a.store.IsReady)

	t.Run("options", func(t *testing.T) {
		res := call(t, a, http.MethodGet, "/api/v1/options", nil)
		if data, _ := res.body["data"].([]any); res.status != http.StatusOK || len(data) != 3 {
			t.Errorf("GET /options = %d %v", res.status, res.body)
		}
	})
	t.Run("recommendations", func(t *testing.T) {
		res := call(t, a, http.MethodGet, "/api/v1/options/friends/recommendations", nil)
		data, _ := res.body["data"].([]any)
		meta, _ := res.body["meta"].(map[string]any)
		if res.status != http.StatusOK || len(data) != domain.MaxListSize {
			t.Fatalf("recommendations = %d, %d movies", res.status, len(data))
		}
		if meta["option_id"] != "friends" || meta["total"] != float64(domain.MaxListSize) || meta["stale"] != false {
			t.Errorf("meta = %v", meta)
		}
		first, _ := data[0].(map[string]any)
		if first["id"] != float64(1) || first["title"] != "Phim 1" || first["rating"] != 7.5 {
			t.Errorf("first movie = %v", first)
		}
		if res.header.Get("Cache-Control") != "public, max-age=300" {
			t.Errorf("Cache-Control = %q", res.header.Get("Cache-Control"))
		}
	})
	t.Run("movie detail", func(t *testing.T) {
		res := call(t, a, http.MethodGet, "/api/v1/movies/40", nil)
		data, _ := res.body["data"].(map[string]any)
		if res.status != http.StatusOK || data["id"] != float64(40) {
			t.Errorf("GET /movies/40 = %d %v", res.status, res.body)
		}
	})
	t.Run("error codes", func(t *testing.T) {
		tests := []struct {
			method, target string
			wantStatus     int
			wantCode       string
		}{
			{http.MethodGet, "/api/v1/movies/41", http.StatusNotFound, "MOVIE_NOT_FOUND"},
			{http.MethodGet, "/api/v1/movies/abc", http.StatusBadRequest, "VALIDATION_ERROR"},
			{http.MethodGet, "/api/v1/options/family/recommendations", http.StatusNotFound, "OPTION_NOT_FOUND"},
			{http.MethodGet, "/api/v1/nope", http.StatusNotFound, "NOT_FOUND"},
			{http.MethodDelete, "/api/v1/options", http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED"},
		}
		for _, tt := range tests {
			res := call(t, a, tt.method, tt.target, nil)
			if res.status != tt.wantStatus || res.code() != tt.wantCode {
				t.Errorf("%s %s = %d %v, want %d %s", tt.method, tt.target, res.status, res.code(), tt.wantStatus, tt.wantCode)
			}
		}
	})
	t.Run("cors", func(t *testing.T) {
		res := call(t, a, http.MethodGet, "/api/v1/options", map[string]string{"Origin": "http://localhost:3000"})
		if res.header.Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Errorf("Access-Control-Allow-Origin = %q", res.header.Get("Access-Control-Allow-Origin"))
		}
	})
	t.Run("admin refresh", func(t *testing.T) {
		const target = "/api/v1/admin/refresh"
		if res := call(t, a, http.MethodPost, target, nil); res.status != http.StatusUnauthorized || res.code() != "UNAUTHORIZED" {
			t.Errorf("no token = %d %v", res.status, res.body)
		}
		wrong := map[string]string{httpx.AdminTokenHeader: "guess"}
		if res := call(t, a, http.MethodPost, target, wrong); res.status != http.StatusUnauthorized {
			t.Errorf("wrong token = %d %v", res.status, res.body)
		}
		admin := map[string]string{httpx.AdminTokenHeader: testAdminToken}
		if res := call(t, a, http.MethodPost, target+"?option_id=family", admin); res.code() != "OPTION_NOT_FOUND" {
			t.Errorf("unknown option = %d %v", res.status, res.body)
		}

		before := fake.discoverCalls.Load()
		res := call(t, a, http.MethodPost, target+"?option_id=friends", admin)
		data, _ := res.body["data"].(map[string]any)
		if accepted, _ := data["accepted"].([]any); res.status != http.StatusAccepted || len(accepted) != 1 || accepted[0] != "friends" {
			t.Fatalf("POST refresh = %d %v", res.status, res.body)
		}
		// The forced refresh ignores LIST_TTL: both discover pages are fetched again.
		waitFor(t, "forced refresh", func() bool { return fake.discoverCalls.Load() == before+2 })
	})
}

func TestAdminDisabledWithoutToken(t *testing.T) {
	a := newTestApp(t, testConfig(t.TempDir(), ""), newFakeTMDB(t))
	res := call(t, a, http.MethodPost, "/api/v1/admin/refresh", map[string]string{httpx.AdminTokenHeader: ""})
	if res.status != http.StatusNotFound || res.code() != "NOT_FOUND" {
		t.Errorf("admin without ADMIN_TOKEN = %d %v, want 404 NOT_FOUND", res.status, res.body)
	}
}

func TestCacheSurvivesRestartWithTMDBDown(t *testing.T) {
	fake := newFakeTMDB(t)
	dir := t.TempDir()
	first := newTestApp(t, testConfig(dir, ""), fake)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		first.job.Run(ctx)
	}()
	waitFor(t, "first app cache ready", first.store.IsReady)
	cancel()
	<-done

	// TMDB goes down; a restarted server still serves everything from disk.
	fake.isDown.Store(true)
	second := newTestApp(t, testConfig(dir, ""), fake)
	if !second.store.IsReady() {
		t.Fatal("restarted server is not ready from the disk cache")
	}
	res := call(t, second, http.MethodGet, "/api/v1/options/solo/recommendations", nil)
	if data, _ := res.body["data"].([]any); res.status != http.StatusOK || len(data) != domain.MaxListSize {
		t.Errorf("recommendations after restart = %d, %d movies", res.status, len(data))
	}
}

func TestRunRejectsInvalidConfig(t *testing.T) {
	err := run(context.Background(), func(string) string { return "" }, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err == nil {
		t.Fatal("run must fail without TMDB_READ_TOKEN")
	}
}
