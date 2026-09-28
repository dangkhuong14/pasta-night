package httpx

import (
	"errors"
	"net/http"

	"pasta_night/be/internal/domain"
)

// Transport-level errors. Domain errors live in internal/domain.
var (
	// ErrNotFound means no route matches the request, including disabled admin routes.
	ErrNotFound = errors.New("route not found")
	// ErrMethodNotAllowed means the path exists but not for this HTTP method.
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrUnauthorized means the admin token is missing or wrong.
	ErrUnauthorized = errors.New("unauthorized")
)

type apiError struct {
	status     int
	code       string
	message    string // default message, English, for developers
	retryAfter string // Retry-After header value, if any
}

// errorMap is the only place that maps errors to HTTP statuses and API error
// codes (API_SPEC.md §6). Add a test row for every new entry.
var errorMap = map[error]apiError{
	domain.ErrValidation: {
		status: http.StatusBadRequest, code: "VALIDATION_ERROR", message: "invalid request parameter",
	},
	ErrUnauthorized: {
		status: http.StatusUnauthorized, code: "UNAUTHORIZED", message: "admin token is missing or invalid",
	},
	ErrNotFound: {
		status: http.StatusNotFound, code: "NOT_FOUND", message: "route not found",
	},
	domain.ErrOptionNotFound: {
		status: http.StatusNotFound, code: "OPTION_NOT_FOUND", message: "option does not exist",
	},
	domain.ErrMovieNotFound: {
		status: http.StatusNotFound, code: "MOVIE_NOT_FOUND", message: "movie is not in any current recommendation list",
	},
	ErrMethodNotAllowed: {
		status: http.StatusMethodNotAllowed, code: "METHOD_NOT_ALLOWED", message: "method not allowed",
	},
	domain.ErrCacheNotReady: {
		status: http.StatusServiceUnavailable, code: "CACHE_NOT_READY",
		message: "recommendations are being prepared, retry later", retryAfter: "30",
	},
}

var internalError = apiError{
	status: http.StatusInternalServerError, code: "INTERNAL_ERROR", message: "internal server error",
}

// errorBody is the error response body (API_SPEC.md §4).
type errorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

// detailedError carries a client-safe message and details for WriteError.
type detailedError struct {
	err     error
	message string
	details map[string]any
}

func (e *detailedError) Error() string {
	if e.message == "" {
		return e.err.Error()
	}
	return e.message + ": " + e.err.Error()
}

func (e *detailedError) Unwrap() error { return e.err }

// WithDetails wraps err with a client-safe message and details. An empty
// message keeps the default message of the error code. Both are sent only
// when err maps to an API error code; a 500 never exposes them.
func WithDetails(err error, message string, details map[string]any) error {
	return &detailedError{err: err, message: message, details: details}
}

// WriteError writes the error response for err. Errors missing from errorMap
// become 500 INTERNAL_ERROR: the cause goes to the request log, never to the client.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	recordError(r, err)
	api, isMapped := lookup(err)
	body := errorBody{Code: api.code, Message: api.message}
	var de *detailedError
	if isMapped && errors.As(err, &de) {
		if de.message != "" {
			body.Message = de.message
		}
		body.Details = de.details
	}
	if api.retryAfter != "" {
		w.Header().Set("Retry-After", api.retryAfter)
	}
	WriteJSON(w, r, api.status, body)
}

func lookup(err error) (apiError, bool) {
	for target, api := range errorMap {
		if errors.Is(err, target) {
			return api, true
		}
	}
	return internalError, false
}
