import type { ViewingOption } from "@/lib/api-types";

import { OptionCard } from "./OptionCard";

type OptionListProps = {
  options: ViewingOption[];
  /** Option ID → brand photo URL, resolved on the server (null = fallback). */
  imageSrcById: Record<string, string | null>;
};

/** Option cards in API order, which is the display order (SCREENS §1). */
export function OptionList({ options, imageSrcById }: OptionListProps) {
  return (
    // One column up to lg; side by side once the page is wide enough (SCREENS §1).
    <ul className="space-y-4 lg:grid lg:grid-cols-3 lg:gap-4 lg:space-y-0">
      {options.map((option, index) => (
        <li key={option.id}>
          <OptionCard
            option={option}
            imageSrc={imageSrcById[option.id] ?? null}
            isPriority={index === 0}
          />
        </li>
      ))}
    </ul>
  );
}
