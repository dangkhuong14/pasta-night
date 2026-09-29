"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { ErrorState } from "@/components/common/ErrorState";
import { errorCopy } from "@/lib/error-copy";

type CacheNotReadyProps = {
  /** From the API's Retry-After header (default 30 s). */
  retryAfterSeconds: number;
};

/**
 * Shown while the backend's first refresh has not finished (SCREENS §2,
 * CACHE_NOT_READY). Re-renders the page every `retryAfterSeconds` until the
 * list is ready; the page then renders the grid and this unmounts.
 */
export function CacheNotReady({ retryAfterSeconds }: CacheNotReadyProps) {
  const router = useRouter();

  useEffect(() => {
    const timer = window.setInterval(
      () => router.refresh(),
      Math.max(retryAfterSeconds, 1) * 1000,
    );
    return () => window.clearInterval(timer);
  }, [router, retryAfterSeconds]);

  const copy = errorCopy("CACHE_NOT_READY");
  return <ErrorState title={copy.title} body={copy.body} />;
}
