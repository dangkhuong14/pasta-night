// Mirrors docs/API_SPEC.md v1. Field names are the API's snake_case on purpose,
// so this file can be diffed against the spec. Nullability follows DATABASE.md §2.

/** GET /options → data[] (API_SPEC §5.1). */
export type ViewingOption = {
  id: string;
  label: string;
  description: string;
  icon: string;
};

/** Watch provider in region VN (DATABASE.md §2 Provider). */
export type Provider = {
  id: number;
  name: string;
  logo_url: string | null;
  type: "flatrate" | "free" | "ads" | "rent" | "buy";
};

/** GET /options/{option_id}/recommendations → data[] (API_SPEC §5.2). */
export type MovieSummary = {
  id: number;
  title: string;
  overview: string;
  poster_url: string | null;
  release_year: number | null;
  rating: number;
  runtime_minutes: number | null;
  genres: string[];
  providers: Provider[];
};

/** One billed actor; profile_url is null when TMDB has no photo (DATABASE.md §2). */
export type CastMember = {
  name: string;
  profile_url: string | null;
};

/**
 * One slide of the detail carousel (DATABASE.md §2). The video, when the
 * movie has one, always comes first; `youtube_key` is set only for it.
 */
export type MediaItem = {
  type: "video" | "image";
  url: string;
  youtube_key: string | null;
};

/** GET /movies/{movie_id} → data (API_SPEC §5.3). */
export type MovieDetail = MovieSummary & {
  original_title: string;
  tagline: string;
  backdrop_url: string | null;
  vote_count: number;
  directors: string[];
  cast: CastMember[];
  media: MediaItem[];
};

export type RecommendationsMeta = {
  option_id: string;
  total: number;
  fetched_at: string;
  stale: boolean;
};

export type MovieMeta = {
  fetched_at: string;
};

/** Empty meta object `{}`. */
export type EmptyMeta = Record<string, never>;

/** Success envelope (API_SPEC §4). */
export type Envelope<TData, TMeta> = {
  data: TData;
  meta: TMeta;
};

/** Error body (API_SPEC §4). */
export type ErrorBody = {
  code: string;
  message: string;
  details: Record<string, unknown> | null;
};

/** API error codes (API_SPEC §6) plus NETWORK_ERROR for requests that never got an answer. */
export type ApiErrorCode =
  | "VALIDATION_ERROR"
  | "UNAUTHORIZED"
  | "NOT_FOUND"
  | "OPTION_NOT_FOUND"
  | "MOVIE_NOT_FOUND"
  | "METHOD_NOT_ALLOWED"
  | "CACHE_NOT_READY"
  | "INTERNAL_ERROR"
  | "NETWORK_ERROR";
