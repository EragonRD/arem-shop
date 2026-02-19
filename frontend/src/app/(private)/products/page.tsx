"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { ProductCard } from "@/components/products/ProductCard";
import { getSession } from "@/lib/auth/session";
import { useI18n } from "@/lib/i18n/useI18n";
import { dataClient } from "@/lib/services/dataClient";
import type { ProductViewModel, UserRole } from "@/lib/types";

export default function ProductsPage() {
  const { t } = useI18n();
  const session = getSession();

  const [products, setProducts] = useState<ProductViewModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const role = (session?.user.role ?? "Admin") as UserRole;

  const loadProducts = async () => {
    if (!session) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      const nextProducts = await dataClient.listProducts(session.token, role);
      setProducts(nextProducts);
    } catch (listError) {
      setError(listError instanceof Error ? listError.message : "failed to list products");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadProducts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleDelete = async (productID: string) => {
    if (!session) {
      return;
    }

    try {
      await dataClient.deleteProduct(session.token, productID);
      setProducts((prev) => prev.filter((product) => product.id !== productID));
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "failed to delete product");
    }
  };

  return (
    <section className="space-y-5">
      <header className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="font-display text-3xl font-semibold text-brand-fg">{t("productsTitle")}</h1>
          <p className="text-sm text-brand-fg/70">{role === "SuperAdmin" ? t("superAdminNote") : t("adminNote")}</p>
        </div>

        <Link
          href="/products/new"
          className="rounded-full bg-gradient-to-r from-brand-primary to-brand-accent px-4 py-2 text-sm font-semibold text-white"
        >
          {t("productsCreate")}
        </Link>
      </header>

      {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}
      {loading ? <p className="text-sm text-brand-fg/70">{t("loading")}</p> : null}

      {!loading && products.length === 0 ? <p className="text-sm text-brand-fg/70">{t("noData")}</p> : null}

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} role={role} onDelete={handleDelete} />
        ))}
      </div>
    </section>
  );
}
