import { cn } from "@/lib/utils";

type GenreChipProps = {
  label: string;
  variant?: "gold" | "neutral";
  className?: string;
};

/** Genre pill (DESIGN-SYSTEM §5): first genre `gold`, second `neutral`. */
export function GenreChip({
  label,
  variant = "neutral",
  className,
}: GenreChipProps) {
  return (
    <span
      className={cn(
        "inline-flex max-w-full items-center truncate rounded-full border px-2 py-0.5 text-[11px] leading-4",
        variant === "gold"
          ? "border-gold-line bg-gold-soft text-primary"
          : "border-border bg-muted text-muted-foreground",
        className,
      )}
    >
      {label}
    </span>
  );
}
