"use client";

import { useI18n } from "@/lib/i18n/useI18n";

export function LanguageToggle() {
  const { locale, setLocale, t } = useI18n();

  return (
    <div className="inline-flex items-center gap-1 rounded-full border border-brand-border/70 bg-white/80 p-1 text-xs shadow-soft">
      <span className="px-2 text-brand-fg/70">{t("language")}</span>
      <button
        type="button"
        onClick={() => setLocale("fr")}
        className={`rounded-full px-2 py-1 font-medium transition ${
          locale === "fr" ? "bg-brand-primary text-white" : "text-brand-fg/70"
        }`}
      >
        FR
      </button>
      <button
        type="button"
        onClick={() => setLocale("en")}
        className={`rounded-full px-2 py-1 font-medium transition ${
          locale === "en" ? "bg-brand-primary text-white" : "text-brand-fg/70"
        }`}
      >
        EN
      </button>
    </div>
  );
}
