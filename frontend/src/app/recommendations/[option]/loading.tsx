import { PageHeader } from "@/components/common/PageHeader";
import { Skeleton } from "@/components/ui/skeleton";

/** Header, title bar, 2×2 card skeletons (SCREENS §2; design/04-loading-state.png). */
export default function Loading() {
  return (
    <>
      <PageHeader backHref="/" />
      <div aria-busy="true" aria-label="Đang tải" className="px-5 pt-6">
        <Skeleton className="h-7 w-52" />
        <Skeleton className="mt-2 h-3 w-36 rounded-full" />
        <div className="mt-6 grid grid-cols-2 gap-4">
          {[0, 1, 2, 3].map((i) => (
            <div key={i}>
              <Skeleton className="aspect-2/3 rounded-lg" />
              <Skeleton className="mt-2 h-3.5 w-3/4" />
              <div className="mt-2 flex gap-1.5">
                <Skeleton className="h-4 w-12 rounded-full" />
                <Skeleton className="h-4 w-12 rounded-full" />
              </div>
            </div>
          ))}
        </div>
      </div>
    </>
  );
}
