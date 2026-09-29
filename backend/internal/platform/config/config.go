// Package config loads and validates the environment and configs/options.yaml.
// Invalid configuration is reported at startup, never at request time.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Defaults for optional environment variables (ARCHITECTURE.md §6).
const (
	DefaultPort                 = 8080
	DefaultTMDBLanguage         = "vi-VN"
	DefaultCacheDir             = "./data/cache"
	DefaultListTTL              = 168 * time.Hour
	DefaultDetailTTL            = 336 * time.Hour
	DefaultRefreshCheckInterval = time.Hour
	DefaultDiscoverPages        = 2
	DefaultTMDBConcurrency      = 5
)

// Config is the validated process configuration.
type Config struct {
	Port                 int
	TMDBReadToken        string
	TMDBLanguage         string
	CacheDir             string
	ListTTL              time.Duration
	DetailTTL            time.Duration
	RefreshCheckInterval time.Duration
	DiscoverPages        int
	TMDBConcurrency      int
	CORSAllowedOrigins   []string
	// AdminToken guards the admin endpoints. Empty disables them (API_SPEC.md §2).
	AdminToken string
}

// IsAdminEnabled reports whether the admin endpoints are registered.
func (c Config) IsAdminEnabled() bool {
	return c.AdminToken != ""
}

// FromEnv reads the configuration through getenv (os.Getenv in production)
// and reports every invalid variable at once. Error messages never include
// secret values.
func FromEnv(getenv func(string) string) (Config, error) {
	p := envParser{getenv: getenv}
	cfg := Config{
		Port:                 p.intInRange("PORT", DefaultPort, 1, 65535),
		TMDBReadToken:        p.required("TMDB_READ_TOKEN"),
		TMDBLanguage:         p.str("TMDB_LANGUAGE", DefaultTMDBLanguage),
		CacheDir:             p.str("CACHE_DIR", DefaultCacheDir),
		ListTTL:              p.duration("LIST_TTL", DefaultListTTL),
		DetailTTL:            p.duration("DETAIL_TTL", DefaultDetailTTL),
		RefreshCheckInterval: p.duration("REFRESH_CHECK_INTERVAL", DefaultRefreshCheckInterval),
		DiscoverPages:        p.positiveInt("DISCOVER_PAGES", DefaultDiscoverPages),
		TMDBConcurrency:      p.positiveInt("TMDB_CONCURRENCY", DefaultTMDBConcurrency),
		CORSAllowedOrigins:   p.origins("CORS_ALLOWED_ORIGINS"),
		AdminToken:           p.str("ADMIN_TOKEN", ""),
	}
	if err := errors.Join(p.errs...); err != nil {
		return Config{}, fmt.Errorf("invalid environment: %w", err)
	}
	return cfg, nil
}

// envParser collects every problem instead of stopping at the first one.
type envParser struct {
	getenv func(string) string
	errs   []error
}

func (p *envParser) fail(key, msg string) {
	p.errs = append(p.errs, fmt.Errorf("%s: %s", key, msg))
}

func (p *envParser) str(key, def string) string {
	if v := strings.TrimSpace(p.getenv(key)); v != "" {
		return v
	}
	return def
}

func (p *envParser) required(key string) string {
	v := p.str(key, "")
	if v == "" {
		p.fail(key, "is required")
	}
	return v
}

func (p *envParser) intInRange(key string, def, minVal, maxVal int) int {
	raw := p.str(key, "")
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < minVal || n > maxVal {
		p.fail(key, fmt.Sprintf("must be an integer between %d and %d, got %q", minVal, maxVal, raw))
		return def
	}
	return n
}

func (p *envParser) positiveInt(key string, def int) int {
	raw := p.str(key, "")
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		p.fail(key, fmt.Sprintf("must be a positive integer, got %q", raw))
		return def
	}
	return n
}

func (p *envParser) duration(key string, def time.Duration) time.Duration {
	raw := p.str(key, "")
	if raw == "" {
		return def
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		p.fail(key, fmt.Sprintf("must be a positive Go duration such as 168h, got %q", raw))
		return def
	}
	return d
}

func (p *envParser) origins(key string) []string {
	origins := []string{}
	for _, part := range strings.Split(p.str(key, ""), ",") {
		origin := strings.ToLower(strings.TrimSpace(part))
		if origin == "" {
			continue
		}
		if err := validateOrigin(origin); err != nil {
			p.fail(key, err.Error())
			continue
		}
		origins = append(origins, origin)
	}
	return origins
}

// validateOrigin accepts exactly what browsers send in the Origin header:
// scheme://host[:port] with no path.
func validateOrigin(origin string) error {
	if origin == "*" {
		return errors.New(`wildcard "*" is not allowed; list the frontend origins`)
	}
	u, err := url.Parse(origin)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" ||
		u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return fmt.Errorf("invalid origin %q: want scheme://host[:port] without a trailing slash", origin)
	}
	return nil
}
