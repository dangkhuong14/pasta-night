// Brand strings and static content the API does not provide (SCREENS.md §0).

export const BRAND_NAME = "Pásta Night";

/** A Pásta Night dish suggested for an option. */
export type PastaPairing = {
  /** Image file name in public/images/dishes/{dishId}.webp (optional asset). */
  dishId: string;
  /** Dish name, e.g. "Carbonara Cổ Điển". */
  name: string;
  /** One-line note shown under the name. */
  note: string;
};

/**
 * Suggested dish per option_id. An option without an entry hides the
 * "GỢI Ý MÓN MỲ HOÀN HẢO" card (SCREENS.md §3 hide rules).
 *
 * TODO(business): add real dishes for netflix-chill, solo, and friends.
 */
export const pastaPairings: Partial<Record<string, PastaPairing>> = {};
