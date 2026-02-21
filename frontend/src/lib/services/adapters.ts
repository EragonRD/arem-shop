import type {
  DashboardViewModel,
  ProductCreatePayload,
  ProductUpdatePayload,
  ProductViewModel,
  PublicProductViewModel,
  TransactionPayload,
  UserRole,
} from "@/lib/types";

type ApiProduct = {
  id: string;
  name: string;
  description: string;
  categoryID: string;
  category: string;
  purchasePrice?: string | number;
  sellingPrice: string | number;
  stock: number;
  imageURL: string;
  shopID: string;
  createdAt: string;
};

type ApiPublicProduct = {
  id: string;
  name: string;
  description: string;
  category: string;
  sellingPrice: string | number;
  stock: number;
  imageURL: string;
  whatsappLink: string;
};

export function toNumber(value: number | string | undefined): number {
  if (typeof value === "number") {
    return value;
  }
  if (!value) {
    return 0;
  }

  const parsed = Number(value);
  if (Number.isNaN(parsed)) {
    return 0;
  }
  return parsed;
}

export function toProductViewModel(product: ApiProduct): ProductViewModel {
  return {
    id: product.id,
    name: product.name,
    description: product.description,
    categoryID: product.categoryID,
    category: product.category,
    purchasePrice: product.purchasePrice !== undefined ? toNumber(product.purchasePrice) : undefined,
    sellingPrice: toNumber(product.sellingPrice),
    stock: product.stock,
    imageURL: product.imageURL,
    shopID: product.shopID,
    createdAt: product.createdAt,
  };
}

export function toPublicProductViewModel(product: ApiPublicProduct): PublicProductViewModel {
  return {
    id: product.id,
    name: product.name,
    description: product.description,
    category: product.category,
    sellingPrice: toNumber(product.sellingPrice),
    stock: product.stock,
    imageURL: product.imageURL,
    whatsappLink: product.whatsappLink,
  };
}

export function toDashboardViewModel(raw: DashboardViewModel): DashboardViewModel {
  return {
    totalSales: toNumber(raw.totalSales),
    totalExpenses: toNumber(raw.totalExpenses),
    netProfit: toNumber(raw.netProfit),
    lowStockProducts: raw.lowStockProducts,
    shopID: raw.shopID,
  };
}

export function toProductPayload(role: UserRole, payload: ProductCreatePayload | ProductUpdatePayload) {
  const basePayload: Record<string, unknown> = {
    name: payload.name,
    description: payload.description,
    categoryID: payload.categoryID,
    sellingPrice: payload.sellingPrice,
    stock: payload.stock,
    imageURL: payload.imageURL,
  };

  if (role === "SuperAdmin" && payload.purchasePrice !== undefined) {
    basePayload.purchasePrice = payload.purchasePrice;
  }

  return basePayload;
}

export function toTransactionPayload(payload: TransactionPayload) {
  const body: Record<string, unknown> = {
    type: payload.type,
    amount: payload.amount,
  };

  if (payload.type === "Sale") {
    body.productID = payload.productID;
    body.quantity = payload.quantity;
  }

  return body;
}
