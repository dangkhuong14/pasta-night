import { useState } from "react";

import { shuffle } from "@/lib/shuffle";

type ShuffleState<T> = {
  order: T[];
  start: number;
  /** Increments on every page so cards remount and replay their enter animation. */
  round: number;
};

/**
 * Pages through `items` (already shuffled by the server) `pageSize` at a time
 * without repeats. When fewer than `pageSize` items remain, the whole list is
 * reshuffled, with the movies on screen moved to the end so the next page
 * never repeats them (SCREENS §2).
 */
export function useShuffle<T>(items: readonly T[], pageSize: number) {
  // Initialized once: later server renders (router.refresh) must not reorder the grid.
  const [state, setState] = useState<ShuffleState<T>>(() => ({
    order: [...items],
    start: 0,
    round: 0,
  }));

  const visible = state.order.slice(state.start, state.start + pageSize);
  // Fewer than one full page in total → hide "Gợi ý khác" (SCREENS §2).
  const shouldShowNextButton = items.length >= pageSize;

  function showNext() {
    setState((prev) => {
      const nextStart = prev.start + pageSize;
      if (prev.order.length - nextStart >= pageSize) {
        return { ...prev, start: nextStart, round: prev.round + 1 };
      }
      const onScreen = prev.order.slice(prev.start, prev.start + pageSize);
      const rest = prev.order.filter((item) => !onScreen.includes(item));
      return {
        order: [...shuffle(rest), ...shuffle(onScreen)],
        start: 0,
        round: prev.round + 1,
      };
    });
  }

  return { visible, round: state.round, shouldShowNextButton, showNext };
}
