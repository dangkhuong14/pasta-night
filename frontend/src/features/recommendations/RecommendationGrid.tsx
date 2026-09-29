"use client";

import { useRef } from "react";

import { PrimaryButton } from "@/components/common/PrimaryButton";
import { StickyActionBar } from "@/components/common/StickyActionBar";
import type { MovieSummary } from "@/lib/api-types";

import { DISPLAY_COUNT } from "./constants";
import { MovieCard } from "./MovieCard";
import { useShuffle } from "./useShuffle";

type RecommendationGridProps = {
  /** The full list (≤ 40), already shuffled on the server for this request. */
  movies: MovieSummary[];
};

/** 2-column grid of 6 movies + "GỢI Ý KHÁC" (SCREENS §2). */
export function RecommendationGrid({ movies }: RecommendationGridProps) {
  const { visible, round, shouldShowNextButton, showNext } = useShuffle(
    movies,
    DISPLAY_COUNT,
  );
  const gridRef = useRef<HTMLDivElement>(null);

  function handleShowNext() {
    showNext();
    const isReducedMotion = window.matchMedia(
      "(prefers-reduced-motion: reduce)",
    ).matches;
    gridRef.current?.scrollIntoView({
      behavior: isReducedMotion ? "auto" : "smooth",
      block: "start",
    });
  }

  return (
    <>
      <div
        ref={gridRef}
        className="grid scroll-mt-20 grid-cols-2 gap-4 px-5 pt-6"
      >
        {visible.map((movie, i) => (
          <MovieCard
            key={`${round}-${movie.id}`}
            movie={movie}
            isPriority={round === 0 && i < 2}
            animationDelayMs={i * 40}
          />
        ))}
      </div>
      {shouldShowNextButton && (
        <StickyActionBar>
          <PrimaryButton onClick={handleShowNext}>Gợi ý khác</PrimaryButton>
        </StickyActionBar>
      )}
    </>
  );
}
