import { BrandWordmark } from "@/components/common/BrandWordmark";
import { Skeleton } from "@/components/ui/skeleton";

/** Landing skeleton: wordmark + H1 + 3 option cards of 176 px (SCREENS §1 States). */
export default function Loading() {
  return (
    <div
      aria-busy="true"
      aria-label="Đang tải"
      className="flex flex-1 flex-col bg-gold-glow px-5 pt-12"
    >
      <div className="flex flex-col items-center">
        <BrandWordmark />
        <Skeleton className="mt-5 h-3 w-36 rounded-full" />
        <Skeleton className="mt-4 h-7 w-64" />
        <Skeleton className="mt-2 h-7 w-40" />
      </div>
      <div className="mt-8 space-y-4">
        {[0, 1, 2].map((i) => (
          <Skeleton key={i} className="h-44 rounded-xl" />
        ))}
      </div>
    </div>
  );
}
