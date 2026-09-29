package refresh

import (
	"cmp"
	"math"
	"net/url"
	"slices"
	"strconv"
	"strings"
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
	profileSize  = "w185" // cast avatars render at 48 px, so this covers 3x screens
	youtubeURL   = "https://www.youtube.com/watch?v="
	maxCast      = 5
	// maxProviders caps the merged list: popular movies are offered by 40+
	// services worldwide, most of them single-country (TMDB_INTEGRATION.md §5).
	maxProviders = 8
)

// mapMovieDetails converts a TMDB movie into a domain.MovieDetail following
// TMDB_INTEGRATION.md §5. ok is false when the movie must be dropped.
func mapMovieDetails(m tmdb.MovieDetails, fetchedAt time.Time) (domain.MovieDetail, bool) {
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
		Providers:      providers(m.WatchProviders),
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

// topCast returns the first maxCast actors in billing order.
func topCast(cast []tmdb.CastMember) []domain.CastMember {
	sorted := slices.Clone(cast)
	slices.SortStableFunc(sorted, func(a, b tmdb.CastMember) int { return cmp.Compare(a.Order, b.Order) })
	members := make([]domain.CastMember, 0, maxCast)
	for _, c := range sorted {
		if len(members) == maxCast {
			break
		}
		if c.Name != "" {
			members = append(members, domain.CastMember{
				Name:       c.Name,
				ProfileURL: imageURL(profileSize, c.ProfilePath),
			})
		}
	}
	return members
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

// providerTypes are the monetization types in preference order: a provider
// offered as flatrate in any region outranks one that is only rentable.
var providerTypes = []domain.ProviderType{
	domain.ProviderFlatrate,
	domain.ProviderFree,
	domain.ProviderAds,
	domain.ProviderRent,
	domain.ProviderBuy,
}

func regionEntries(rp tmdb.RegionProviders) [][]tmdb.Provider {
	return [][]tmdb.Provider{rp.Flatrate, rp.Free, rp.Ads, rp.Rent, rp.Buy}
}

// providerAgg accumulates one provider across every region that offers it.
type providerAgg struct {
	provider domain.Provider
	typeRank int // index into providerTypes; lowest wins
	priority int // lowest display_priority seen
	regions  int // how many regions offer it
}

// providers merges the watch providers of every region, because TMDB has no
// data at all for some regions (VN included) and per-region lists are tiny.
// Providers are ranked by how many regions offer them, which surfaces the
// global platforms customers recognize and drops single-country services
// (TMDB_INTEGRATION.md §5).
func providers(wp tmdb.WatchProviders) []domain.Provider {
	byID := make(map[int]*providerAgg, len(wp.Results))
	// Regions are walked in a fixed order so the output never depends on Go's
	// random map iteration: the same movie must always map to the same list.
	regions := make([]string, 0, len(wp.Results))
	for region := range wp.Results {
		regions = append(regions, region)
	}
	slices.Sort(regions)

	for _, region := range regions {
		countedInRegion := make(map[int]bool)
		for typeRank, entries := range regionEntries(wp.Results[region]) {
			for _, p := range entries {
				if p.ProviderID <= 0 || p.ProviderName == "" {
					continue
				}
				agg, ok := byID[p.ProviderID]
				if !ok {
					agg = &providerAgg{
						provider: domain.Provider{
							ID:      p.ProviderID,
							Name:    p.ProviderName,
							LogoURL: imageURL(logoSize, p.LogoPath),
							Type:    providerTypes[typeRank],
						},
						typeRank: typeRank,
						priority: p.DisplayPriority,
					}
					byID[p.ProviderID] = agg
				}
				if typeRank < agg.typeRank {
					agg.typeRank = typeRank
					agg.provider.Type = providerTypes[typeRank]
				}
				agg.priority = min(agg.priority, p.DisplayPriority)
				if agg.provider.LogoURL == nil {
					agg.provider.LogoURL = imageURL(logoSize, p.LogoPath)
				}
				if !countedInRegion[p.ProviderID] {
					countedInRegion[p.ProviderID] = true
					agg.regions++
				}
			}
		}
	}

	all := make([]*providerAgg, 0, len(byID))
	for _, agg := range byID {
		all = append(all, agg)
	}
	slices.SortFunc(all, func(a, b *providerAgg) int {
		return cmp.Or(
			cmp.Compare(b.regions, a.regions), // most widely available first
			cmp.Compare(a.typeRank, b.typeRank),
			cmp.Compare(a.priority, b.priority),
			cmp.Compare(a.provider.ID, b.provider.ID), // deterministic tie-break
		)
	})
	// TMDB gives one service several IDs (9 and 119 are both "Amazon Prime
	// Video"), which would show the same name twice on the card. The list is
	// already ranked, so the first entry of a name is the best one to keep.
	out := make([]domain.Provider, 0, maxProviders)
	seenNames := make(map[string]bool, len(all))
	for _, agg := range all {
		if len(out) == maxProviders {
			break
		}
		name := strings.ToLower(strings.TrimSpace(agg.provider.Name))
		if seenNames[name] {
			continue
		}
		seenNames[name] = true
		out = append(out, agg.provider)
	}
	return out
}
