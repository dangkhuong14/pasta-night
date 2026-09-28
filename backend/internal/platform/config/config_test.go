package config

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func envOf(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

func TestFromEnvDefaults(t *testing.T) {
	cfg, err := FromEnv(envOf(map[string]string{"TMDB_READ_TOKEN": "token"}))
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	want := Config{
		Port:                 8080,
		TMDBReadToken:        "token",
		TMDBLanguage:         "vi-VN",
		WatchRegion:          "VN",
		CacheDir:             "./data/cache",
		ListTTL:              168 * time.Hour,
		DetailTTL:            336 * time.Hour,
		RefreshCheckInterval: time.Hour,
		DiscoverPages:        2,
		TMDBConcurrency:      5,
		CORSAllowedOrigins:   []string{},
	}
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("FromEnv =\n%+v\nwant\n%+v", cfg, want)
	}
	if cfg.IsAdminEnabled() {
		t.Error("admin endpoints must be disabled when ADMIN_TOKEN is not set")
	}
}

func TestFromEnvOverrides(t *testing.T) {
	cfg, err := FromEnv(envOf(map[string]string{
		"PORT":                   "9000",
		"TMDB_READ_TOKEN":        "  token  ",
		"TMDB_LANGUAGE":          "en-US",
		"WATCH_REGION":           "US",
		"CACHE_DIR":              "/data/cache",
		"LIST_TTL":               "24h",
		"DETAIL_TTL":             "48h",
		"REFRESH_CHECK_INTERVAL": "30m",
		"DISCOVER_PAGES":         "3",
		"TMDB_CONCURRENCY":       "8",
		"CORS_ALLOWED_ORIGINS":   "http://localhost:3000, https://Pasta.Example.com,",
		"ADMIN_TOKEN":            "admin-token",
	}))
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	want := Config{
		Port:                 9000,
		TMDBReadToken:        "token",
		TMDBLanguage:         "en-US",
		WatchRegion:          "US",
		CacheDir:             "/data/cache",
		ListTTL:              24 * time.Hour,
		DetailTTL:            48 * time.Hour,
		RefreshCheckInterval: 30 * time.Minute,
		DiscoverPages:        3,
		TMDBConcurrency:      8,
		CORSAllowedOrigins:   []string{"http://localhost:3000", "https://pasta.example.com"},
		AdminToken:           "admin-token",
	}
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("FromEnv =\n%+v\nwant\n%+v", cfg, want)
	}
	if !cfg.IsAdminEnabled() {
		t.Error("admin endpoints must be enabled when ADMIN_TOKEN is set")
	}
}

func TestFromEnvErrors(t *testing.T) {
	withToken := func(extra map[string]string) map[string]string {
		vars := map[string]string{"TMDB_READ_TOKEN": "token"}
		for k, v := range extra {
			vars[k] = v
		}
		return vars
	}
	tests := []struct {
		name    string
		env     map[string]string
		wantErr string
	}{
		{"missing tmdb token", map[string]string{}, "TMDB_READ_TOKEN: is required"},
		{"blank tmdb token", map[string]string{"TMDB_READ_TOKEN": "   "}, "TMDB_READ_TOKEN: is required"},
		{"port not a number", withToken(map[string]string{"PORT": "http"}), "PORT"},
		{"port out of range", withToken(map[string]string{"PORT": "70000"}), "PORT"},
		{"duration in days", withToken(map[string]string{"LIST_TTL": "7d"}), "LIST_TTL"},
		{"zero duration", withToken(map[string]string{"REFRESH_CHECK_INTERVAL": "0s"}), "REFRESH_CHECK_INTERVAL"},
		{"negative duration", withToken(map[string]string{"DETAIL_TTL": "-1h"}), "DETAIL_TTL"},
		{"zero pages", withToken(map[string]string{"DISCOVER_PAGES": "0"}), "DISCOVER_PAGES"},
		{"negative concurrency", withToken(map[string]string{"TMDB_CONCURRENCY": "-1"}), "TMDB_CONCURRENCY"},
		{"lowercase region", withToken(map[string]string{"WATCH_REGION": "vn"}), "WATCH_REGION"},
		{"wildcard origin", withToken(map[string]string{"CORS_ALLOWED_ORIGINS": "*"}), "CORS_ALLOWED_ORIGINS"},
		{"origin with slash", withToken(map[string]string{"CORS_ALLOWED_ORIGINS": "http://localhost:3000/"}), "CORS_ALLOWED_ORIGINS"},
		{"origin without scheme", withToken(map[string]string{"CORS_ALLOWED_ORIGINS": "localhost:3000"}), "CORS_ALLOWED_ORIGINS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromEnv(envOf(tt.env))
			if err == nil {
				t.Fatalf("FromEnv succeeded, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q does not contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestFromEnvReportsEveryError(t *testing.T) {
	_, err := FromEnv(envOf(map[string]string{"PORT": "x", "LIST_TTL": "y"}))
	if err == nil {
		t.Fatal("want error")
	}
	for _, key := range []string{"TMDB_READ_TOKEN", "PORT", "LIST_TTL"} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("error %q does not mention %s", err, key)
		}
	}
}

func TestFromEnvNeverEchoesSecrets(t *testing.T) {
	_, err := FromEnv(envOf(map[string]string{
		"TMDB_READ_TOKEN": "tmdb-secret-value",
		"ADMIN_TOKEN":     "admin-secret-value",
		"PORT":            "not-a-port",
	}))
	if err == nil {
		t.Fatal("want error")
	}
	for _, secret := range []string{"tmdb-secret-value", "admin-secret-value"} {
		if strings.Contains(err.Error(), secret) {
			t.Errorf("error %q leaks a secret", err)
		}
	}
}
