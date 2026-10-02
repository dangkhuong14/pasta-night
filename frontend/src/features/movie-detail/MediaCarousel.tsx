"use client";

import {
  useCallback,
  useEffect,
  useRef,
  useState,
  useSyncExternalStore,
} from "react";
import Image from "next/image";

import { Film } from "lucide-react";

import type { MediaItem } from "@/lib/api-types";
import { cn } from "@/lib/utils";

/** How long an image stays before the carousel moves on (SCREENS.md §3). */
const AUTO_ADVANCE_MS = 5000;
/** Quiet period after a touch or swipe before auto-advance resumes. */
const RESUME_AFTER_MS = 6000;
/** How long a smooth scroll of ours keeps firing scroll events. */
const SMOOTH_SCROLL_MS = 700;
/** Movement below this is a click, not a drag. */
const DRAG_SLOP_PX = 4;

/** A mouse drag in progress across the carousel. */
type Drag = {
  startX: number;
  startScroll: number;
  hasMoved: boolean;
};

const REDUCED_MOTION = "(prefers-reduced-motion: reduce)";

function subscribeReducedMotion(onChange: () => void) {
  const query = window.matchMedia(REDUCED_MOTION);
  query.addEventListener("change", onChange);
  return () => query.removeEventListener("change", onChange);
}

/**
 * The setting has to be read during render, to decide whether the trailer
 * autoplays at all. The server cannot know it, so it assumes motion is
 * welcome and hydration corrects it.
 */
function useReducedMotion() {
  return useSyncExternalStore(
    subscribeReducedMotion,
    () => window.matchMedia(REDUCED_MOTION).matches,
    () => false,
  );
}

/** Vietnamese label for one dot: "Trailer", or the still's number among stills. */
function slideLabel(slides: MediaItem[], index: number) {
  if (slides[index]?.type === "video") return "Trailer";
  const nth = slides
    .slice(0, index + 1)
    .filter((s) => s.type === "image").length;
  return `Ảnh ${nth}`;
}

type MediaCarouselProps = {
  /** Slides in order; the video, if any, is first. */
  media: MediaItem[];
  /** Shown while the detail request is still in flight, and if media is empty. */
  fallbackImageUrl: string | null;
  /** For the alt text of the fallback and of each still. */
  movieTitle: string;
};

/**
 * The detail sheet's hero: a horizontally scrolled carousel of the trailer
 * and stills (SCREENS.md §3). Images advance on their own; the video slide
 * does not, so a trailer is never cut off. Everything stands still for
 * customers who ask for reduced motion.
 */
export function MediaCarousel({
  media,
  fallbackImageUrl,
  movieTitle,
}: MediaCarouselProps) {
  const scrollerRef = useRef<HTMLDivElement>(null);
  const [current, setCurrent] = useState(0);
  // The video slide starts behind a transparent layer so a swipe over it
  // reaches the carousel; a tap hands the player over to the customer.
  const [isPlayerHandedOver, setIsPlayerHandedOver] = useState(false);
  const isReducedMotion = useReducedMotion();
  // Set while the customer is scrolling or has just touched the carousel.
  const pausedUntilRef = useRef(0);
  // Set while an auto-advance is scrolling, whose scroll events must not be
  // mistaken for the customer taking over.
  const selfScrollUntilRef = useRef(0);

  // The scroll position is the single source of truth for which slide shows,
  // so a swipe and an auto-advance cannot disagree.
  const handleScroll = useCallback(() => {
    const el = scrollerRef.current;
    if (el === null) return;
    if (Date.now() > selfScrollUntilRef.current) {
      pausedUntilRef.current = Date.now() + RESUME_AFTER_MS;
    }
    const index = Math.round(el.scrollLeft / el.clientWidth);
    setCurrent((prev) => (prev === index ? prev : index));
  }, []);

  const slides = media;

  const goTo = useCallback((index: number) => {
    const el = scrollerRef.current;
    if (el === null) return;
    // Tapping a dot is the customer steering: hold off auto-advance as a
    // swipe would.
    pausedUntilRef.current = Date.now() + RESUME_AFTER_MS;
    selfScrollUntilRef.current = Date.now() + SMOOTH_SCROLL_MS;
    el.scrollTo({ left: index * el.clientWidth, behavior: "smooth" });
  }, []);

  // A mouse cannot swipe: `overflow-x-auto` gives it only the scrollbar, which
  // this carousel hides. So a mouse drag moves the scroller by hand.
  const dragRef = useRef<Drag | null>(null);
  // True between the end of a drag and the click it would otherwise fire.
  const wasDragRef = useRef(false);
  const [isDragging, setIsDragging] = useState(false);

  const handlePointerDown = (event: React.PointerEvent<HTMLDivElement>) => {
    pausedUntilRef.current = Date.now() + RESUME_AFTER_MS;
    if (event.pointerType === "touch") return; // a finger scrolls natively
    const el = scrollerRef.current;
    if (el === null) return;
    wasDragRef.current = false;
    dragRef.current = {
      startX: event.clientX,
      startScroll: el.scrollLeft,
      hasMoved: false,
    };
  };

  const handlePointerMove = (event: React.PointerEvent<HTMLDivElement>) => {
    const drag = dragRef.current;
    const el = scrollerRef.current;
    if (drag === null || el === null) return;
    if (!drag.hasMoved) {
      if (Math.abs(event.clientX - drag.startX) <= DRAG_SLOP_PX) return;
      drag.hasMoved = true;
      setIsDragging(true);
      // Captured only now that it is a drag: capturing on pointerdown would
      // retarget the click and cost the video slide its tap-to-control. From
      // here it keeps the moves coming when the pointer leaves the carousel,
      // and over the video, where the iframe would otherwise take them.
      el.setPointerCapture(event.pointerId);
    }
    el.scrollLeft = drag.startScroll - (event.clientX - drag.startX);
  };

  const endDrag = () => {
    const drag = dragRef.current;
    const el = scrollerRef.current;
    dragRef.current = null;
    if (drag === null || el === null || !drag.hasMoved) return;
    wasDragRef.current = true;
    setIsDragging(false);
    // Snap scrolling does not settle a scrollLeft we wrote ourselves.
    goTo(Math.round(el.scrollLeft / el.clientWidth));
  };

  const isOnVideo = slides[current]?.type === "video";
  const shouldAutoAdvance = !isReducedMotion && slides.length > 1 && !isOnVideo;

  useEffect(() => {
    if (!shouldAutoAdvance) return;
    const timer = window.setInterval(() => {
      const el = scrollerRef.current;
      // Skip this tick while the customer is interacting or the tab is hidden.
      if (el === null || document.hidden || Date.now() < pausedUntilRef.current)
        return;
      const index = Math.round(el.scrollLeft / el.clientWidth);
      // Wrap to the first image rather than the video, which would replay it.
      const firstImage = slides.findIndex((item) => item.type === "image");
      const next =
        index + 1 >= slides.length ? Math.max(firstImage, 0) : index + 1;
      selfScrollUntilRef.current = Date.now() + SMOOTH_SCROLL_MS;
      el.scrollTo({ left: next * el.clientWidth, behavior: "smooth" });
    }, AUTO_ADVANCE_MS);
    return () => window.clearInterval(timer);
  }, [shouldAutoAdvance, slides]);

  if (slides.length === 0) {
    return <SingleImage src={fallbackImageUrl} alt={movieTitle} />;
  }

  return (
    <>
      <div
        ref={scrollerRef}
        onScroll={handleScroll}
        onPointerDown={handlePointerDown}
        onPointerMove={handlePointerMove}
        onPointerUp={endDrag}
        onPointerCancel={endDrag}
        className={cn(
          "flex h-full [scrollbar-width:none] overflow-x-auto overflow-y-hidden overscroll-x-contain [&::-webkit-scrollbar]:hidden",
          // Snapping fights a scrollLeft written by hand, so it is off mid-drag.
          isDragging ? "cursor-grabbing" : "cursor-grab snap-x snap-mandatory",
        )}
      >
        {slides.map((item, index) => (
          <div
            key={item.url}
            className="relative h-full w-full shrink-0 snap-center"
          >
            {item.type === "video" && item.youtube_key !== null ? (
              <VideoSlide
                youtubeKey={item.youtube_key}
                title={movieTitle}
                isAutoplay={!isReducedMotion}
                isHandedOver={isPlayerHandedOver}
                onHandOver={() => {
                  // A drag that happened to end on the video must not be read
                  // as a tap on it.
                  if (!wasDragRef.current) setIsPlayerHandedOver(true);
                }}
              />
            ) : (
              <Image
                src={item.url}
                alt={`${movieTitle} — ${slideLabel(slides, index).toLowerCase()}`}
                fill
                sizes="(max-width: 448px) 100vw, 448px"
                priority={index === 0}
                // Without this a mouse drag starts the browser's own image
                // drag instead of scrolling the carousel.
                draggable={false}
                className="object-cover"
              />
            )}
          </div>
        ))}
      </div>

      {slides.length > 1 && (
        <div className="absolute right-0 bottom-5 left-0 z-10 flex justify-center">
          {slides.map((item, index) => (
            <button
              key={item.url}
              type="button"
              onClick={() => goTo(index)}
              aria-label={slideLabel(slides, index)}
              aria-current={index === current}
              // The dot is 6 px; the button around it is what a thumb can hit.
              className="flex h-8 items-center px-1.5 focus-visible:outline-none"
            >
              <span
                className={cn(
                  "h-1.5 rounded-full transition-all",
                  index === current
                    ? "w-4 bg-primary"
                    : "w-1.5 bg-foreground/40",
                )}
              />
            </button>
          ))}
        </div>
      )}
    </>
  );
}

type VideoSlideProps = {
  youtubeKey: string;
  title: string;
  isAutoplay: boolean;
  /** True once the customer tapped the slide to control the player. */
  isHandedOver: boolean;
  onHandOver: () => void;
};

/**
 * YouTube is the only source TMDB gives for trailers, so the slide is an
 * embed. `youtube-nocookie.com` keeps Google from setting its usual cookies,
 * and autoplay has to be muted or browsers block it outright.
 */
function VideoSlide({
  youtubeKey,
  title,
  isAutoplay,
  isHandedOver,
  onHandOver,
}: VideoSlideProps) {
  const params = new URLSearchParams({
    mute: "1",
    playsinline: "1",
    rel: "0",
    modestbranding: "1",
    iv_load_policy: "3", // no annotation overlays
    // `rel=0` no longer removes YouTube's suggestions; since 2018 it only keeps
    // them to the same channel. Looping is what stops the end screen: the
    // trailer restarts instead of ending. `loop` needs `playlist` to name what
    // to repeat.
    loop: "1",
    playlist: youtubeKey,
    autoplay: isAutoplay ? "1" : "0",
  });
  return (
    <>
      <iframe
        src={`https://www.youtube-nocookie.com/embed/${encodeURIComponent(youtubeKey)}?${params.toString()}`}
        title={`Trailer ${title}`}
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
        allowFullScreen
        className="h-full w-full border-0"
      />
      {!isHandedOver && (
        // A cross-origin iframe keeps every touch that lands on it, so a swipe
        // over the player would never reach the carousel. This layer takes the
        // gesture instead, and a tap on it gives the player back.
        <button
          type="button"
          onClick={onHandOver}
          aria-label="Chạm để điều khiển trailer"
          className="absolute inset-0 focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-inset"
        />
      )}
    </>
  );
}

/** The poster, or a placeholder, while the detail request is still running. */
function SingleImage({ src, alt }: { src: string | null; alt: string }) {
  if (src === null) {
    return (
      <div className="flex h-full w-full items-center justify-center">
        <Film
          aria-hidden
          className="size-10 text-muted-foreground"
          strokeWidth={1.5}
        />
      </div>
    );
  }
  return (
    <Image
      src={src}
      alt={alt}
      fill
      sizes="(max-width: 448px) 100vw, 448px"
      priority
      className="object-cover"
    />
  );
}
