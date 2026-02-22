"use client";

import React, { useEffect, useMemo, useState } from "react";

import { useI18n } from "@/lib/i18n/useI18n";
import type { CategoryViewModel, ProductCreatePayload, ProductUpdatePayload, UserRole } from "@/lib/types";
import { apiClient } from "@/lib/services/api";

export interface ProductFormValues {
  name: string;
  description: string;
  categoryID: string;
  purchasePrice?: number;
  sellingPrice: number;
  stock: number;
  imageURL: string;
}

interface ProductFormProps {
  token: string;
  role: UserRole;
  mode: "create" | "edit";
  initialValues?: ProductFormValues;
  onSubmit: (payload: ProductCreatePayload | ProductUpdatePayload) => Promise<void>;
}

const defaultValues: ProductFormValues = {
  name: "",
  description: "",
  categoryID: "",
  sellingPrice: 0,
  stock: 0,
  imageURL: "",
};

export function ProductForm({ token, role, mode, initialValues, onSubmit }: ProductFormProps) {
  const { t } = useI18n();
  const [values, setValues] = useState<ProductFormValues>(initialValues ?? defaultValues);
  const [submitting, setSubmitting] = useState(false);
  const [categories, setCategories] = useState<CategoryViewModel[]>([]);
  const [error, setError] = useState("");

  const canEditPurchasePrice = role === "SuperAdmin";

  const submitLabel = useMemo(() => {
    if (mode === "create") {
      return t("create");
    }
    return t("update");
  }, [mode, t]);

  useEffect(() => {
    if (!token) return;

    apiClient
      .listCategories(token)
      .then((data) => {
        setCategories(data);
        if (!values.categoryID && data.length > 0) {
          setValues((prev) => ({ ...prev, categoryID: data[0].id }));
        }
      })
      .catch((err) => {
        console.error("Failed to load categories", err);
        setError("Could not load categories. Please try again later.");
      });
  }, [token, values.categoryID]);

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

  const [uploadingImage, setUploadingImage] = useState(false);

  const handleImageUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file || !token) return;

    try {
      setUploadingImage(true);
      setError("");
      // Call the API client directly
      const { url } = await apiClient.uploadImage(token, file);
      handleChange("imageURL", url);
    } catch (err: any) {
      setError(err?.message || "Failed to upload image");
    } finally {
      setUploadingImage(false);
      // Reset the file input so the same file could be selected again if needed
      event.target.value = "";
    }
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError("");

    if (!values.name.trim() || !values.categoryID.trim()) {
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
        categoryID: values.categoryID.trim(),
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
        <label className="space-y-1 text-sm md:col-span-2">
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
          <select
            value={values.categoryID}
            onChange={(event) => handleChange("categoryID", event.target.value)}
            className="w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring bg-white"
            required
            disabled={categories.length === 0}
          >
            <option value="" disabled>
              {categories.length === 0 ? "Loading categories..." : "Select a category"}
            </option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </label>

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

        <label className="space-y-1 text-sm md:col-span-2">
          <span className="text-brand-fg/80">{t("description")}</span>
          <textarea
            value={values.description}
            onChange={(event) => handleChange("description", event.target.value)}
            className="h-24 w-full rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
          />
        </label>

        <div className="space-y-1 text-sm md:col-span-2">
          <span className="text-brand-fg/80">{t("imageURL") || "Image"}</span>
          <div className="flex gap-2 items-center">
            {values.imageURL && (
              <img
                src={values.imageURL}
                alt="Product preview"
                className="h-12 w-12 rounded-lg object-cover border border-brand-border bg-brand-secondary/30 flex-shrink-0"
                onError={(e) => (e.currentTarget.style.display = "none")}
                onLoad={(e) => (e.currentTarget.style.display = "block")}
              />
            )}
            <input
              value={values.imageURL}
              onChange={(event) => handleChange("imageURL", event.target.value)}
              className="flex-1 rounded-xl border border-brand-border px-3 py-2 outline-none ring-brand-primary/40 focus:ring"
              placeholder="https://..."
            />
            <label className={`shrink-0 cursor-pointer rounded-xl bg-brand-secondary/50 px-4 py-2 text-sm font-medium text-brand-fg/80 hover:bg-brand-secondary transition border border-brand-border ${uploadingImage ? 'opacity-50 pointer-events-none' : ''}`}>
              {uploadingImage ? "Uploading..." : "Upload local file"}
              <input
                type="file"
                accept="image/jpeg,image/png,image/webp"
                className="hidden"
                disabled={uploadingImage}
                onChange={handleImageUpload}
              />
            </label>
          </div>
        </div>
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
