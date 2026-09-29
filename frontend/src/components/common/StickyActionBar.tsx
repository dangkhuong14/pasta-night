import type { ReactNode } from "react";

type StickyActionBarProps = {
  children: ReactNode;
};

/**
 * Fixed bottom bar inside the page column, with a fade above it and room for
 * the iOS home indicator (DESIGN-SYSTEM §5). `data-slot` lets the layout add
 * bottom padding so the footer is never hidden behind the bar.
 */
export function StickyActionBar({ children }: StickyActionBarProps) {
  return (
    <div
      data-slot="sticky-action-bar"
      className="fixed inset-x-0 bottom-0 z-20 mx-auto w-full max-w-md bg-linear-to-t from-background via-background/90 to-transparent px-4 pt-8 pb-[calc(1rem+env(safe-area-inset-bottom))]"
    >
      {children}
    </div>
  );
}
