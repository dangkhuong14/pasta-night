"use client";

import type { MouseEvent } from "react";
import Image from "next/image";

import { Film } from "lucide-react";

import { GenreChip } from "@/components/common/GenreChip";
import { RatingBadge } from "@/components/common/RatingBadge";
import type { MovieSummary } from "@/lib/api-types";

export type MovieCardProps = {
  movie: MovieSummary;
  /** Eager-load the first posters above the fold. */
  isPriority: boolean;
  /** Stagger for the enter animation (DESIGN-SYSTEM §4). */
  animationDelayMs: number;
};

/**
 * Poster card (DESIGN-SYSTEM §5 MovieCard). Tapping it adds `?movie={id}` with
 * `history.pushState`, which Next.js syncs into `useSearchParams`: the detail
 * sheet opens instantly without a server round trip, and back closes it.
 */
export function MovieCard({
  movie,
  isPriority,
  animationDelayMs,
}: MovieCardProps) {
  const href = `?movie=${movie.id}`;

  function handleClick(event: MouseEvent<HTMLAnchorElement>) {
    // Keep new-tab / new-window behavior for modified clicks.
    if (
      event.button !== 0 ||
      event.metaKey ||
      event.ctrlKey ||
      event.shiftKey ||
      event.altKey
    )
      return;
    event.preventDefault();
    window.history.pushState(null, "", href);
  }

  return (
    <a
      href={href}
      onClick={handleClick}
      style={{ animationDelay: `${animationDelayMs}ms` }}
      className="group block animate-fade-up rounded-lg focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-4 focus-visible:ring-offset-background focus-visible:outline-none motion-reduce:animate-none"
    >
      <div className="relative aspect-2/3 overflow-hidden rounded-lg border border-border bg-muted transition group-active:scale-[0.98] motion-reduce:transition-none">
        {movie.poster_url !== null ? (
          <Image
            src={movie.poster_url}
            alt={movie.title}
            fill
            sizes="(max-width: 448px) 50vw, 224px"
            priority={isPriority}
            className="object-cover"
          />
        ) : (
          <div className="flex h-full flex-col items-center justify-center gap-2 p-3 text-center">
            <Film
              aria-hidden
              className="size-8 text-muted-foreground"
              strokeWidth={1.5}
            />
            <span className="line-clamp-3 text-xs text-muted-foreground">
              {movie.title}
            </span>
          </div>
        )}
        <RatingBadge rating={movie.rating} className="absolute top-2 right-2" />
      </div>
      <h3 className="mt-2 truncate text-sm font-semibold text-foreground">
        {movie.title}
      </h3>
      {movie.genres.length > 0 && (
        <div className="mt-1.5 flex gap-1.5 overflow-hidden">
          {movie.genres.slice(0, 2).map((genre, i) => (
            <GenreChip
              key={genre}
              label={genre}
              variant={i === 0 ? "gold" : "neutral"}
            />
          ))}
        </div>
      )}
    </a>
  );
}
