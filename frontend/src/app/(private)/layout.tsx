"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { AppShell } from "@/components/layout/AppShell";
import { shouldRedirectToLogin } from "@/lib/auth/guards";
import { useI18n } from "@/lib/i18n/useI18n";

export default function PrivateLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const { t } = useI18n();
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const onUnauthorized = () => {
      router.replace("/login");
    };

    if (typeof window !== "undefined") {
      window.addEventListener("arem:unauthorized", onUnauthorized);
    }

    if (shouldRedirectToLogin(pathname)) {
      router.replace("/login");
      return () => {
        if (typeof window !== "undefined") {
          window.removeEventListener("arem:unauthorized", onUnauthorized);
        }
      };
    }

    setReady(true);

    return () => {
      if (typeof window !== "undefined") {
        window.removeEventListener("arem:unauthorized", onUnauthorized);
      }
    };
  }, [pathname, router]);

  if (!ready) {
    return (
      <main className="flex min-h-screen items-center justify-center text-sm text-brand-fg/70">{t("loading")}</main>
    );
  }

  return <AppShell>{children}</AppShell>;
}
