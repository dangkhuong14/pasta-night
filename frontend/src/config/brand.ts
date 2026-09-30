// Brand strings and static content the API does not provide (SCREENS.md §0).

export const BRAND_NAME = "Pásta Night";

/** Shown under the wordmark in the about sheet (SCREENS.md §6). */
export const BRAND_TAGLINE = "Mỗi phần mỳ đi kèm một buổi tối trọn vẹn.";

/** "CÂU CHUYỆN": one paragraph per entry. */
export const BRAND_STORY: string[] = [
  "Khởi đầu từ một căn bếp nhỏ đam mê ẩm thực Ý, Pásta Night không chỉ muốn mang đến những sợi mỳ tươi ngon nhất mà còn muốn tạo nên một trải nghiệm hoàn chỉnh cho tâm hồn.",
  "Chúng tôi tin rằng mọi bữa ăn ngon đều xứng đáng đi kèm với một bộ phim hay. Đó là lý do mỗi phần mỳ gửi đi đều mang theo một gợi ý điện ảnh riêng biệt, được chọn lọc kỹ lưỡng để phù hợp với hương vị bạn thưởng thức.",
];

/** One opening-hours row: the days on the left, the times on the right. */
export type OpeningHours = {
  days: string;
  hours: string;
};

/**
 * Shop details for the about sheet and the footer. Every field is optional in
 * practice: an empty string hides its row or icon, so the sheet never shows a
 * blank label (SCREENS.md §6 hide rules).
 *
 * TODO(business): replace these placeholders with the real shop details before
 * launch. They come from the UX Pilot mockup, not from the business.
 */
export const shopInfo = {
  address: "123 Đường Pasta, Quận 1, TP. Hồ Chí Minh",
  /** Opens in Google Maps; empty hides the "CHỈ ĐƯỜNG" button. */
  mapUrl: "",
  openingHours: [
    { days: "Thứ 2 – Thứ 6", hours: "10:00 – 22:00" },
    { days: "Thứ 7 – CN", hours: "09:00 – 23:00" },
  ] satisfies OpeningHours[],
  /** Digits and spaces; the tel: link strips the spaces. */
  hotline: "0123 456 789",
};

/** Social profiles; an empty URL hides that icon. */
export const socialLinks = {
  facebook: "",
  instagram: "",
  tiktok: "",
};

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
