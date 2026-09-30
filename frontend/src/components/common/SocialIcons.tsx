import type { SVGProps } from "react";

// lucide-react v1 dropped its brand icons, so the three social glyphs are
// inline. They are drawn on the same 24px grid as lucide so sizes match.

type IconProps = SVGProps<SVGSVGElement>;

const base = {
  viewBox: "0 0 24 24",
  fill: "currentColor",
  xmlns: "http://www.w3.org/2000/svg",
  "aria-hidden": true,
} as const;

export function FacebookIcon(props: IconProps) {
  return (
    <svg {...base} {...props}>
      <path d="M22 12.06C22 6.5 17.52 2 12 2S2 6.5 2 12.06c0 5.02 3.66 9.18 8.44 9.94v-7.03H7.9v-2.91h2.54V9.85c0-2.52 1.49-3.91 3.77-3.91 1.09 0 2.24.2 2.24.2v2.46h-1.26c-1.24 0-1.63.78-1.63 1.57v1.89h2.78l-.45 2.91h-2.33V22c4.78-.76 8.44-4.92 8.44-9.94Z" />
    </svg>
  );
}

export function InstagramIcon(props: IconProps) {
  return (
    <svg
      {...base}
      fill="none"
      stroke="currentColor"
      strokeWidth={1.6}
      {...props}
    >
      <rect x="3" y="3" width="18" height="18" rx="5" />
      <circle cx="12" cy="12" r="4" />
      <circle cx="17.3" cy="6.7" r="1.1" fill="currentColor" stroke="none" />
    </svg>
  );
}

export function TiktokIcon(props: IconProps) {
  return (
    <svg {...base} {...props}>
      <path d="M16.5 2h-3v13.2a2.7 2.7 0 1 1-2.2-2.65V9.5A5.75 5.75 0 1 0 16.5 15.2V8.9a6.6 6.6 0 0 0 4 1.35V7.2a3.75 3.75 0 0 1-4-3.75V2Z" />
    </svg>
  );
}
