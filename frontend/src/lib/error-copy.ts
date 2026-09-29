/** Vietnamese copy for a user-facing error state. */
export type ErrorCopy = {
  title: string;
  body: string;
};

/**
 * Copy per API error code (API_SPEC §6), plus UI-only states. Never show the
 * API `message` to users: it is English and meant for developers.
 */
export const ERROR_COPY: Record<string, ErrorCopy> = {
  CACHE_NOT_READY: {
    title: "Đang chuẩn bị gợi ý…",
    body: "Vui lòng đợi trong giây lát.",
  },
  NO_RECOMMENDATIONS: {
    title: "Chưa có gợi ý",
    body: "Hiện chưa có phim cho lựa chọn này. Bạn thử lựa chọn khác nhé.",
  },
  DEFAULT: {
    title: "Rất tiếc, đã có lỗi xảy ra",
    body: "Vui lòng thử lại sau ít phút.",
  },
};

/** Returns the copy for an error code, falling back to the generic message. */
export function errorCopy(code?: string): ErrorCopy {
  return (
    (code !== undefined ? ERROR_COPY[code] : undefined) ?? ERROR_COPY.DEFAULT
  );
}
