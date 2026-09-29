import { SectionLabel } from "@/components/common/SectionLabel";
import { Skeleton } from "@/components/ui/skeleton";
import { initials } from "@/lib/format";

const MAX_CAST = 5;
const SKELETON_COUNT = 4;

type CastListProps = {
  /** null while the detail is loading. */
  cast: string[] | null;
};

/**
 * "DIỄN VIÊN CHÍNH": initials avatars, since API v1 has no cast photos
 * (SCREENS §0). The parent hides it when the cast is empty or failed to load.
 */
export function CastList({ cast }: CastListProps) {
  return (
    <section>
      <SectionLabel as="h3">Diễn viên chính</SectionLabel>
      <ul className="-mx-5 mt-3 flex gap-4 overflow-x-auto px-5 pb-1">
        {cast === null
          ? Array.from({ length: SKELETON_COUNT }, (_, i) => (
              <li
                key={i}
                className="flex w-16 shrink-0 flex-col items-center gap-2"
              >
                <Skeleton className="size-12 rounded-full" />
                <Skeleton className="h-2.5 w-12 rounded-full" />
              </li>
            ))
          : cast.slice(0, MAX_CAST).map((name) => (
              <li
                key={name}
                className="flex w-16 shrink-0 flex-col items-center gap-2 text-center"
              >
                <span
                  aria-hidden
                  className="flex size-12 items-center justify-center rounded-full border border-gold-line bg-muted font-serif text-sm text-primary"
                >
                  {initials(name)}
                </span>
                <span className="line-clamp-2 text-xs leading-tight text-muted-foreground">
                  {name}
                </span>
              </li>
            ))}
      </ul>
    </section>
  );
}
