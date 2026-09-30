"use client";

import {
  useCallback,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import Image from "next/image";
import { usePathname, useRouter, useSearchParams } from "next/navigation";

import { CalendarDays, Clock, X } from "lucide-react";

import { GenreChip } from "@/components/common/GenreChip";
import { RatingBadge } from "@/components/common/RatingBadge";
import { SectionLabel } from "@/components/common/SectionLabel";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerTitle,
} from "@/components/ui/drawer";
import type { MovieDetail, MovieSummary } from "@/lib/api-types";
import { formatRuntime } from "@/lib/format";

import { CastList } from "./CastList";
import { PastaPairingCard, type PastaPairingView } from "./PastaPairingCard";
import { PromoCard } from "./PromoCard";
import { ProvidersCard } from "./ProvidersCard";
import { SheetActions } from "./SheetActions";
import { useMovieDetail } from "./useMovieDetail";

const SHOP_URL = process.env.NEXT_PUBLIC_SHOP_URL ?? "";

type MovieDetailSheetProps = {
  /** The page's recommendation list: summaries render instantly from it. */
  movies: MovieSummary[];
  pairing: PastaPairingView | null;
  /**
   * The site footer. The sheet covers the whole screen, so the page's own
   * footer is hidden behind it; the route passes one in to close the page
   * off here too (SCREENS §3).
   */
  footer?: ReactNode;
};

/** `?movie=` must be a positive integer, else it is ignored (PROJECT-RULES §3). */
function parseMovieId(value: string | null): number | null {
  if (value === null || !/^[1-9]\d{0,9}$/.test(value)) return null;
  return Number(value);
}

/** Detail bottom sheet driven by `?movie={id}` (SCREENS §3). */
export function MovieDetailSheet({
  movies,
  pairing,
  footer,
}: MovieDetailSheetProps) {
  const router = useRouter();
  const pathname = usePathname();
  const movieId = parseMovieId(useSearchParams().get("movie"));
  const detailState = useMovieDetail(movieId);

  // Keep showing the last movie while the sheet animates closed.
  const [shownId, setShownId] = useState(movieId);
  if (movieId !== null && movieId !== shownId) setShownId(movieId);

  // A ?movie= present at first render came from a deep link: closing replaces
  // the URL instead of going back (which would leave the site).
  const deepLinkIdRef = useRef(movieId);

  const close = useCallback(() => {
    if (movieId !== null && movieId === deepLinkIdRef.current) {
      deepLinkIdRef.current = null;
      window.history.replaceState(null, "", pathname);
    } else {
      window.history.back();
    }
  }, [movieId, pathname]);

  const summary = movies.find((m) => m.id === shownId) ?? null;
  const detail = detailState.status === "loaded" ? detailState.detail : null;
  const isMissingFromList =
    movieId !== null && movies.every((m) => m.id !== movieId);

  useEffect(() => {
    if (detailState.status !== "error") return;
    if (detailState.code === "MOVIE_NOT_FOUND") {
      // The list changed under us: close and fetch the current list.
      close();
      router.refresh();
    } else if (isMissingFromList) {
      // Nothing to show without the detail: close silently.
      close();
    }
  }, [detailState, isMissingFromList, close, router]);

  const movie: MovieSummary | MovieDetail | null = detail ?? summary;
  const isOpen = movieId !== null && movie !== null;

  return (
    <Drawer
      open={isOpen}
      onOpenChange={(open) => {
        if (!open) close();
      }}
    >
      {/* Full height on mobile (DESIGN-SYSTEM §5). `!` beats the shadcn drawer's
          data-[swipe-axis=y] sizing, which caps the sheet at 100dvh - 6rem. */}
      <DrawerContent className="mx-auto h-dvh! max-h-dvh! w-full max-w-md rounded-none! border-t-0! bg-card">
        {movie !== null && (
          <DetailBody
            movie={movie}
            detail={detail?.id === movie.id ? detail : null}
            isDetailFailed={detailState.status === "error"}
            pairing={pairing}
            footer={footer}
          />
        )}
      </DrawerContent>
    </Drawer>
  );
}

type DetailBodyProps = {
  movie: MovieSummary | MovieDetail;
  detail: MovieDetail | null;
  isDetailFailed: boolean;
  pairing: PastaPairingView | null;
  footer?: ReactNode;
};

function DetailBody({
  movie,
  detail,
  isDetailFailed,
  pairing,
  footer,
}: DetailBodyProps) {
  // Poster first; the backdrop arrives with GET /movies/{id} (SCREENS §3 Data).
  const heroSrc = detail?.backdrop_url ?? movie.poster_url;
  const firstGenre = movie.genres[0];
  // null = still loading (skeleton); [] = hide the section.
  const cast = detail !== null ? detail.cast : isDetailFailed ? [] : null;

  return (
    <div className="h-full overflow-y-auto overscroll-contain">
      <div className="relative h-65 w-full shrink-0 bg-linear-to-b from-card to-background">
        {heroSrc !== null && (
          <Image
            src={heroSrc}
            alt=""
            fill
            sizes="(max-width: 448px) 100vw, 448px"
            priority
            className="object-cover"
          />
        )}
        <div
          aria-hidden
          className="absolute inset-0 bg-linear-to-t from-card via-card/10 to-transparent"
        />
        <DrawerClose
          aria-label="Đóng"
          className="absolute top-4 left-4 flex size-11 items-center justify-center rounded-full bg-black/60 text-primary backdrop-blur transition hover:text-gold-bright focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
        >
          <X aria-hidden className="size-5" strokeWidth={1.5} />
        </DrawerClose>
      </div>

      <div className="relative -mt-6 space-y-6 rounded-t-3xl bg-card px-5 pt-3 pb-10">
        <div aria-hidden className="mx-auto h-1 w-10 rounded-full bg-muted" />

        <div>
          <div className="flex items-start justify-between gap-3">
            <DrawerTitle className="font-serif text-[26px] leading-tight font-normal text-primary">
              {movie.title}
            </DrawerTitle>
            <RatingBadge
              rating={movie.rating}
              variant="outline"
              className="mt-1"
            />
          </div>
          <DrawerDescription className="sr-only">
            Thông tin chi tiết về phim {movie.title}
          </DrawerDescription>
          <MetaRow
            runtime={movie.runtime_minutes}
            year={movie.release_year}
            genre={firstGenre}
          />
        </div>

        {movie.providers.length > 0 && (
          <ProvidersCard providers={movie.providers} />
        )}
        {pairing !== null && <PastaPairingCard pairing={pairing} />}

        {movie.overview !== "" && (
          <section>
            <SectionLabel as="h3">Nội dung</SectionLabel>
            <p className="mt-2 text-sm leading-relaxed text-foreground/90">
              {movie.overview}
            </p>
          </section>
        )}

        {(cast === null || cast.length > 0) && <CastList cast={cast} />}
        {SHOP_URL !== "" && <PromoCard shopUrl={SHOP_URL} />}
        <SheetActions movieTitle={movie.title} />
      </div>
      {footer}
    </div>
  );
}

type MetaRowProps = {
  runtime: number | null;
  year: number | null;
  genre: string | undefined;
};

/** Clock + runtime · calendar + year · first genre; each hidden when missing (SCREENS §3). */
function MetaRow({ runtime, year, genre }: MetaRowProps) {
  if (runtime === null && year === null && genre === undefined) return null;
  return (
    <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-2 text-xs text-muted-foreground">
      {runtime !== null && (
        <span className="inline-flex items-center gap-1.5">
          <Clock aria-hidden className="size-3.5" strokeWidth={1.5} />
          {formatRuntime(runtime)}
        </span>
      )}
      {year !== null && (
        <span className="inline-flex items-center gap-1.5">
          <CalendarDays aria-hidden className="size-3.5" strokeWidth={1.5} />
          {year}
        </span>
      )}
      {genre !== undefined && <GenreChip label={genre} variant="neutral" />}
    </div>
  );
}
