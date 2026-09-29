import type { Metadata, Viewport } from "next";
import { Inter, Playfair_Display } from "next/font/google";

import { Footer } from "@/components/common/Footer";
import { BRAND_NAME } from "@/config/brand";

import "./globals.css";

// `vietnamese` subset: without it diacritics fall back to a system font (DESIGN-SYSTEM §1).
const playfair = Playfair_Display({
  subsets: ["latin", "vietnamese"],
  variable: "--font-playfair",
  display: "swap",
});

const inter = Inter({
  subsets: ["latin", "vietnamese"],
  variable: "--font-inter",
  display: "swap",
});

export const metadata: Metadata = {
  title: `${BRAND_NAME} — Gợi ý phim tối nay`,
  description:
    "Chọn người cùng xem, nhận gợi ý phim hợp tâm trạng cho bữa mỳ tối nay.",
};

export const viewport: Viewport = {
  themeColor: "#0a0a0a",
  colorScheme: "dark",
  viewportFit: "cover",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="vi" className={`${playfair.variable} ${inter.variable}`}>
      <body>
        {/* Pages with a StickyActionBar get room below the footer so the bar never hides it. */}
        <div className="mx-auto flex min-h-dvh w-full max-w-md flex-col has-data-[slot=sticky-action-bar]:pb-24">
          <main className="flex flex-1 flex-col">{children}</main>
          <Footer />
        </div>
      </body>
    </html>
  );
}
