"use client";

import { useMemo, useState } from "react";

import { useI18n } from "@/lib/i18n/useI18n";
import type { ProductCreatePayload, ProductUpdatePayload, UserRole } from "@/lib/types";

export interface ProductFormValues {
  name: string;
  description: string;
  category: string;
  purchasePrice?: number;
  sellingPrice: number;
  stock: number;
  imageURL: string;
}

interface ProductFormProps {
  role: UserRole;
  mode: "create" | "edit";
  initialValues?: ProductFormValues;
  onSubmit: (payload: ProductCreatePayload | ProductUpdatePayload) => Promise<void>;
}

const defaultValues: ProductFormValues = {
  name: "",
  description: "",
  category: "",
  sellingPrice: 0,
  stock: 0,
  imageURL: "",
};

export function ProductForm({ role, mode, initialValues, onSubmit }: ProductFormProps) {
  const { t } = useI18n();
  const [values, setValues] = useState<ProductFormValues>(initialValues ?? defaultValues);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const canEditPurchasePrice = role === "SuperAdmin";

  const submitLabel = useMemo(() => {
    if (mode === "create") {
      return t("create");
    }
    return t("update");
  }, [mode, t]);

  const handleChange = (field: keyof ProductFormValues, value: string) => {
    if (field === "sellingPrice" || field === "purchasePrice") {
      setValues((prev) => ({ ...prev, [field]: value === "" ? undefined : Number(value) }));
      return;
    }

    if (field === "stock") {
      setValues((prev) => ({ ...prev, stock: Number(value) }));
      return;
    }

    setValues((prev) => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError("");

    if (!values.name.trim() || !values.category.trim()) {
      setError("name and category are required");
      return;
    }
    if (values.sellingPrice <= 0) {
      setError("sellingPrice must be greater than 0");
      return;
    }
    if (values.stock < 0) {
      setError("stock cannot be negative");
      return;
    }
    if (canEditPurchasePrice && mode === "create" && (!values.purchasePrice || values.purchasePrice <= 0)) {
      setError("purchasePrice must be greater than 0 for SuperAdmin");
      return;
    }
    if (canEditPurchasePrice && values.purchasePrice !== undefined && values.purchasePrice <= 0) {
      setError("purchasePrice must be greater than 0");
      return;
    }

    try {
      setSubmitting(true);
      await onSubmit({
        name: values.name.trim(),
        description: values.description.trim(),
        category: values.category.trim(),
        sellingPrice: Number(values.sellingPrice),
        stock: Number(values.stock),
        imageURL: values.imageURL.trim(),
        ...(canEditPurchasePrice && values.purchasePrice !== undefined
          ? { purchasePrice: Number(values.purchasePrice) }
          : {}),
      });
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
          <span className="text-brand-fg/80">{t("name")}</span>
          <input
            value={values.name}
            onChange={(event) => handleChange("name", event.target.value)}
            className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
            required
          />
        </label>

        <label className="space-y-1 text-sm">
          <span className="text-brand-fg/80">{t("category")}</span>
          <input
            value={values.category}
            onChange={(event) => handleChange("category", event.target.value)}
            className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
            required
          />
        </label>

        <label className="space-y-1 text-sm md:col-span-2">
          <span className="text-brand-fg/80">{t("description")}</span>
          <textarea
            value={values.description}
            onChange={(event) => handleChange("description", event.target.value)}
            className="h-24 w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
          />
        </label>

        <label className="space-y-1 text-sm">
          <span className="text-brand-fg/80">{t("sellingPrice")}</span>
          <input
            type="number"
            min="0"
            step="0.01"
            value={values.sellingPrice}
            onChange={(event) => handleChange("sellingPrice", event.target.value)}
            className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
            required
          />
        </label>

        {canEditPurchasePrice ? (
          <label className="space-y-1 text-sm">
            <span className="text-brand-fg/80">{t("purchasePrice")}</span>
            <input
              type="number"
              min="0"
              step="0.01"
              value={values.purchasePrice ?? ""}
              onChange={(event) => handleChange("purchasePrice", event.target.value)}
              className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
            />
          </label>
        ) : (
          <div className="rounded-xl border border-dashed border-brand-border bg-brand-secondary/35 p-3 text-xs text-brand-fg/70">
            {t("adminNote")}
          </div>
        )}

        <label className="space-y-1 text-sm">
          <span className="text-brand-fg/80">{t("stock")}</span>
          <input
            type="number"
            min="0"
            value={values.stock}
            onChange={(event) => handleChange("stock", event.target.value)}
            className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
            required
          />
        </label>

        <label className="space-y-1 text-sm">
          <span className="text-brand-fg/80">{t("imageURL")}</span>
          <input
            value={values.imageURL}
            onChange={(event) => handleChange("imageURL", event.target.value)}
            className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
          />
        </label>
      </div>

      {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}

      <button
        type="submit"
        disabled={submitting}
        className="rounded-full bg-gradient-to-r from-brand-primary to-brand-accent px-5 py-2 text-sm font-semibold text-white transition hover:opacity-95 disabled:opacity-60"
      >
        {submitting ? `${t("loading")}` : submitLabel}
      </button>
    </form>
  );
}
