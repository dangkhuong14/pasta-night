package domain

import "errors"

// Sentinel errors shared across layers. platform/httpx maps them to API error codes.
var (
	// ErrOptionNotFound means the option ID is not configured in options.yaml.
	ErrOptionNotFound = errors.New("option not found")
	// ErrMovieNotFound means the movie is not in any current recommendation list.
	ErrMovieNotFound = errors.New("movie not found")
	// ErrCacheNotReady means the first refresh for an option has not finished yet.
	ErrCacheNotReady = errors.New("cache not ready")
	// ErrValidation means a path or query parameter is malformed.
	ErrValidation = errors.New("validation failed")
)
