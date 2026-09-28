package tmdb

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinels for TMDB failures. Callers branch on them with errors.Is
// (TMDB_INTEGRATION.md §6).
var (
	// ErrUnauthorized means TMDB rejected the token (bad or suspended).
	ErrUnauthorized = errors.New("tmdb: unauthorized")
	// ErrNotFound means the resource does not exist, e.g. a movie removed from TMDB.
	ErrNotFound = errors.New("tmdb: not found")
	// ErrRateLimited means TMDB still answered 429 after the allowed retries.
	ErrRateLimited = errors.New("tmdb: rate limited")
	// ErrBadRequest means TMDB rejected our parameters: a bug on our side.
	ErrBadRequest = errors.New("tmdb: bad request")
	// ErrUnavailable means TMDB or the network is down, timed out, or sent an unreadable body.
	ErrUnavailable = errors.New("tmdb: unavailable")
)

// APIError is an error response from TMDB. It unwraps to one of the sentinels.
type APIError struct {
	HTTPStatus int
	Code       int    // TMDB status_code
	Message    string // TMDB status_message
}

func (e *APIError) Error() string {
	return fmt.Sprintf("tmdb: http %d, status_code %d: %s", e.HTTPStatus, e.Code, e.Message)
}

// Unwrap maps the HTTP status to a sentinel. Any 4xx not listed in
// TMDB_INTEGRATION.md §6 counts as our bug (ErrBadRequest).
func (e *APIError) Unwrap() error {
	switch {
	case e.HTTPStatus == http.StatusUnauthorized:
		return ErrUnauthorized
	case e.HTTPStatus == http.StatusNotFound:
		return ErrNotFound
	case e.HTTPStatus == http.StatusTooManyRequests:
		return ErrRateLimited
	case e.HTTPStatus >= http.StatusInternalServerError:
		return ErrUnavailable
	default:
		return ErrBadRequest
	}
}

// errorResponse is TMDB's error body.
type errorResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}
