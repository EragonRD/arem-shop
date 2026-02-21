"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { ProductForm, type ProductFormValues } from "@/components/products/ProductForm";
import { getSession } from "@/lib/auth/session";
import { useI18n } from "@/lib/i18n/useI18n";
import { dataClient } from "@/lib/services/dataClient";

interface EditProductPageProps {
  params: {
    id: string;
  };
}

export default function EditProductPage({ params }: EditProductPageProps) {
  const { t } = useI18n();
  const router = useRouter();
  const session = getSession();

  const [initialValues, setInitialValues] = useState<ProductFormValues | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      if (!session) {
        return;
      }

      try {
        const product = await dataClient.getProductByID(session.token, params.id, session.user.role);
        setInitialValues({
          name: product.name,
          description: product.description,
          categoryID: product.categoryID,
          purchasePrice: product.purchasePrice,
          sellingPrice: product.sellingPrice,
          stock: product.stock,
          imageURL: product.imageURL,
        });
      } catch (loadError) {
        setError(loadError instanceof Error ? loadError.message : "failed to load product");
      }
    };

    void load();
  }, [params.id, session]);

  if (!session) {
    return null;
  }

  if (!initialValues) {
    return (
      <section className="space-y-4">
        <h1 className="font-display text-3xl font-semibold text-brand-fg">{t("productsEdit")}</h1>
        {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}
        {!error ? <p className="text-sm text-brand-fg/70">{t("loading")}</p> : null}
      </section>
    );
  }

  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="font-display text-3xl font-semibold text-brand-fg">{t("productsEdit")}</h1>
        <Link href="/products" className="text-sm text-brand-primary underline underline-offset-2">
          {t("backToProducts")}
        </Link>
      </div>

      {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}

      <ProductForm
        token={session.token}
        role={session.user.role}
        mode="edit"
        initialValues={initialValues}
        onSubmit={async (payload) => {
          await dataClient.updateProduct(session.token, params.id, session.user.role, payload);
          router.push("/products");
        }}
      />
    </section>
  );
}
