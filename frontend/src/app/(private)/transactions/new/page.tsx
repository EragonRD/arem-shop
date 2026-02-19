"use client";

import { useEffect, useState } from "react";

import { TransactionForm } from "@/components/transactions/TransactionForm";
import { getSession } from "@/lib/auth/session";
import { useI18n } from "@/lib/i18n/useI18n";
import { dataClient } from "@/lib/services/dataClient";
import type { ProductViewModel } from "@/lib/types";

export default function NewTransactionPage() {
  const { t } = useI18n();
  const session = getSession();

  const [products, setProducts] = useState<ProductViewModel[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      if (!session) {
        return;
      }

      try {
        const list = await dataClient.listProducts(session.token, session.user.role);
        setProducts(list);
      } catch (loadError) {
        setError(loadError instanceof Error ? loadError.message : "failed to load products");
      }
    };

    void load();
  }, [session]);

  if (!session) {
    return null;
  }

  return (
    <section className="space-y-4">
      <h1 className="font-display text-3xl font-semibold text-brand-fg">{t("transactionsTitle")}</h1>
      {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}

      <TransactionForm
        products={products}
        onSubmit={async (payload) => {
          await dataClient.createTransaction(session.token, payload);
        }}
      />
    </section>
  );
}
