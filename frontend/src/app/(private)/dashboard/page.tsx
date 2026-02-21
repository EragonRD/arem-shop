"use client";

import { useEffect, useState } from "react";

import { StatCard } from "@/components/common/StatCard";
import { getSession } from "@/lib/auth/session";
import { formatCurrency } from "@/lib/format";
import { useI18n } from "@/lib/i18n/useI18n";
import { dataClient } from "@/lib/services/dataClient";
import type { DashboardViewModel } from "@/lib/types";

export default function DashboardPage() {
  const { t } = useI18n();
  const session = getSession();
  const [data, setData] = useState<DashboardViewModel | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchData = async () => {
      if (!session) {
        return;
      }
      if (session.user.role !== "SuperAdmin") {
        return;
      }

      try {
        const dashboard = await dataClient.getDashboard(session.token);
        setData(dashboard);
      } catch (dashboardError) {
        setError(dashboardError instanceof Error ? dashboardError.message : "failed to load dashboard");
      }
    };

    void fetchData();
  }, [session]);

  if (!session || session.user.role !== "SuperAdmin") {
    return (
      <section className="rounded-2xl border border-brand-border/70 bg-white/90 p-6 text-brand-fg shadow-soft">
        <h1 className="font-display text-2xl font-semibold">{t("dashboardTitle")}</h1>
        <p className="mt-3 text-sm text-brand-fg/75">{t("forbidden")}</p>
      </section>
    );
  }

  return (
    <section className="space-y-5">
      <header>
        <h1 className="font-display text-3xl font-semibold text-brand-fg">{t("dashboardTitle")}</h1>
        <p className="text-sm text-brand-fg/70">{t("superAdminNote")}</p>
      </header>

      {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <StatCard label="Total Sales" value={formatCurrency(data?.totalSales ?? 0)} />
        <StatCard label="Total Expenses" value={formatCurrency(data?.totalExpenses ?? 0)} />
        <StatCard label="Net Profit" value={formatCurrency(data?.netProfit ?? 0)} />
        <StatCard label="Low Stock" value={String(data?.lowStockProducts ?? 0)} />
      </div>

      <div className="rounded-2xl border border-brand-border/70 bg-white/90 p-6 shadow-soft">
        <h2 className="font-display text-xl font-semibold text-brand-fg">Paramètres de la Boutique</h2>
        <p className="mb-4 text-sm text-brand-fg/70">Modifiez le nom et le numéro WhatsApp de votre boutique.</p>
        <form
          onSubmit={async (e) => {
            e.preventDefault();
            const formData = new FormData(e.currentTarget);
            const name = formData.get("shopName") as string;
            const whatsAppNumber = formData.get("whatsAppNumber") as string;

            try {
              await dataClient.updateShop(session.token, { name, whatsAppNumber });

              // Mettre à jour la session locale pour affichage immédiat
              const newSession = {
                ...session,
                user: {
                  ...session.user,
                  shopName: name,
                  whatsAppNumber: whatsAppNumber,
                }
              };
              localStorage.setItem("arem_user", JSON.stringify(newSession));
              window.location.reload();
            } catch (err) {
              setError(err instanceof Error ? err.message : "Erreur lors de la mise à jour");
            }
          }}
          className="flex flex-col max-w-sm gap-4"
        >
          <div>
            <label className="mb-1 block text-sm font-medium text-brand-fg">Nom de la boutique</label>
            <input
              type="text"
              name="shopName"
              defaultValue={session.user.shopName || ""}
              placeholder="Ex: Ma Boutique"
              className="w-full rounded-xl border border-brand-border bg-brand-bg px-3 py-2 text-sm text-brand-fg focus:border-brand-primary focus:outline-none focus:ring-1 focus:ring-brand-primary"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-brand-fg">Numéro WhatsApp</label>
            <input
              type="text"
              name="whatsAppNumber"
              defaultValue={session.user.whatsAppNumber || ""}
              placeholder="Ex: +33612345678"
              className="w-full rounded-xl border border-brand-border bg-brand-bg px-3 py-2 text-sm text-brand-fg focus:border-brand-primary focus:outline-none focus:ring-1 focus:ring-brand-primary"
            />
          </div>
          <button
            type="submit"
            className="rounded-xl bg-brand-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-brand-primary/90 mt-2"
          >
            Enregistrer
          </button>
        </form>
      </div>
    </section>
  );
}
