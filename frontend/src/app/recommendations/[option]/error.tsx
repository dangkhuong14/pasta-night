"use client";

import { useEffect } from "react";

import { ErrorState } from "@/components/common/ErrorState";
import { PageHeader } from "@/components/common/PageHeader";
import { PrimaryButton } from "@/components/common/PrimaryButton";
import { errorCopy } from "@/lib/error-copy";

type ErrorPageProps = {
  error: Error & { digest?: string };
  retry: () => void;
};

/**
 * Recommendations error boundary (SCREENS §2 "other errors", §4). Server
 * errors arrive without details in production, so the copy is the default.
 */
export default function ErrorPage({ error, retry }: ErrorPageProps) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  const copy = errorCopy();
  return (
    <>
      <PageHeader backHref="/" />
      <ErrorState
        title={copy.title}
        body={copy.body}
        action={<PrimaryButton onClick={retry}>Thử lại</PrimaryButton>}
      />
    </>
  );
}
