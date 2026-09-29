import { cn } from "@/lib/utils";

type BrandWordmarkProps = {
  className?: string;
};

/** "PÁSTA NIGHT" wordmark (DESIGN-SYSTEM §3 Wordmark). */
export function BrandWordmark({ className }: BrandWordmarkProps) {
  return (
    <span
      className={cn(
        "font-serif text-xl tracking-[0.3em] text-primary uppercase",
        className,
      )}
    >
      Pásta Night
    </span>
  );
}
