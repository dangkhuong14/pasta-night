import { OutlineButton } from "@/components/common/OutlineButton";

type PromoCardProps = {
  shopUrl: string;
};

/** "Cần thêm gia vị?" shop promo (DESIGN-SYSTEM §5). Hidden when NEXT_PUBLIC_SHOP_URL is unset. */
export function PromoCard({ shopUrl }: PromoCardProps) {
  return (
    <section className="flex items-center gap-4 rounded-xl border border-gold-line bg-linear-to-br from-gold-soft to-transparent p-5">
      <div className="min-w-0 flex-1">
        <h3 className="font-serif text-lg text-primary">Cần thêm gia vị?</h3>
        <p className="mt-1 text-xs leading-relaxed text-muted-foreground">
          Mua ngay Pásta Night để trải nghiệm phim thêm trọn vẹn.
        </p>
      </div>
      <OutlineButton href={shopUrl}>Mua ngay</OutlineButton>
    </section>
  );
}
