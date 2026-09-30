"use client";

import { useEffect, useRef, useState } from "react";

import { ArrowLeft, Check, Share2 } from "lucide-react";

import { DrawerClose } from "@/components/ui/drawer";
import { cn } from "@/lib/utils";

const ACTION_CLASSES =
  "inline-flex h-10 w-full items-center justify-center gap-2 rounded-full border border-border bg-muted/40 px-4 text-xs text-muted-foreground transition hover:border-gold-line hover:text-primary focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none active:scale-[0.98] motion-reduce:transition-none motion-reduce:active:scale-100";

const COPIED_RESET_MS = 2000;

type SheetActionsProps = {
  movieTitle: string;
};

/**
 * Closing note plus the two actions under it (SCREENS §3): share this
 * recommendation, or go back to the grid. Sharing uses the Web Share sheet on
 * phones — the device this app is opened on — and falls back to copying the
 * link where the browser has no share sheet.
 */
export function SheetActions({ movieTitle }: SheetActionsProps) {
  const [isCopied, setIsCopied] = useState(false);
  const timerRef = useRef<number | undefined>(undefined);

  useEffect(() => () => window.clearTimeout(timerRef.current), []);

  async function handleShare() {
    // The URL carries ?movie={id}, so the link reopens this exact sheet.
    const url = window.location.href;
    const payload = {
      title: movieTitle,
      text: `Tối nay xem "${movieTitle}" nhé!`,
      url,
    };

    if (typeof navigator.share === "function") {
      try {
        await navigator.share(payload);
        return;
      } catch {
        return; // the customer dismissed the share sheet, or it failed; nothing to report
      }
    }
    try {
      await navigator.clipboard.writeText(url);
      setIsCopied(true);
      timerRef.current = window.setTimeout(
        () => setIsCopied(false),
        COPIED_RESET_MS,
      );
    } catch {
      // Clipboard blocked (insecure origin, denied permission): leave the label as it was.
    }
  }

  return (
    <section className="space-y-4 pt-2">
      <p className="text-center font-serif text-sm text-muted-foreground italic">
        Chúc bạn ngon miệng.
      </p>
      <div className="grid grid-cols-2 gap-3">
        <button
          type="button"
          onClick={handleShare}
          className={cn(ACTION_CLASSES)}
        >
          {isCopied ? (
            <>
              <Check aria-hidden className="size-4" strokeWidth={1.5} />
              Đã sao chép
            </>
          ) : (
            <>
              <Share2 aria-hidden className="size-4" strokeWidth={1.5} />
              Chia sẻ gợi ý
            </>
          )}
        </button>
        <DrawerClose className={cn(ACTION_CLASSES)}>
          <ArrowLeft aria-hidden className="size-4" strokeWidth={1.5} />
          Chọn phim khác
        </DrawerClose>
      </div>
    </section>
  );
}
