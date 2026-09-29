import Image from "next/image";

/**
 * Brand line + TMDB attribution on every page (SCREENS §5). The notice is the
 * wording required by TMDB's API Terms of Use: keep it in English, and keep
 * the TMDB logo less prominent than the Pásta Night wordmark.
 */
export function Footer() {
  return (
    <footer className="flex flex-col items-center gap-3 px-5 pt-10 pb-8 text-[10px] text-muted-foreground">
      <p className="flex items-center gap-3 tracking-[0.3em] uppercase">
        <span aria-hidden className="h-px w-10 bg-border" />
        Pásta Night
        <span aria-hidden className="h-px w-10 bg-border" />
      </p>
      <div className="flex flex-col items-center gap-1.5 text-center">
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
