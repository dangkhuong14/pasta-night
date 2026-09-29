"use client";

import { useEffect } from "react";

import { ErrorState } from "@/components/common/ErrorState";
import { PrimaryButton } from "@/components/common/PrimaryButton";
import { errorCopy } from "@/lib/error-copy";

type ErrorPageProps = {
  error: Error & { digest?: string };
  retry: () => void;
};

/**
 * Landing error boundary (SCREENS §1, §4). Server errors reach the client
 * without their details in production, so the copy is always the default.
 */
export default function ErrorPage({ error, retry }: ErrorPageProps) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  const copy = errorCopy();
  return (
    <ErrorState
      title={copy.title}
      body={copy.body}
      action={<PrimaryButton onClick={retry}>Thử lại</PrimaryButton>}
    />
  );
}
