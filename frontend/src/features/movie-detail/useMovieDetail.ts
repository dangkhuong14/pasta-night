import { useEffect, useState } from "react";

import { ApiError, getMovie } from "@/lib/api";
import type { MovieDetail } from "@/lib/api-types";

export type MovieDetailState =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "loaded"; detail: MovieDetail }
  | { status: "error"; code: string };

type FetchResult = { movieId: number } & (
  { status: "loaded"; detail: MovieDetail } | { status: "error"; code: string }
);

// In-memory for the browser session (ARCHITECTURE §5): reopening a movie is instant.
const detailCache = new Map<number, MovieDetail>();

/** Fetches `GET /movies/{id}` when the sheet opens (SCREENS §3). */
export function useMovieDetail(movieId: number | null): MovieDetailState {
  const [result, setResult] = useState<FetchResult | null>(null);

  useEffect(() => {
    if (movieId === null || detailCache.has(movieId)) return;
    let isCancelled = false;
    getMovie(movieId)
      .then(({ data }) => {
        detailCache.set(movieId, data);
        if (!isCancelled)
          setResult({ movieId, status: "loaded", detail: data });
      })
      .catch((error: unknown) => {
        const code = error instanceof ApiError ? error.code : "INTERNAL_ERROR";
        if (!isCancelled) setResult({ movieId, status: "error", code });
      });
    return () => {
      isCancelled = true;
    };
  }, [movieId]);

  if (movieId === null) return { status: "idle" };
  const cached = detailCache.get(movieId);
  if (cached !== undefined) return { status: "loaded", detail: cached };
  if (result?.movieId === movieId) return result;
  return { status: "loading" };
}
