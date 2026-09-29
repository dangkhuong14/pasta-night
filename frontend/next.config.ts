import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  poweredByHeader: false,
  images: {
    // TMDB CDN only: posters, backdrops, provider logos (PROJECT-RULES §4 Security).
    // Images are proxied by the Next.js server, so browsers never call image.tmdb.org directly.
    remotePatterns: [new URL("https://image.tmdb.org/t/p/**")],
  },
};

export default nextConfig;
