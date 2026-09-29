import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

import { SectionLabel } from "./SectionLabel";

type InfoCardProps = {
  label: string;
  /** `gold` for the pasta pairing card (DESIGN-SYSTEM §5 PastaPairingCard). */
  variant?: "default" | "gold";
  children: ReactNode;
  className?: string;
};

/** Bordered card with a section label on top (DESIGN-SYSTEM §5 InfoCard). */
export function InfoCard({
  label,
  variant = "default",
  children,
  className,
}: InfoCardProps) {
  return (
    <section
      className={cn(
        "rounded-xl border bg-card p-4",
        variant === "gold" ? "border-gold-line" : "border-border",
        className,
      )}
    >
      <SectionLabel as="h3" variant={variant === "gold" ? "gold" : "muted"}>
        {label}
      </SectionLabel>
      <div className="mt-3">{children}</div>
    </section>
  );
}
