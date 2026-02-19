"use client";

import Link from "next/link";

import type { ProductViewModel, UserRole } from "@/lib/types";

interface ProductCardProps {
  product: ProductViewModel;
  role: UserRole;
  onDelete?: (id: string) => void;
}

function formatPrice(value: number): string {
  return new Intl.NumberFormat("fr-FR", { style: "currency", currency: "EUR" }).format(value);
}

export function ProductCard({ product, role, onDelete }: ProductCardProps) {
  return (
    <article className="overflow-hidden rounded-2xl border border-brand-border/70 bg-white/90 shadow-soft transition hover:-translate-y-0.5">
      <div className="relative h-44 w-full overflow-hidden bg-brand-secondary/35">
        {product.imageURL ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={product.imageURL}
            alt={product.name}
            className="h-full w-full object-cover transition duration-300 hover:scale-105"
          />
        ) : (
          <div className="flex h-full items-center justify-center text-sm text-brand-fg/60">No image</div>
        )}
      </div>

      <div className="space-y-3 p-4">
        <div>
          <h3 className="font-display text-lg font-semibold text-brand-fg">{product.name}</h3>
          <p className="text-sm text-brand-fg/75">{product.description}</p>
        </div>

        <div className="grid grid-cols-2 gap-2 text-sm">
          <p className="rounded-lg bg-brand-secondary/35 px-2 py-1 text-brand-fg/80">{product.category}</p>
          <p className="rounded-lg bg-brand-secondary/35 px-2 py-1 text-brand-fg/80">Stock: {product.stock}</p>
          <p className="rounded-lg bg-brand-secondary/35 px-2 py-1 text-brand-fg/80">Sale: {formatPrice(product.sellingPrice)}</p>
          {role === "SuperAdmin" && product.purchasePrice !== undefined ? (
            <p className="rounded-lg bg-brand-secondary/35 px-2 py-1 text-brand-fg/80">Buy: {formatPrice(product.purchasePrice)}</p>
          ) : null}
        </div>

        <div className="flex items-center gap-2">
          <Link
            href={`/products/${product.id}/edit`}
            className="rounded-full border border-brand-border px-3 py-1.5 text-sm font-medium text-brand-fg transition hover:border-brand-primary hover:text-brand-primary"
          >
            Edit
          </Link>
          <button
            type="button"
            onClick={() => onDelete?.(product.id)}
            className="rounded-full border border-red-200 px-3 py-1.5 text-sm font-medium text-red-600 transition hover:bg-red-50"
          >
            Delete
          </button>
        </div>
      </div>
    </article>
  );
}
