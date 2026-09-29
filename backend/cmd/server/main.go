// Command server runs the pasta_night API. It serves recommendations from
// the in-memory cache and keeps the cache fresh from TMDB in the background.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"pasta_night/be/internal/cache"
	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/movie"
	"pasta_night/be/internal/option"
	"pasta_night/be/internal/platform/config"
	"pasta_night/be/internal/platform/httpx"
	"pasta_night/be/internal/platform/tmdb"
	"pasta_night/be/internal/refresh"
)

const (
	optionsPath     = "configs/options.yaml"
	shutdownTimeout = 10 * time.Second
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := run(ctx, os.Getenv, log)
	stop()
	if err != nil {
		log.Error("server stopped", "err", err)
		os.Exit(1)
	}
	log.Info("server stopped")
}

// run loads the configuration, wires the app, and serves until ctx is done.
func run(ctx context.Context, getenv func(string) string, log *slog.Logger) error {
	cfg, err := config.FromEnv(getenv)
	if err != nil {
		return err
	}
	options, err := config.LoadOptions(optionsPath)
	if err != nil {
		return fmt.Errorf("load %s: %w", optionsPath, err)
	}
	client := tmdb.NewClient(tmdb.DefaultBaseURL, cfg.TMDBReadToken, cfg.TMDBLanguage)
	a, err := newApp(ctx, cfg, options, client, log)
	if err != nil {
		return err
	}
	return a.serve(ctx)
}

// app is the wired server. Tests build it with a fake TMDB.
type app struct {
	cfg     config.Config
	options []domain.Option
	handler http.Handler
	store   *cache.Store
	job     *refresh.Job
	log     *slog.Logger
}

// newApp loads the cache and wires every layer (ARCHITECTURE.md §5 Startup).
func newApp(ctx context.Context, cfg config.Config, options []domain.Option, client refresh.TMDB, log *slog.Logger) (*app, error) {
	optionIDs := make([]string, len(options))
	for i, o := range options {
		optionIDs[i] = o.ID
	}
	store := cache.NewStore(cfg.CacheDir, optionIDs, log)
	stats, err := store.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("load cache: %w", err)
	}
	log.Info("cache loaded", "lists", stats.Lists, "details", stats.Details, "skipped", stats.Skipped)

	refresher := refresh.NewRefresher(client, store, refresh.Settings{
		DiscoverPages: cfg.DiscoverPages,
		Concurrency:   cfg.TMDBConcurrency,
		DetailTTL:     cfg.DetailTTL,
	}, log)
	job := refresh.NewJob(refresher, store, options, cfg.RefreshCheckInterval, cfg.ListTTL, log)

	mux := http.NewServeMux()
	option.NewHandler(option.NewService(options)).Register(mux)
	movieSvc := movie.NewService(movie.NewRepository(store), cfg.ListTTL+cfg.RefreshCheckInterval)
	movie.NewHandler(movieSvc, options).Register(mux)
	if cfg.IsAdminEnabled() {
		refresh.NewHandler(job, options).Register(mux, httpx.RequireAdminToken(cfg.AdminToken))
	} else {
		log.Info("admin endpoints disabled: ADMIN_TOKEN is not set")
	}
	mux.Handle("GET "+httpx.HealthzPath, httpx.Healthz(store.IsReady))

	handler := httpx.Chain(httpx.Fallback(mux),
		httpx.RequestID,
		httpx.LogRequests(log),
		httpx.Recover,
		httpx.CORS(cfg.CORSAllowedOrigins),
	)
	return &app{cfg: cfg, options: options, handler: handler, store: store, job: job, log: log}, nil
}

// serve runs the HTTP server and the refresh job until ctx is done or the
// server fails, then stops the job and drains requests for up to shutdownTimeout.
func (a *app) serve(ctx context.Context) error {
	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(a.cfg.Port),
		Handler:           a.handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ErrorLog:          slog.NewLogLogger(a.log.Handler(), slog.LevelWarn),
	}

	jobCtx, stopJob := context.WithCancel(ctx)
	defer stopJob()
	jobDone := make(chan struct{})
	go func() {
		defer close(jobDone)
		a.job.Run(jobCtx)
	}()

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.ListenAndServe() }()
	a.log.Info("server started",
		"port", a.cfg.Port,
		"cache_dir", a.cfg.CacheDir,
		"options", len(a.options),
		"admin_enabled", a.cfg.IsAdminEnabled(),
		"cors_origins", len(a.cfg.CORSAllowedOrigins),
	)

	var err error
	select {
	case <-ctx.Done():
		a.log.Info("shutting down")
	case err = <-serveErr: // ListenAndServe never returns nil
		err = fmt.Errorf("http server: %w", err)
	}
	stopJob()

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()
	if shutdownErr := srv.Shutdown(shutdownCtx); shutdownErr != nil {
		err = errors.Join(err, fmt.Errorf("shutdown http server: %w", shutdownErr))
	}
	<-jobDone
	return err
}
