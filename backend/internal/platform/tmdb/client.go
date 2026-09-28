// Package tmdb is the HTTP client for the two TMDB endpoints the refresh job
// uses. It knows TMDB and nothing about the domain (TMDB_INTEGRATION.md §3).
package tmdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DefaultBaseURL is the TMDB API v3 root.
const DefaultBaseURL = "https://api.themoviedb.org/3"

const (
	requestTimeout    = 10 * time.Second
	maxRetries        = 2 // for 429 only
	defaultRetryAfter = 5 * time.Second
	maxBodyBytes      = 10 << 20
)

// Client calls TMDB. It is safe for concurrent use; callers bound concurrency.
type Client struct {
	baseURL    string
	token      string
	language   string
	httpClient *http.Client
}

// NewClient returns a client for baseURL (DefaultBaseURL in production) that
// authenticates with token and sends language on every request.
func NewClient(baseURL, token, language string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		language:   language,
		httpClient: &http.Client{Timeout: requestTimeout},
	}
}

// DiscoverQuery holds the /discover/movie filters (TMDB_INTEGRATION.md §4.1).
type DiscoverQuery struct {
	Page               int
	WithGenres         string // genre IDs joined by "|" (or) or "," (and)
	VoteAverageGTE     float64
	VoteCountGTE       int
	SortBy             string
	WithWatchProviders string // provider IDs joined by "|"; empty = no provider filter
	WatchRegion        string // sent only with WithWatchProviders
}

// Discover fetches one page of /discover/movie.
func (c *Client) Discover(ctx context.Context, q DiscoverQuery) (DiscoverResponse, error) {
	params := url.Values{}
	params.Set("include_adult", "false")
	params.Set("include_video", "false")
	params.Set("with_genres", q.WithGenres)
	params.Set("vote_average.gte", strconv.FormatFloat(q.VoteAverageGTE, 'f', -1, 64))
	params.Set("vote_count.gte", strconv.Itoa(q.VoteCountGTE))
	params.Set("sort_by", q.SortBy)
	params.Set("page", strconv.Itoa(q.Page))
	if q.WithWatchProviders != "" {
		params.Set("with_watch_providers", q.WithWatchProviders)
		params.Set("watch_region", q.WatchRegion)
		params.Set("with_watch_monetization_types", "flatrate")
	}
	var resp DiscoverResponse
	if err := c.get(ctx, "/discover/movie", params, &resp); err != nil {
		return DiscoverResponse{}, fmt.Errorf("get /discover/movie page %d: %w", q.Page, err)
	}
	return resp, nil
}

// MovieDetails fetches /movie/{id} with credits, videos, and watch providers
// appended, in one call.
func (c *Client) MovieDetails(ctx context.Context, id int) (MovieDetails, error) {
	params := url.Values{}
	params.Set("append_to_response", "credits,videos,watch/providers")
	// With language=vi-VN alone most movies return no trailers.
	params.Set("include_video_language", "vi,en")
	var movie MovieDetails
	if err := c.get(ctx, "/movie/"+strconv.Itoa(id), params, &movie); err != nil {
		return MovieDetails{}, fmt.Errorf("get /movie/%d: %w", id, err)
	}
	return movie, nil
}

// get performs a GET and decodes the JSON body into out. Only 429 is
// retried, waiting Retry-After, at most maxRetries times.
func (c *Client) get(ctx context.Context, path string, params url.Values, out any) error {
	params.Set("language", c.language)
	rawURL := c.baseURL + path + "?" + params.Encode()
	for attempt := 0; ; attempt++ {
		wait, err := c.do(ctx, rawURL, out)
		if err == nil || !errors.Is(err, ErrRateLimited) || attempt == maxRetries {
			return err
		}
		if err := sleep(ctx, wait); err != nil {
			return err
		}
	}
}

// do performs one request. On 429 it also returns how long TMDB asked us to wait.
func (c *Client) do(ctx context.Context, rawURL string, out any) (time.Duration, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, transportError(ctx, err)
	}
	defer func() { _ = resp.Body.Close() }() // read-only body; close error is irrelevant

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return 0, transportError(ctx, fmt.Errorf("read body: %w", err))
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var wait time.Duration
		if resp.StatusCode == http.StatusTooManyRequests {
			wait = parseRetryAfter(resp.Header.Get("Retry-After"))
		}
		return wait, parseAPIError(resp.StatusCode, body)
	}
	if err := json.Unmarshal(body, out); err != nil {
		return 0, fmt.Errorf("%w: decode body: %w", ErrUnavailable, err)
	}
	return 0, nil
}

// transportError reports our own cancellation as-is, so callers can tell a
// shutdown from a TMDB outage, and everything else as ErrUnavailable.
func transportError(ctx context.Context, err error) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	return fmt.Errorf("%w: %w", ErrUnavailable, err)
}

func parseAPIError(status int, body []byte) *APIError {
	apiErr := &APIError{HTTPStatus: status}
	var parsed errorResponse
	if err := json.Unmarshal(body, &parsed); err == nil {
		apiErr.Code = parsed.StatusCode
		apiErr.Message = parsed.StatusMessage
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(status)
	}
	return apiErr
}

// parseRetryAfter reads Retry-After as seconds or an HTTP date, falling
// back to defaultRetryAfter.
func parseRetryAfter(v string) time.Duration {
	v = strings.TrimSpace(v)
	if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		return max(time.Until(t), 0)
	}
	return defaultRetryAfter
}

func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
