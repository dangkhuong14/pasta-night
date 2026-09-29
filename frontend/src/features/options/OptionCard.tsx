import Image from "next/image";
import Link from "next/link";

import type { ViewingOption } from "@/lib/api-types";

import { OptionIcon } from "./OptionIcon";

export type OptionCardProps = {
  option: ViewingOption;
  /** `/images/options/{id}.webp` when the brand photo exists, else null (fallback). */
  imageSrc: string | null;
  /** Load eagerly: only the first card is above the fold. */
  isPriority?: boolean;
};

/** Full-width card linking to the option's recommendations (DESIGN-SYSTEM §5 OptionCard). */
export function OptionCard({
  option,
  imageSrc,
  isPriority = false,
}: OptionCardProps) {
  return (
    <Link
      href={`/recommendations/${encodeURIComponent(option.id)}`}
      className="group relative block h-44 overflow-hidden rounded-xl border border-border bg-card transition focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none active:scale-[0.98] motion-reduce:transition-none motion-reduce:active:scale-100"
    >
      {imageSrc !== null ? (
        <Image
          src={imageSrc}
          alt=""
          fill
          sizes="(max-width: 448px) 100vw, 448px"
          priority={isPriority}
          className="object-cover transition duration-500 group-hover:scale-105 motion-reduce:transition-none"
        />
      ) : (
        <div
          aria-hidden
          className="absolute inset-0 flex items-start justify-end bg-gold-glow p-5"
        >
          <OptionIcon
            name={option.icon}
            className="size-16 text-primary/25"
            strokeWidth={1}
          />
        </div>
      )}
      <div
        aria-hidden
        className="absolute inset-0 bg-linear-to-t from-black/90 via-black/40 to-transparent"
      />
      <div className="absolute inset-x-0 bottom-0 p-4">
        <h2 className="font-serif text-xl text-primary">{option.label}</h2>
        <p className="mt-1 text-[13px] leading-relaxed text-foreground/90">
          {option.description}
        </p>
      </div>
    </Link>
  );
}
