// The only module that calls the backend (PROJECT-RULES §3 Data fetching).
// Contract: docs/API_SPEC.md v1.
import type {
  EmptyMeta,
  Envelope,
  ErrorBody,
  MovieDetail,
  MovieMeta,
  MovieSummary,
  RecommendationsMeta,
  ViewingOption,
} from "./api-types";

const API_BASE_URL = (process.env.NEXT_PUBLIC_API_BASE_URL ?? "").replace(
  /\/+$/,
  "",
);

/** Server-side cache lifetime for GETs, matching the API's `Cache-Control: max-age=300`. */
const REVALIDATE_SECONDS = 300;
/** Fail fast instead of hanging a page render on a stuck backend. */
const REQUEST_TIMEOUT_MS = 10_000;
/** API_SPEC §6: CACHE_NOT_READY sends `Retry-After: 30`; used when the header is missing. */
export const DEFAULT_RETRY_AFTER_SECONDS = 30;

/** A non-2xx API response, or `NETWORK_ERROR` (status 0) when no response arrived. */
export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string, // API error code, e.g. "OPTION_NOT_FOUND"
    public details: unknown,
    public retryAfter?: number, // seconds, from Retry-After
  ) {
    super(code);
    this.name = "ApiError";
  }
}

/** True when no backend is configured: every call is answered from `mock-data.ts`. */
export function isMockMode(): boolean {
  return API_BASE_URL === "";
}

async function request<TData, TMeta>(
  path: string,
  init?: RequestInit,
): Promise<Envelope<TData, TMeta>> {
  let res: Response;
  try {
    res = await fetch(`${API_BASE_URL}${path}`, {
      ...init,
      headers: { Accept: "application/json" },
      signal: AbortSignal.timeout(REQUEST_TIMEOUT_MS),
    });
  } catch {
    throw new ApiError(0, "NETWORK_ERROR", null);
  }

  if (!res.ok) {
    const body = (await res
      .json()
      .catch(() => null)) as Partial<ErrorBody> | null;
    const code = typeof body?.code === "string" ? body.code : "INTERNAL_ERROR";
    throw new ApiError(
      res.status,
      code,
      body?.details ?? null,
      parseRetryAfter(res.headers.get("Retry-After")),
    );
  }
  return (await res.json()) as Envelope<TData, TMeta>;
}

function parseRetryAfter(value: string | null): number | undefined {
  if (value === null) return undefined;
  const seconds = Number.parseInt(value, 10);
  return Number.isFinite(seconds) && seconds >= 0 ? seconds : undefined;
}

// Loaded lazily so mock data never ships in bundles built against a real API.
const loadMock = () => import("./mock-data");

/** GET /options (API_SPEC §5.1). */
export async function getOptions(): Promise<
  Envelope<ViewingOption[], EmptyMeta>
> {
  if (isMockMode()) {
    const { MOCK_OPTIONS } = await loadMock();
    return { data: MOCK_OPTIONS, meta: {} };
  }
  return request<ViewingOption[], EmptyMeta>("/options", {
    next: { revalidate: REVALIDATE_SECONDS },
  });
}

/** GET /options/{option_id}/recommendations (API_SPEC §5.2). */
export async function getRecommendations(
  optionId: string,
): Promise<Envelope<MovieSummary[], RecommendationsMeta>> {
  if (isMockMode()) {
    const { MOCK_FETCHED_AT, MOCK_MOVIES, MOCK_OPTIONS } = await loadMock();
    if (!MOCK_OPTIONS.some((o) => o.id === optionId)) {
      throw new ApiError(404, "OPTION_NOT_FOUND", { option_id: optionId });
    }
    const data = MOCK_MOVIES.map(toSummary);
    return {
      data,
      meta: {
        option_id: optionId,
        total: data.length,
        fetched_at: MOCK_FETCHED_AT,
        stale: false,
      },
    };
  }
  return request<MovieSummary[], RecommendationsMeta>(
    `/options/${encodeURIComponent(optionId)}/recommendations`,
    {
      next: { revalidate: REVALIDATE_SECONDS },
    },
  );
}

/** GET /movies/{movie_id} (API_SPEC §5.3). Called from the browser when the detail sheet opens. */
export async function getMovie(
  movieId: number,
): Promise<Envelope<MovieDetail, MovieMeta>> {
  if (isMockMode()) {
    const { MOCK_FETCHED_AT, MOCK_MOVIES } = await loadMock();
    const movie = MOCK_MOVIES.find((m) => m.id === movieId);
    if (!movie)
      throw new ApiError(404, "MOVIE_NOT_FOUND", { movie_id: movieId });
    return { data: movie, meta: { fetched_at: MOCK_FETCHED_AT } };
  }
  return request<MovieDetail, MovieMeta>(`/movies/${movieId}`);
}

/** The API's MovieSummary is a subset of MovieDetail (API_SPEC §5.2). */
function toSummary(m: MovieDetail): MovieSummary {
  return {
    id: m.id,
    title: m.title,
    overview: m.overview,
    poster_url: m.poster_url,
    release_year: m.release_year,
    rating: m.rating,
    runtime_minutes: m.runtime_minutes,
    genres: m.genres,
    providers: m.providers,
  };
}
