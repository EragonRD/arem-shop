"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { useI18n } from "@/lib/i18n/useI18n";

const navItems = [
  { href: "/dashboard", key: "navDashboard" as const },
  { href: "/products", key: "navProducts" as const },
  { href: "/products/new", key: "navProductsNew" as const },
  { href: "/transactions/new", key: "navTransaction" as const },
];

export function Sidebar() {
  const pathname = usePathname();
  const { t } = useI18n();

  return (
    <aside className="hidden min-h-screen w-64 flex-col border-r border-brand-border/60 bg-white/75 p-6 backdrop-blur lg:flex">
      <div className="mb-8">
        <p className="font-display text-2xl font-bold text-brand-fg">Arem</p>
        <p className="text-sm text-brand-fg/65">{t("subtitle")}</p>
      </div>

      <nav className="space-y-2">
        {navItems.map((item) => {
          const active = pathname === item.href || pathname.startsWith(`${item.href}/`);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`block rounded-xl px-3 py-2 text-sm font-medium transition ${
                active
                  ? "bg-gradient-to-r from-brand-primary to-brand-accent text-white shadow-soft"
                  : "text-brand-fg/75 hover:bg-brand-secondary/40"
              }`}
            >
              {t(item.key)}
            </Link>
          );
        })}
      </nav>

      <div className="mt-auto rounded-xl border border-brand-border/70 bg-brand-secondary/45 p-3 text-xs text-brand-fg/70">
        {t("storefrontHint")}
      </div>
    </aside>
  );
}
