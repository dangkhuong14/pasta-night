import { Star } from "lucide-react";

import { formatRating } from "@/lib/format";
import { cn } from "@/lib/utils";

type RatingBadgeProps = {
  rating: number;
  /** `overlay` sits on a poster; `outline` sits in the detail title row. */
  variant?: "overlay" | "outline";
  className?: string;
};

/** Pill with a gold star and the rating, always 1 decimal (DESIGN-SYSTEM §5). */
export function RatingBadge({
  rating,
  variant = "overlay",
  className,
}: RatingBadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex shrink-0 items-center gap-1 rounded-full font-medium tabular-nums",
        variant === "overlay"
          ? "bg-black/60 px-2 py-0.5 text-[11px] text-foreground backdrop-blur"
          : "border border-gold-line bg-gold-soft px-2.5 py-1 text-[13px] text-primary",
        className,
      )}
    >
      <Star
        aria-hidden
        className="size-3 fill-gold-bright text-gold-bright"
        strokeWidth={1.5}
      />
      <span>
        <span className="sr-only">Điểm </span>
        {formatRating(rating)}
      </span>
    </span>
  );
}
