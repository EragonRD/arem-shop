"use client";

import { useRouter } from "next/navigation";

import { LanguageToggle } from "@/components/common/LanguageToggle";
import { clearSession, getUser } from "@/lib/auth/session";
import { useI18n } from "@/lib/i18n/useI18n";

export function TopNav() {
  const router = useRouter();
  const { t } = useI18n();
  const user = getUser();

  const handleLogout = () => {
    clearSession();
    router.replace("/login");
  };

  return (
    <header className="sticky top-0 z-20 flex items-center justify-between border-b border-brand-border/60 bg-white/75 px-4 py-3 backdrop-blur lg:px-8">
      <div>
        <p className="font-display text-lg font-semibold text-brand-fg truncate max-w-[150px]">
          {user?.shopName || t("brand")}
        </p>
        <p className="text-xs text-brand-fg/70">
          {t("role")}: {user?.role ?? "-"}
        </p>
      </div>

      <div className="flex items-center gap-3">
        <LanguageToggle />
        <button
          type="button"
          onClick={handleLogout}
          className="rounded-full border border-brand-border bg-white px-4 py-2 text-sm font-medium text-brand-fg transition hover:border-brand-primary hover:text-brand-primary"
        >
          {t("logout")}
        </button>
      </div>
    </header>
  );
}
