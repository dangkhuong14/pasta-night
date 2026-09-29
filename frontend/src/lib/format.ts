/**
 * Formats a runtime in minutes the way the design shows it (SCREENS.md §3):
 * 105 → "1h 45p", 60 → "1h", 45 → "45p".
 */
export function formatRuntime(minutes: number): string {
  const hours = Math.floor(minutes / 60);
  const rest = minutes % 60;
  if (hours === 0) return `${rest}p`;
  if (rest === 0) return `${hours}h`;
  return `${hours}h ${rest}p`;
}

/** Always one decimal, so badges keep the same width: 8 → "8.0". */
export function formatRating(rating: number): string {
  return rating.toFixed(1);
}

/**
 * Up to two initials for a cast avatar (no photos in API v1):
 * "Keanu Reeves" → "KR", "Zendaya" → "Z".
 */
export function initials(name: string): string {
  const words = name.trim().split(/\s+/).filter(Boolean);
  if (words.length === 0) return "";
  const first = Array.from(words[0])[0] ?? "";
  const last =
    words.length > 1 ? (Array.from(words[words.length - 1])[0] ?? "") : "";
  return (first + last).toLocaleUpperCase("vi");
}
