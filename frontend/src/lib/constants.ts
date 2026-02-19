import type { Locale } from "@/lib/types";

export const STORAGE_KEYS = {
  locale: "arem_locale",
  token: "arem_token",
  user: "arem_user",
} as const;

export const DEFAULT_LOCALE: Locale = "fr";
export const SUPPORTED_LOCALES: Locale[] = ["fr", "en"];
