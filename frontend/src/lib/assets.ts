import "server-only";

import { existsSync } from "node:fs";
import path from "node:path";

const PUBLIC_DIR = path.join(process.cwd(), "public");

/**
 * Returns the public URL of an optional brand asset, or null when the file is
 * not in `public/`. Brand photos (option backgrounds, dishes) are added later
 * by the business; until then the UI falls back per DESIGN-SYSTEM §6.
 * Paths come partly from API IDs, so anything resolving outside `public/` is null.
 */
export function publicAsset(relativePath: string): string | null {
  const file = path.resolve(PUBLIC_DIR, relativePath);
  if (!file.startsWith(PUBLIC_DIR + path.sep)) return null;
  return existsSync(file) ? `/${relativePath}` : null;
}
