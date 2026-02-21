"use client";

import React, { useEffect, useState } from "react";

import { useI18n } from "@/lib/i18n/useI18n";
import type { ProductViewModel, TransactionPayload, TransactionType } from "@/lib/types";

interface TransactionFormProps {
  products: ProductViewModel[];
  onSubmit: (payload: TransactionPayload) => Promise<void>;
}

export function TransactionForm({ products, onSubmit }: TransactionFormProps) {
  const { t } = useI18n();
  const [type, setType] = useState<TransactionType>("Sale");
  const [productID, setProductID] = useState("");
  const [quantity, setQuantity] = useState(1);
  const [amount, setAmount] = useState(0);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  // Initialize productID when products become available
  useEffect(() => {
    if (!productID && products.length > 0) {
      setProductID(products[0].id);
    }
  }, [products, productID]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError("");

    if (amount <= 0) {
      setError("amount must be greater than 0");
      return;
    }

    if (type === "Sale") {
      if (!productID) {
        setError("product is required for sale");
        return;
      }
      if (quantity <= 0) {
        setError("quantity must be greater than 0");
        return;
      }
    }

    try {
      setSubmitting(true);
      await onSubmit({
        type,
        amount,
        ...(type === "Sale" ? { productID, quantity } : {}),
      });
      setAmount(0);
      setQuantity(1);
    } catch (submitError) {
      setError(submitError instanceof Error ? submitError.message : "unexpected error");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 rounded-2xl border border-brand-border/70 bg-white/90 p-6 shadow-soft">
      <div className="grid gap-4 md:grid-cols-2">
        <label className="space-y-1 text-sm">
          <span className="text-brand-fg/80">{t("type")}</span>
          <select
            value={type}
            onChange={(event) => setType(event.target.value as TransactionType)}
            className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
          >
            <option value="Sale">Sale</option>
            <option value="Expense">Expense</option>
            <option value="Withdrawal">Withdrawal</option>
          </select>
        </label>

        <label className="space-y-1 text-sm">
          <span className="text-brand-fg/80">{t("amount")}</span>
          <input
            type="number"
            min="0"
            step="0.01"
            value={amount}
            onChange={(event) => setAmount(Number(event.target.value))}
            className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
          />
        </label>

        {type === "Sale" ? (
          <>
            <label className="space-y-1 text-sm">
              <span className="text-brand-fg/80">{t("product")}</span>
              <select
                value={productID}
                onChange={(event) => setProductID(event.target.value)}
                className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
              >
                {products.map((product) => (
                  <option key={product.id} value={product.id}>
                    {product.name}
                  </option>
                ))}
              </select>
            </label>

            <label className="space-y-1 text-sm">
              <span className="text-brand-fg/80">{t("quantity")}</span>
              <input
                type="number"
                min="1"
                value={quantity}
                onChange={(event) => setQuantity(Number(event.target.value))}
                className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
              />
            </label>
          </>
        ) : null}
      </div>

      {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}

      <button
        type="submit"
        disabled={submitting}
        className="rounded-full bg-gradient-to-r from-brand-primary to-brand-accent px-5 py-2 text-sm font-semibold text-white transition hover:opacity-95 disabled:opacity-60"
      >
        {submitting ? t("loading") : t("create")}
      </button>
    </form>
  );
}
