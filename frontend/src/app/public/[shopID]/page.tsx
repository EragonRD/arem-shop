"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

import { useI18n } from "@/lib/i18n/useI18n";
import { dataClient } from "@/lib/services/dataClient";
import type { PublicProductViewModel } from "@/lib/types";

function formatPrice(value: number): string {
  return new Intl.NumberFormat("fr-FR", { style: "currency", currency: "EUR" }).format(value);
}

export default function PublicShopPage() {
  const { shopID } = useParams<{ shopID: string }>();
  const { t } = useI18n();
  const [products, setProducts] = useState<PublicProductViewModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        setLoading(true);
        const list = await dataClient.listPublicProducts(shopID);
        setProducts(list);
      } catch (catalogError) {
        setError(catalogError instanceof Error ? catalogError.message : "failed to load public catalog");
      } finally {
        setLoading(false);
      }
    };

    void load();
  }, [shopID]);

  return (
    <main className="mx-auto min-h-screen max-w-6xl space-y-6 px-4 py-8">
      <header className="space-y-2">
        <p className="font-display text-3xl font-semibold text-brand-fg">{t("publicTitle")}</p>
        <p className="text-sm text-brand-fg/70">Shop: {shopID}</p>
        <div className="text-sm text-brand-primary underline underline-offset-2">
          <Link href="/login">Back office</Link>
        </div>
      </header>

      {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}
      {loading ? <p className="text-sm text-brand-fg/70">{t("loading")}</p> : null}
      {!loading && products.length === 0 ? <p className="text-sm text-brand-fg/70">{t("noData")}</p> : null}

      <section className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {products.map((product) => (
          <article key={product.id} className="overflow-hidden rounded-2xl border border-brand-border/70 bg-white/90 shadow-soft">
            <div className="h-44 bg-brand-secondary/35">
              {product.imageURL ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img src={product.imageURL} alt={product.name} className="h-full w-full object-cover" />
              ) : null}
            </div>
            <div className="space-y-2 p-4">
              <h2 className="font-display text-xl font-semibold text-brand-fg">{product.name}</h2>
              <p className="text-sm text-brand-fg/75">{product.description}</p>
              <p className="text-sm text-brand-fg/80">{product.category}</p>
              <p className="text-sm text-brand-fg/80">Stock: {product.stock}</p>
              <p className="text-base font-semibold text-brand-primary">{formatPrice(product.sellingPrice)}</p>
              <a
                href={product.whatsappLink}
                target="_blank"
                rel="noreferrer"
                className="inline-block rounded-full bg-gradient-to-r from-brand-primary to-brand-accent px-4 py-2 text-sm font-medium text-white"
              >
                WhatsApp
              </a>
            </div>
          </article>
        ))}
      </section>
    </main>
  );
}
