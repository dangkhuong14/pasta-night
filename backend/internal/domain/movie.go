package domain

import "time"

// MaxListSize is the maximum number of movies in one option's list (DATABASE.md §2).
const MaxListSize = 40

// MovieList is the ranked list of movie IDs for one option. The JSON tags
// define the cache file format; every ID has a MovieDetail in the same snapshot.
type MovieList struct {
	OptionID  string    `json:"option_id"`
	FetchedAt time.Time `json:"fetched_at"`
	MovieIDs  []int     `json:"movie_ids"`
}

// MovieDetail is everything the API can return about one movie. The JSON tags
// define the cache file format; pointer fields are nullable and encode as null.
type MovieDetail struct {
	ID             int          `json:"id"`
	Title          string       `json:"title"`
	OriginalTitle  string       `json:"original_title"`
	Overview       string       `json:"overview"`
	Tagline        string       `json:"tagline"`
	PosterURL      *string      `json:"poster_url"`
	BackdropURL    *string      `json:"backdrop_url"`
	ReleaseYear    *int         `json:"release_year"`
	Rating         float64      `json:"rating"`
	VoteCount      int          `json:"vote_count"`
	RuntimeMinutes *int         `json:"runtime_minutes"`
	Genres         []string     `json:"genres"`
	Directors      []string     `json:"directors"`
	Cast           []CastMember `json:"cast"`
	TrailerURL     *string      `json:"trailer_url"`
	Providers      []Provider   `json:"providers"`
	FetchedAt      time.Time    `json:"fetched_at"`
}

// CastMember is one billed actor. ProfileURL is nil when TMDB has no photo,
// and the UI falls back to the actor's initials.
type CastMember struct {
	Name       string  `json:"name"`
	ProfileURL *string `json:"profile_url"`
}

// ProviderType is how a watch provider offers a movie.
type ProviderType string

// Provider types, in the order the refresh mapper walks them.
const (
	ProviderFlatrate ProviderType = "flatrate"
	ProviderFree     ProviderType = "free"
	ProviderAds      ProviderType = "ads"
	ProviderRent     ProviderType = "rent"
	ProviderBuy      ProviderType = "buy"
)

// Provider is a streaming or rental service that offers a movie somewhere in
// the world (TMDB_INTEGRATION.md §5 provider rules).
type Provider struct {
	ID      int          `json:"id"`
	Name    string       `json:"name"`
	LogoURL *string      `json:"logo_url"`
	Type    ProviderType `json:"type"`
}
