"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

import { saveSession } from "@/lib/auth/session";
import { useI18n } from "@/lib/i18n/useI18n";
import { dataClient, dataMode } from "@/lib/services/dataClient";

const DEFAULT_SHOP_ID = "11111111-1111-1111-1111-111111111111";

export default function LoginPage() {
  const router = useRouter();
  const { t } = useI18n();

  const [email, setEmail] = useState("owner@shopdemo.com");
  const [password, setPassword] = useState("ChangeMe123!");
  const [shopID, setShopID] = useState(DEFAULT_SHOP_ID);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError("");

    try {
      setSubmitting(true);
      const session = await dataClient.login(email, password, shopID);
      saveSession(session);
      router.replace("/dashboard");
    } catch (loginError) {
      setError(loginError instanceof Error ? loginError.message : "login failed");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <main className="relative flex min-h-screen items-center justify-center overflow-hidden p-4">
      <div className="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(circle_at_10%_20%,rgba(97,87,234,0.2),transparent_35%),radial-gradient(circle_at_90%_10%,rgba(238,150,214,0.2),transparent_35%),radial-gradient(circle_at_40%_80%,rgba(97,87,234,0.12),transparent_40%)]" />

      <section className="animate-fade-up w-full max-w-md rounded-3xl border border-brand-border/70 bg-white/90 p-8 shadow-soft backdrop-blur">
        <div className="mb-6">
          <p className="font-display text-3xl font-bold text-brand-fg">Arem Shop</p>
          <p className="text-sm text-brand-fg/70">{t("loginTitle")}</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <label className="block space-y-1 text-sm">
            <span>{t("loginEmail")}</span>
            <input
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
              type="email"
              required
            />
          </label>

          <label className="block space-y-1 text-sm">
            <span>{t("loginPassword")}</span>
            <input
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
              type="password"
              required
            />
          </label>

          <label className="block space-y-1 text-sm">
            <span>{t("loginShopID")}</span>
            <input
              value={shopID}
              onChange={(event) => setShopID(event.target.value)}
              className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
              required
            />
          </label>

          {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}

          <button
            type="submit"
            disabled={submitting}
            className="w-full rounded-full bg-gradient-to-r from-brand-primary to-brand-accent px-5 py-2.5 text-sm font-semibold text-white transition hover:opacity-95 disabled:opacity-60"
          >
            {submitting ? t("loading") : t("loginSubmit")}
          </button>
        </form>

        <div className="mt-5 flex items-center justify-between text-xs text-brand-fg/65">
          <span>Mode: {dataMode}</span>
          <Link href={`/public/${shopID || DEFAULT_SHOP_ID}`} className="underline underline-offset-2">
            {t("navPublic")}
          </Link>
        </div>
      </section>
    </main>
  );
}
