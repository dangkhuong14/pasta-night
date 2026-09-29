import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

type OutlineButtonProps = {
  children: ReactNode;
  /** External URL; opens in a new tab. */
  href: string;
  className?: string;
};

/** Small gold-outlined pill link, e.g. "MUA NGAY" (DESIGN-SYSTEM §5). */
export function OutlineButton({
  children,
  href,
  className,
}: OutlineButtonProps) {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className={cn(
        "inline-flex h-9 min-w-11 shrink-0 items-center justify-center rounded-full border border-gold-line px-4 text-[11px] font-semibold tracking-[0.15em] text-primary uppercase transition hover:bg-gold-soft focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none active:scale-[0.98] motion-reduce:transition-none",
        className,
      )}
    >
      {children}
    </a>
  );
}
