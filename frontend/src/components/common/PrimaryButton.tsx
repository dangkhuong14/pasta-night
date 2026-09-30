import type { ReactNode } from "react";
import Link from "next/link";

import { cn } from "@/lib/utils";

const PRIMARY_CLASSES =
  "inline-flex h-12 w-full items-center justify-center rounded-full bg-primary px-6 text-[13px] font-semibold tracking-[0.15em] text-primary-foreground uppercase shadow-gold transition hover:bg-gold-bright focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background focus-visible:outline-none active:scale-[0.98] motion-reduce:transition-none motion-reduce:active:scale-100";

type PrimaryButtonProps = {
  children: ReactNode;
  className?: string;
} & (
  | { href: string; isExternal?: boolean; onClick?: never }
  | { href?: never; isExternal?: never; onClick: () => void }
);

/**
 * The one gold call to action per screen (DESIGN-SYSTEM §5). Renders a Link
 * when `href` is set, or a plain anchor opening a new tab for `isExternal`.
 */
export function PrimaryButton({
  children,
  className,
  href,
  isExternal = false,
  onClick,
}: PrimaryButtonProps) {
  if (href !== undefined) {
    if (isExternal) {
      return (
        <a
          href={href}
          target="_blank"
          rel="noopener noreferrer"
          className={cn(PRIMARY_CLASSES, className)}
        >
          {children}
        </a>
      );
    }
    return (
      <Link href={href} className={cn(PRIMARY_CLASSES, className)}>
        {children}
      </Link>
    );
  }
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(PRIMARY_CLASSES, className)}
    >
      {children}
    </button>
  );
}
