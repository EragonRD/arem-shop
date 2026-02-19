"use client";

import { useI18nContext } from "@/lib/i18n/I18nProvider";

export function useI18n() {
  return useI18nContext();
}
