"use client";

import type { MouseEvent } from "react";

import { ABOUT_PARAM } from "./constants";

/**
 * Opens the about sheet by adding `?about=1` with `history.pushState`, the
 * same trick `MovieCard` uses: no server round trip, and back closes it.
 * Falling back to a plain link keeps it working before hydration.
 */
export function AboutLink() {
  function handleClick(event: MouseEvent<HTMLAnchorElement>) {
    if (
      event.button !== 0 ||
      event.metaKey ||
      event.ctrlKey ||
      event.shiftKey ||
      event.altKey
    )
      return;
    event.preventDefault();
    // Replaces the query rather than adding to it, so only one sheet is ever
    // open: the two Drawers are siblings in the layout and cannot nest. From
    // inside the movie sheet this swaps to the about sheet, and back returns
    // to the movie because pushState left a history entry.
    window.history.pushState(null, "", `?${ABOUT_PARAM}=1`);
  }

  return (
    <a
      href={`?${ABOUT_PARAM}=1`}
      onClick={handleClick}
      className="rounded-full text-[13px] text-primary transition hover:text-gold-bright focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
    >
      Về Pásta Night
    </a>
  );
}
