// Package domain holds the types and sentinel errors shared by every layer.
// It imports nothing from this module.
package domain

// GenreMode controls how an option's genre IDs are combined in a discover query.
type GenreMode string

// Genre modes accepted by ViewingOption discover.genre_mode (DATABASE.md §2).
const (
	GenreModeAnd GenreMode = "and"
	GenreModeOr  GenreMode = "or"
)

// Option is a viewing option ("who are you watching with") loaded from
// configs/options.yaml. Its ID appears in frontend routes, so it never changes.
type Option struct {
	ID          string
	Label       string
	Description string
	Icon        string
	Discover    DiscoverParams
}

// DiscoverParams are the TMDB discover filters that build an option's movie list.
type DiscoverParams struct {
	GenreIDs         []int
	GenreMode        GenreMode
	MinVoteAverage   float64
	MinVoteCount     int
	WatchProviderIDs []int
	SortBy           string
}
