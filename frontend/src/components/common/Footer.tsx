import type { ReactNode } from "react";
import Image from "next/image";

import { MapPin, Phone } from "lucide-react";

import { shopInfo, socialLinks } from "@/config/brand";

import { BrandWordmark } from "./BrandWordmark";
import { FacebookIcon, InstagramIcon, TiktokIcon } from "./SocialIcons";

type FooterProps = {
  /** The "Về Pásta Night" link. Passed in so this stays free of feature imports. */
  aboutLink?: ReactNode;
};

/** Contact and social icons; each is skipped when its destination is unset. */
const contactIcons = [
  {
    key: "phone",
    label: "Gọi hotline",
    href: telHref(shopInfo.hotline),
    Icon: Phone,
  },
  { key: "map", label: "Xem địa chỉ", href: shopInfo.mapUrl, Icon: MapPin },
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
].filter((item) => item.href !== "");

function telHref(hotline: string): string {
  const digits = hotline.replace(/\s/g, "");
  return digits === "" ? "" : `tel:${digits}`;
}

/**
 * Brand line, contact and social icons, the about link, and the TMDB
 * attribution on every page (SCREENS.md §5). The notice is the wording
 * required by TMDB's API Terms of Use: keep it in English, and keep the TMDB
 * logo less prominent than the Pásta Night wordmark.
 */
export function Footer({ aboutLink }: FooterProps) {
  return (
    <footer className="flex flex-col items-center gap-4 px-5 pt-8 pb-8 text-center">
      {/* Sets the footer off from the page, and from the sheet body when the
          movie detail sheet renders its own copy (SCREENS.md §5). */}
      <span aria-hidden className="mb-6 h-px w-full bg-gold-line/50" />
      <BrandWordmark className="text-base" />

      {contactIcons.length > 0 && (
        <ul className="flex items-center gap-2">
          {contactIcons.map(({ key, label, href, Icon }) => (
            <li key={key}>
              <a
                href={href}
                aria-label={label}
                {...(href.startsWith("http")
                  ? { target: "_blank", rel: "noopener noreferrer" }
                  : {})}
                className="flex size-11 items-center justify-center rounded-full text-primary transition hover:text-gold-bright focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
              >
                <Icon className="size-5" strokeWidth={1.5} />
              </a>
            </li>
          ))}
        </ul>
      )}

      {aboutLink}

      <div className="mt-2 flex flex-col items-center gap-1.5 text-[10px] text-muted-foreground">
        <Image
          src="/images/brand/tmdb-logo.svg"
          alt="TMDB"
          width={77}
          height={10}
          unoptimized
        />
        <p lang="en">
          This website uses TMDB and the TMDB APIs but is not endorsed,
          certified, or otherwise approved by TMDB.
        </p>
      </div>
    </footer>
  );
}
