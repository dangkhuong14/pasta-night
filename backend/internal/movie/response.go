package movie

import (
	"time"

	"pasta_night/be/internal/domain"
)

// API response shapes (API_SPEC.md §5.2, §5.3). They differ from the cache
// format: no schema_version, and fetched_at moves to meta.

type movieSummary struct {
	ID             int                `json:"id"`
	Title          string             `json:"title"`
	Overview       string             `json:"overview"`
	PosterURL      *string            `json:"poster_url"`
	ReleaseYear    *int               `json:"release_year"`
	Rating         float64            `json:"rating"`
	RuntimeMinutes *int               `json:"runtime_minutes"`
	Genres         []string           `json:"genres"`
	Providers      []providerResponse `json:"providers"`
}

type movieDetail struct {
	ID             int                `json:"id"`
	Title          string             `json:"title"`
	OriginalTitle  string             `json:"original_title"`
	Overview       string             `json:"overview"`
	Tagline        string             `json:"tagline"`
	PosterURL      *string            `json:"poster_url"`
	BackdropURL    *string            `json:"backdrop_url"`
	ReleaseYear    *int               `json:"release_year"`
	Rating         float64            `json:"rating"`
	VoteCount      int                `json:"vote_count"`
	RuntimeMinutes *int               `json:"runtime_minutes"`
	Genres         []string           `json:"genres"`
	Directors      []string           `json:"directors"`
	Cast           []string           `json:"cast"`
	TrailerURL     *string            `json:"trailer_url"`
	Providers      []providerResponse `json:"providers"`
}

type providerResponse struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	LogoURL *string `json:"logo_url"`
	Type    string  `json:"type"`
}

type recommendationsMeta struct {
	OptionID  string `json:"option_id"`
	Total     int    `json:"total"`
	FetchedAt string `json:"fetched_at"`
	IsStale   bool   `json:"stale"`
}

type detailMeta struct {
	FetchedAt string `json:"fetched_at"`
}

func newMovieSummaries(movies []domain.MovieDetail) []movieSummary {
	out := make([]movieSummary, 0, len(movies))
	for _, m := range movies {
		out = append(out, movieSummary{
			ID:             m.ID,
			Title:          m.Title,
			Overview:       m.Overview,
			PosterURL:      m.PosterURL,
			ReleaseYear:    m.ReleaseYear,
			Rating:         m.Rating,
			RuntimeMinutes: m.RuntimeMinutes,
			Genres:         nonNil(m.Genres),
			Providers:      newProviders(m.Providers),
		})
	}
	return out
}

func newMovieDetail(m domain.MovieDetail) movieDetail {
	return movieDetail{
		ID:             m.ID,
		Title:          m.Title,
		OriginalTitle:  m.OriginalTitle,
		Overview:       m.Overview,
		Tagline:        m.Tagline,
		PosterURL:      m.PosterURL,
		BackdropURL:    m.BackdropURL,
		ReleaseYear:    m.ReleaseYear,
		Rating:         m.Rating,
		VoteCount:      m.VoteCount,
		RuntimeMinutes: m.RuntimeMinutes,
		Genres:         nonNil(m.Genres),
		Directors:      nonNil(m.Directors),
		Cast:           nonNil(m.Cast),
		TrailerURL:     m.TrailerURL,
		Providers:      newProviders(m.Providers),
	}
}

func newProviders(providers []domain.Provider) []providerResponse {
	out := make([]providerResponse, 0, len(providers))
	for _, p := range providers {
		out = append(out, providerResponse{ID: p.ID, Name: p.Name, LogoURL: p.LogoURL, Type: string(p.Type)})
	}
	return out
}

func newRecommendationsMeta(recs Recommendations) recommendationsMeta {
	return recommendationsMeta{
		OptionID:  recs.OptionID,
		Total:     len(recs.Movies),
		FetchedAt: formatTime(recs.FetchedAt),
		IsStale:   recs.IsStale,
	}
}

// formatTime renders timestamps as RFC 3339 UTC (API_SPEC.md §3).
func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// nonNil keeps arrays as [] instead of null (API_SPEC.md §3).
func nonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
