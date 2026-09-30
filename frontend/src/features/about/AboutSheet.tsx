"use client";

import { useCallback, useRef } from "react";
import { usePathname, useSearchParams } from "next/navigation";

import { Clock, MapPin, Phone, X } from "lucide-react";

import { OutlineButton } from "@/components/common/OutlineButton";
import { PrimaryButton } from "@/components/common/PrimaryButton";
import { SectionLabel } from "@/components/common/SectionLabel";
import {
  FacebookIcon,
  InstagramIcon,
  TiktokIcon,
} from "@/components/common/SocialIcons";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerTitle,
} from "@/components/ui/drawer";
import {
  BRAND_STORY,
  BRAND_TAGLINE,
  shopInfo,
  socialLinks,
} from "@/config/brand";

import { ABOUT_PARAM } from "./constants";

const SHOP_URL = process.env.NEXT_PUBLIC_SHOP_URL ?? "";

const socials = [
  {
    key: "facebook",
    label: "Facebook",
    href: socialLinks.facebook,
    Icon: FacebookIcon,
  },
  {
    key: "instagram",
    label: "Instagram",
    href: socialLinks.instagram,
    Icon: InstagramIcon,
  },
  {
    key: "tiktok",
    label: "TikTok",
    href: socialLinks.tiktok,
    Icon: TiktokIcon,
  },
].filter((s) => s.href !== "");

/**
 * "Về Pásta Night": the brand story and shop details, as a sheet over
 * whatever page the customer is on (SCREENS.md §6). Driven by `?about=1`,
 * exactly like the movie detail sheet is driven by `?movie=`.
 */
export function AboutSheet() {
  const pathname = usePathname();
  const isOpen = useSearchParams().get(ABOUT_PARAM) !== null;

  // A param present at first render came from a shared link: closing then
  // replaces the URL instead of going back, which would leave the site.
  const wasOpenOnMount = useRef(isOpen);

  const close = useCallback(() => {
    if (wasOpenOnMount.current) {
      wasOpenOnMount.current = false;
      window.history.replaceState(null, "", pathname);
    } else {
      window.history.back();
    }
  }, [pathname]);

  return (
    <Drawer
      open={isOpen}
      onOpenChange={(open) => {
        if (!open) close();
      }}
    >
      <DrawerContent className="mx-auto max-h-[85dvh] w-full max-w-md rounded-t-3xl! bg-card">
        <div className="relative overflow-y-auto overscroll-contain px-5 pt-3 pb-8">
          <div
            aria-hidden
            className="mx-auto h-1 w-10 rounded-full bg-gold-line"
          />
          <DrawerClose
            aria-label="Đóng"
            className="absolute top-2 right-3 flex size-11 items-center justify-center rounded-full text-muted-foreground transition hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
          >
            <X aria-hidden className="size-5" strokeWidth={1.5} />
          </DrawerClose>

          <header className="mt-4 text-center">
            <DrawerTitle className="font-serif text-xl tracking-[0.3em] text-primary uppercase">
              Pásta Night
            </DrawerTitle>
            <DrawerDescription className="mt-3 font-serif text-lg leading-snug text-foreground italic">
              “{BRAND_TAGLINE}”
            </DrawerDescription>
          </header>

          <section className="mt-8">
            <SectionLabel as="h3">Câu chuyện</SectionLabel>
            <div className="mt-3 space-y-4 text-sm leading-relaxed text-muted-foreground">
              {BRAND_STORY.map((paragraph) => (
                <p key={paragraph.slice(0, 32)}>{paragraph}</p>
              ))}
            </div>
          </section>

          <dl className="mt-8 divide-y divide-border border-y border-border">
            {shopInfo.address !== "" && (
              <InfoRow
                icon={
                  <MapPin aria-hidden className="size-4" strokeWidth={1.5} />
                }
                label="Địa chỉ"
              >
                <div className="flex items-start justify-between gap-3">
                  <p className="text-sm text-foreground">{shopInfo.address}</p>
                  {shopInfo.mapUrl !== "" && (
                    <OutlineButton href={shopInfo.mapUrl}>
                      Chỉ đường
                    </OutlineButton>
                  )}
                </div>
              </InfoRow>
            )}

            {shopInfo.openingHours.length > 0 && (
              <InfoRow
                icon={
                  <Clock aria-hidden className="size-4" strokeWidth={1.5} />
                }
                label="Giờ mở cửa"
              >
                <ul className="space-y-1">
                  {shopInfo.openingHours.map((row) => (
                    <li
                      key={row.days}
                      className="flex justify-between gap-4 text-sm"
                    >
                      <span className="text-muted-foreground">{row.days}</span>
                      <span className="text-foreground tabular-nums">
                        {row.hours}
                      </span>
                    </li>
                  ))}
                </ul>
              </InfoRow>
            )}

            {shopInfo.hotline !== "" && (
              <InfoRow
                icon={
                  <Phone aria-hidden className="size-4" strokeWidth={1.5} />
                }
                label="Hotline"
              >
                <a
                  href={`tel:${shopInfo.hotline.replace(/\s/g, "")}`}
                  className="text-sm text-foreground transition hover:text-primary focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
                >
                  {shopInfo.hotline}
                </a>
              </InfoRow>
            )}
          </dl>

          {socials.length > 0 && (
            <ul className="mt-8 flex justify-center gap-4">
              {socials.map(({ key, label, href, Icon }) => (
                <li key={key}>
                  <a
                    href={href}
                    target="_blank"
                    rel="noopener noreferrer"
                    aria-label={label}
                    className="flex size-11 items-center justify-center rounded-full border border-gold-line text-primary transition hover:bg-gold-soft focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
                  >
                    <Icon className="size-5" strokeWidth={1.5} />
                  </a>
                </li>
              ))}
            </ul>
          )}

          {SHOP_URL !== "" && (
            <div className="mt-8">
              <PrimaryButton href={SHOP_URL} isExternal>
                Đặt thêm một phần
              </PrimaryButton>
            </div>
          )}
        </div>
      </DrawerContent>
    </Drawer>
  );
}

type InfoRowProps = {
  icon: React.ReactNode;
  label: string;
  children: React.ReactNode;
};

function InfoRow({ icon, label, children }: InfoRowProps) {
  return (
    <div className="flex gap-3 py-4">
      <span className="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-full border border-gold-line text-primary">
        {icon}
      </span>
      <div className="min-w-0 flex-1">
        <dt>
          <SectionLabel>{label}</SectionLabel>
        </dt>
        <dd className="mt-1.5">{children}</dd>
      </div>
    </div>
  );
}
