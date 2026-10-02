package tmdb

// DTOs decode only the fields we use (TMDB_INTEGRATION.md §4). JSON null and
// missing fields both decode to zero values; the refresh mapper turns zero
// values into null where the domain allows it. DTOs never leave
// platform/tmdb and refresh.

// DiscoverResponse is one page of GET /discover/movie.
type DiscoverResponse struct {
	Page       int             `json:"page"`
	TotalPages int             `json:"total_pages"`
	Results    []DiscoverMovie `json:"results"`
}

// DiscoverMovie is one discover result; phase 1 needs only IDs.
type DiscoverMovie struct {
	ID    int  `json:"id"`
	Adult bool `json:"adult"`
}

// MovieDetails is GET /movie/{id} with credits, videos, and watch/providers appended.
type MovieDetails struct {
	ID             int            `json:"id"`
	Adult          bool           `json:"adult"`
	Title          string         `json:"title"`
	OriginalTitle  string         `json:"original_title"`
	Overview       string         `json:"overview"`
	Tagline        string         `json:"tagline"`
	PosterPath     string         `json:"poster_path"`
	BackdropPath   string         `json:"backdrop_path"`
	ReleaseDate    string         `json:"release_date"` // "YYYY-MM-DD" or ""
	VoteAverage    float64        `json:"vote_average"`
	VoteCount      int            `json:"vote_count"`
	Runtime        int            `json:"runtime"`
	Genres         []Genre        `json:"genres"`
	Credits        Credits        `json:"credits"`
	Videos         Videos         `json:"videos"`
	Images         Images         `json:"images"`
	WatchProviders WatchProviders `json:"watch/providers"` // the JSON key contains a slash
}

// Images is the appended images object. Only backdrops are used: they are
// 16:9, which is the shape of the detail sheet's carousel.
type Images struct {
	Backdrops []Image `json:"backdrops"`
}

// Image is one artwork file. Language is "" for textless art, which is what
// the carousel prefers (TMDB_INTEGRATION.md §5).
type Image struct {
	FilePath    string  `json:"file_path"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	AspectRatio float64 `json:"aspect_ratio"`
	VoteAverage float64 `json:"vote_average"`
	Language    string  `json:"iso_639_1"`
}

// Genre is a localized TMDB genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Credits is the appended credits object.
type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

// CastMember is one actor.
type CastMember struct {
	Name        string `json:"name"`
	Order       int    `json:"order"`        // billing order, 0 = top
	ProfilePath string `json:"profile_path"` // "" when TMDB has no photo
}

// CrewMember is one crew member.
type CrewMember struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

// Videos is the appended videos object.
type Videos struct {
	Results []Video `json:"results"`
}

// Video is one video; trailers are those with Site "YouTube" and Type "Trailer".
type Video struct {
	Key      string `json:"key"`
	Site     string `json:"site"` // "YouTube", "Vimeo"
	Type     string `json:"type"` // "Trailer", "Teaser", "Clip", ...
	Official bool   `json:"official"`
	Language string `json:"iso_639_1"` // "vi", "en"
}

// WatchProviders is the appended watch/providers object.
type WatchProviders struct {
	Results map[string]RegionProviders `json:"results"` // key = region code, e.g. "VN"
}

// RegionProviders lists the providers of one region by monetization type.
type RegionProviders struct {
	Flatrate []Provider `json:"flatrate"`
	Free     []Provider `json:"free"`
	Ads      []Provider `json:"ads"`
	Rent     []Provider `json:"rent"`
	Buy      []Provider `json:"buy"`
}

// Provider is one watch provider entry.
type Provider struct {
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name"`
	LogoPath        string `json:"logo_path"`
	DisplayPriority int    `json:"display_priority"`
}
