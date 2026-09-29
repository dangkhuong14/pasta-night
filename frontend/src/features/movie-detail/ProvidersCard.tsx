import Image from "next/image";

import { InfoCard } from "@/components/common/InfoCard";
import type { Provider } from "@/lib/api-types";

type ProvidersCardProps = {
  providers: Provider[];
};

/**
 * "CÓ MẶT TRÊN" with provider logos (DESIGN-SYSTEM §5 ProviderLogo). The
 * JustWatch caption is required by TMDB's watch-provider license.
 * The parent hides this card when `providers` is empty.
 */
export function ProvidersCard({ providers }: ProvidersCardProps) {
  return (
    <InfoCard label="Có mặt trên">
      <ul className="flex flex-wrap gap-x-5 gap-y-3">
        {providers.map((p) => (
          <li
            key={p.id}
            className="flex w-14 flex-col items-center gap-1.5 text-center"
          >
            {p.logo_url !== null ? (
              <Image
                src={p.logo_url}
                alt=""
                width={32}
                height={32}
                className="size-8 rounded-md border border-border"
              />
            ) : null}
            <span className="line-clamp-2 text-[10px] leading-tight text-muted-foreground">
              {p.name}
            </span>
          </li>
        ))}
      </ul>
      <p className="mt-3 text-[10px] text-muted-foreground">Nguồn: JustWatch</p>
    </InfoCard>
  );
}
