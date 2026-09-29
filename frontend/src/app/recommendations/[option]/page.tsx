import type { ReactNode } from "react";
import { redirect } from "next/navigation";
import { connection } from "next/server";

import { ErrorState } from "@/components/common/ErrorState";
import { PageHeader } from "@/components/common/PageHeader";
import { PrimaryButton } from "@/components/common/PrimaryButton";
import { SectionLabel } from "@/components/common/SectionLabel";
import { pastaPairings } from "@/config/brand";
import { MovieDetailSheet } from "@/features/movie-detail/MovieDetailSheet";
import type { PastaPairingView } from "@/features/movie-detail/PastaPairingCard";
import { CacheNotReady } from "@/features/recommendations/CacheNotReady";
import { RecommendationGrid } from "@/features/recommendations/RecommendationGrid";
import {
  ApiError,
  DEFAULT_RETRY_AFTER_SECONDS,
  getRecommendations,
} from "@/lib/api";
import type { MovieSummary } from "@/lib/api-types";
import { publicAsset } from "@/lib/assets";
import { errorCopy } from "@/lib/error-copy";
import { shuffle } from "@/lib/shuffle";

/** Recommendations for one option + the detail sheet (SCREENS §2, §3; ARCHITECTURE §4). */
export default async function RecommendationsPage({
  params,
}: PageProps<"/recommendations/[option]">) {
  // Per-request render: each visit gets its own shuffle; the API response itself is cached for 5 minutes.
  await connection();
  const { option } = await params;

  let movies: MovieSummary[];
  try {
    ({ data: movies } = await getRecommendations(option));
  } catch (error) {
    if (error instanceof ApiError && error.code === "OPTION_NOT_FOUND")
      redirect("/");
    if (error instanceof ApiError && error.code === "CACHE_NOT_READY") {
      return (
        <Shell>
          <CacheNotReady
            retryAfterSeconds={error.retryAfter ?? DEFAULT_RETRY_AFTER_SECONDS}
          />
        </Shell>
      );
    }
    throw error; // → error.tsx
  }

  if (movies.length === 0) {
    const copy = errorCopy("NO_RECOMMENDATIONS");
    return (
      <Shell>
        <ErrorState
          title={copy.title}
          body={copy.body}
          action={<PrimaryButton href="/">Chọn lựa khác</PrimaryButton>}
        />
      </Shell>
    );
  }

  return (
    <Shell>
      <section className="px-5 pt-6">
        <h1 className="font-serif text-2xl text-primary">Dành riêng cho bạn</h1>
        <SectionLabel className="mt-1">
          Những bộ phim hoàn hảo cho tối nay
        </SectionLabel>
      </section>
      {/* The initial shuffle happens here, per request, so SSR and hydration agree (SCREENS §2). */}
      <RecommendationGrid movies={shuffle(movies)} />
      <MovieDetailSheet movies={movies} pairing={resolvePairing(option)} />
    </Shell>
  );
}

function Shell({ children }: { children: ReactNode }) {
  return (
    <>
      <PageHeader backHref="/" />
      {children}
    </>
  );
}

/** Pasta pairing for the option, with its dish photo if the asset exists (SCREENS §3). */
function resolvePairing(optionId: string): PastaPairingView | null {
  const pairing = pastaPairings[optionId];
  if (pairing === undefined) return null;
  return {
    name: pairing.name,
    note: pairing.note,
    imageSrc: publicAsset(`images/dishes/${pairing.dishId}.webp`),
  };
}
