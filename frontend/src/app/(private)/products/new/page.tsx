"use client";

import { useRouter } from "next/navigation";

import { ProductForm } from "@/components/products/ProductForm";
import { getSession } from "@/lib/auth/session";
import { useI18n } from "@/lib/i18n/useI18n";
import { dataClient } from "@/lib/services/dataClient";

export default function NewProductPage() {
  const router = useRouter();
  const { t } = useI18n();
  const session = getSession();

  if (!session) {
    return null;
  }

  return (
    <section className="space-y-4">
      <h1 className="font-display text-3xl font-semibold text-brand-fg">{t("productsCreate")}</h1>
      <ProductForm
        role={session.user.role}
        mode="create"
        onSubmit={async (payload) => {
          await dataClient.createProduct(session.token, session.user.role, payload);
          router.push("/products");
        }}
      />
    </section>
  );
}
