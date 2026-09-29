import Image from "next/image";

import { InfoCard } from "@/components/common/InfoCard";

/** Pairing resolved on the server from `config/brand.ts` + `public/images/dishes/`. */
export type PastaPairingView = {
  name: string;
  note: string;
  /** `/images/dishes/{dishId}.webp`, or null to hide the image slot. */
  imageSrc: string | null;
};

type PastaPairingCardProps = {
  pairing: PastaPairingView;
};

/** "GỢI Ý MÓN MỲ HOÀN HẢO" (DESIGN-SYSTEM §5 PastaPairingCard). */
export function PastaPairingCard({ pairing }: PastaPairingCardProps) {
  return (
    <InfoCard label="Gợi ý món mỳ hoàn hảo" variant="gold">
      <div className="flex items-center gap-3">
        {pairing.imageSrc !== null && (
          <Image
            src={pairing.imageSrc}
            alt={pairing.name}
            width={56}
            height={56}
            className="size-14 shrink-0 rounded-lg object-cover"
          />
        )}
        <div className="min-w-0">
          <p className="font-serif text-base text-foreground">{pairing.name}</p>
          <p className="mt-0.5 text-xs text-muted-foreground">{pairing.note}</p>
        </div>
      </div>
    </InfoCard>
  );
}
