package refresh

// allowedProviderIDs are the watch providers a customer in Vietnam can
// realistically use, by TMDB provider ID. Everything else is dropped while
// mapping (TMDB_INTEGRATION.md §5 provider rules).
//
// Why an allowlist at all: TMDB has no watch-provider data for region VN, so
// providers are merged across every region. Without this filter a movie lists
// services that do not operate here at all (Rakuten TV, Viaplay, Sky Store,
// SF Anytime, maxdome…), which is worse than showing nothing.
//
// These services operate in Vietnam but TMDB does not know them, so they can
// never appear: FPT Play, VieON, Galaxy Play, TV360, K+, Danet, Bilibili.
// Re-check when TMDB adds region VN; the IDs come from
// GET /watch/providers/movie.
var allowedProviderIDs = map[int]bool{
	8:   true, // Netflix
	9:   true, // Amazon Prime Video
	119: true, // Amazon Prime Video (second ID TMDB uses for the same service)
	2:   true, // Apple TV Store (rent and buy)
	350: true, // Apple TV (the subscription, formerly Apple TV+)
	192: true, // YouTube
	283: true, // Crunchyroll
	344: true, // Rakuten Viki
	581: true, // iQIYI
	623: true, // WeTV
}
