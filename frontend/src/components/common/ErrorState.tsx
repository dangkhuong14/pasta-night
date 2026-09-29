import type { ReactNode } from "react";

import { Film } from "lucide-react";

type ErrorStateProps = {
  title: string;
  body: string;
  /** Usually a PrimaryButton; omitted while auto-retrying. */
  action?: ReactNode;
};

/** Shared error state, centered ~40% from the top (SCREENS §4). */
export function ErrorState({ title, body, action }: ErrorStateProps) {
  return (
    <div
      role="alert"
      className="flex flex-1 flex-col items-center justify-center px-8 pt-10 pb-[20dvh] text-center"
    >
      <Film aria-hidden className="size-12 text-primary" strokeWidth={1.5} />
      <h2 className="mt-5 font-serif text-xl text-foreground">{title}</h2>
      <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
        {body}
      </p>
      {action !== undefined && (
        <div className="mt-8 w-full max-w-60">{action}</div>
      )}
    </div>
  );
}
