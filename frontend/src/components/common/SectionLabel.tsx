import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

type SectionLabelProps = {
  children: ReactNode;
  variant?: "muted" | "gold";
  as?: "p" | "h2" | "h3";
  className?: string;
};

/** Eyebrow / section label: 11 px uppercase, wide tracking (DESIGN-SYSTEM §3). */
export function SectionLabel({
  children,
  variant = "muted",
  as: Tag = "p",
  className,
}: SectionLabelProps) {
  return (
    <Tag
      className={cn(
        "text-[11px] font-medium tracking-[0.2em] uppercase",
        variant === "gold" ? "text-primary" : "text-muted-foreground",
        className,
      )}
    >
      {children}
    </Tag>
  );
}
