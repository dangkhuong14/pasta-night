package refresh

import (
	"cmp"
	"math"
	"net/url"
	"slices"
	"strconv"
	"time"

	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/tmdb"
)

// Mapping constants (TMDB_INTEGRATION.md §2, §5).
const (
	imageBaseURL = "https://image.tmdb.org/t/p/"
	posterSize   = "w500"
	backdropSize = "w1280"
	logoSize     = "w92"
	youtubeURL   = "https://www.youtube.com/watch?v="
	maxCast      = 5
)

// mapMovieDetails converts a TMDB movie into a domain.MovieDetail following
// TMDB_INTEGRATION.md §5. ok is false when the movie must be dropped.
func mapMovieDetails(m tmdb.MovieDetails, region string, fetchedAt time.Time) (domain.MovieDetail, bool) {
	if m.ID <= 0 || m.Adult {
		return domain.MovieDetail{}, false
	}
	title := m.Title
	if title == "" {
		title = m.OriginalTitle
	}
	if title == "" {
		return domain.MovieDetail{}, false
	}
	return domain.MovieDetail{
		ID:             m.ID,
		Title:          title,
		OriginalTitle:  m.OriginalTitle,
		Overview:       m.Overview,
		Tagline:        m.Tagline,
		PosterURL:      imageURL(posterSize, m.PosterPath),
		BackdropURL:    imageURL(backdropSize, m.BackdropPath),
		ReleaseYear:    releaseYear(m.ReleaseDate),
		Rating:         math.Round(m.VoteAverage*10) / 10,
		VoteCount:      m.VoteCount,
		RuntimeMinutes: positive(m.Runtime),
		Genres:         genreNames(m.Genres),
		Directors:      directors(m.Credits.Crew),
		Cast:           topCast(m.Credits.Cast),
		TrailerURL:     trailerURL(m.Videos.Results),
		Providers:      providers(m.WatchProviders, region),
		FetchedAt:      fetchedAt,
	}, true
}

// imageURL returns nil for an empty TMDB image path.
func imageURL(size, path string) *string {
	if path == "" {
		return nil
	}
	u := imageBaseURL + size + path
	return &u
}

// releaseYear reads the year of a "YYYY-MM-DD" date; "" or invalid → nil.
func releaseYear(date string) *int {
	if len(date) < 4 {
		return nil
	}
	year, err := strconv.Atoi(date[:4])
	if err != nil || year <= 0 {
		return nil
	}
	return &year
}

// positive returns nil for TMDB's "unknown" zero.
func positive(n int) *int {
	if n <= 0 {
		return nil
	}
	return &n
}

func genreNames(genres []tmdb.Genre) []string {
	names := make([]string, 0, len(genres))
	for _, g := range genres {
		if g.Name != "" {
			names = append(names, g.Name)
		}
	}
	return names
}

// directors keeps TMDB order and drops duplicate names (one person can be
// credited twice).
func directors(crew []tmdb.CrewMember) []string {
	names := []string{}
	seen := make(map[string]bool)
	for _, c := range crew {
		if c.Job != "Director" || c.Name == "" || seen[c.Name] {
			continue
		}
		seen[c.Name] = true
		names = append(names, c.Name)
	}
	return names
}

// topCast returns the first maxCast names in billing order.
func topCast(cast []tmdb.CastMember) []string {
	sorted := slices.Clone(cast)
	slices.SortStableFunc(sorted, func(a, b tmdb.CastMember) int { return cmp.Compare(a.Order, b.Order) })
	names := make([]string, 0, maxCast)
	for _, c := range sorted {
		if len(names) == maxCast {
			break
		}
		if c.Name != "" {
			names = append(names, c.Name)
		}
	}
	return names
}

// trailerURL picks the best YouTube trailer: language vi, then en, then
// others; official first within a language; TMDB order breaks ties.
func trailerURL(videos []tmdb.Video) *string {
	best := -1
	for i, v := range videos {
		if v.Site != "YouTube" || v.Type != "Trailer" || v.Key == "" {
			continue
		}
		if best == -1 || trailerRank(v) < trailerRank(videos[best]) {
			best = i
		}
	}
	if best == -1 {
		return nil
	}
	u := youtubeURL + url.QueryEscape(videos[best].Key)
	return &u
}

func trailerRank(v tmdb.Video) int {
	rank := 4
	switch v.Language {
	case "vi":
		rank = 0
	case "en":
		rank = 2
	}
	if !v.Official {
		rank++
	}
	return rank
}

// providers maps the WATCH_REGION providers: types walked in order, the
// first type seen wins for a provider, then a stable sort by display_priority.
func providers(wp tmdb.WatchProviders, region string) []domain.Provider {
	rp, ok := wp.Results[region]
	if !ok {
		return []domain.Provider{}
	}
	type ranked struct {
		provider domain.Provider
		priority int
	}
	groups := []struct {
		typ     domain.ProviderType
		entries []tmdb.Provider
	}{
		{domain.ProviderFlatrate, rp.Flatrate},
		{domain.ProviderFree, rp.Free},
		{domain.ProviderAds, rp.Ads},
		{domain.ProviderRent, rp.Rent},
		{domain.ProviderBuy, rp.Buy},
	}
	var all []ranked
	seen := make(map[int]bool)
	for _, g := range groups {
		for _, p := range g.entries {
			if p.ProviderID <= 0 || p.ProviderName == "" || seen[p.ProviderID] {
				continue
			}
			seen[p.ProviderID] = true
			all = append(all, ranked{
				provider: domain.Provider{
					ID:      p.ProviderID,
					Name:    p.ProviderName,
					LogoURL: imageURL(logoSize, p.LogoPath),
					Type:    g.typ,
				},
				priority: p.DisplayPriority,
			})
		}
	}
	slices.SortStableFunc(all, func(a, b ranked) int { return cmp.Compare(a.priority, b.priority) })
	out := make([]domain.Provider, 0, len(all))
	for _, r := range all {
		out = append(out, r.provider)
	}
	return out
}
