import Link from "next/link";

import { ChevronLeft } from "lucide-react";

import { BrandWordmark } from "./BrandWordmark";

type PageHeaderProps = {
  backHref: string;
};

/** Sticky header: gold back chevron + centered wordmark (DESIGN-SYSTEM §5). */
export function PageHeader({ backHref }: PageHeaderProps) {
  return (
    <header className="sticky top-0 z-30 grid h-14 grid-cols-[44px_1fr_44px] items-center border-b border-border bg-background/80 px-2 backdrop-blur">
      <Link
        href={backHref}
        aria-label="Quay lại"
        className="flex size-11 items-center justify-center rounded-full text-primary transition hover:text-gold-bright focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
      >
        <ChevronLeft aria-hidden className="size-6" strokeWidth={1.5} />
      </Link>
      <div className="text-center">
        <BrandWordmark className="text-lg" />
      </div>
    </header>
  );
}
