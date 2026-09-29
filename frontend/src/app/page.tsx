import { connection } from "next/server";

import { BrandWordmark } from "@/components/common/BrandWordmark";
import { SectionLabel } from "@/components/common/SectionLabel";
import { OptionList } from "@/features/options/OptionList";
import { getOptions } from "@/lib/api";
import { publicAsset } from "@/lib/assets";

/** Landing (SCREENS §1): "who are you watching with tonight?" */
export default async function HomePage() {
  // Render per request so `next build` never needs the API; the fetch itself is cached for 5 minutes.
  await connection();
  const { data: options } = await getOptions();
  if (options.length === 0) {
    // SCREENS §1: an empty option list is a misconfiguration → error.tsx.
    throw new Error("GET /options returned no options");
  }
  const imageSrcById = Object.fromEntries(
    options.map(
      (o) => [o.id, publicAsset(`images/options/${o.id}.webp`)] as const,
    ),
  );

  return (
    <div className="flex flex-1 flex-col bg-gold-glow px-5 pt-12">
      <header className="text-center">
        <BrandWordmark />
        <SectionLabel className="mt-5">Gợi ý phim tối nay</SectionLabel>
        <h1 className="mt-3 font-serif text-[28px] leading-tight text-foreground">
          Tối nay bạn xem phim cùng ai?
        </h1>
      </header>
      <div className="mt-8">
        <OptionList options={options} imageSrcById={imageSrcById} />
      </div>
    </div>
  );
}
