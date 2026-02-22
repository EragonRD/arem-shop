import { mockDashboard } from "@/lib/mocks/dashboard";
import { mockProductsSeed } from "@/lib/mocks/products";
import { mockTransactions } from "@/lib/mocks/transactions";
import type {
  AuthSession,
  DashboardViewModel,
  ProductCreatePayload,
  ProductUpdatePayload,
  ProductViewModel,
  PublicProductViewModel,
  TransactionPayload,
  UserRole,
  CategoryViewModel,
} from "@/lib/types";
import type { DataClient } from "@/lib/services/types";

const DEFAULT_SHOP_ID = "11111111-1111-1111-1111-111111111111";

let products: ProductViewModel[] = structuredClone(mockProductsSeed);

function delay(ms = 160): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

function toBase64Url(raw: string): string {
  if (typeof window !== "undefined" && typeof window.btoa === "function") {
    return window.btoa(raw).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
  }

  return Buffer.from(raw).toString("base64url");
}

function createMockToken(role: UserRole, shopID: string, email: string): string {
  const header = toBase64Url(JSON.stringify({ alg: "HS256", typ: "JWT" }));
  const payload = toBase64Url(
    JSON.stringify({
      userID: "mock-user-id",
      email,
      role,
      shopID,
      exp: Math.floor(Date.now() / 1000) + 60 * 60,
    }),
  );

  return `${header}.${payload}.mock-signature`;
}

function roleFromEmail(email: string): UserRole {
  const normalized = email.trim().toLowerCase();
  if (normalized.startsWith("admin") || normalized.includes("+admin") || normalized.includes("admin@")) {
    return "Admin";
  }
  return "SuperAdmin";
}

function validateProductPayload(
  role: UserRole,
  payload: ProductCreatePayload | ProductUpdatePayload,
  mode: "create" | "update",
): void {
  if (!payload.name.trim()) {
    throw new Error("name is required");
  }
  if (!payload.categoryID.trim()) {
    throw new Error("category is required");
  }
  if (payload.sellingPrice <= 0) {
    throw new Error("sellingPrice must be greater than 0");
  }
  if (payload.stock < 0) {
    throw new Error("stock cannot be negative");
  }

  if (role === "Admin") {
    if (payload.purchasePrice !== undefined) {
      throw new Error("purchasePrice is not allowed for Admin");
    }
    return;
  }

  if (mode === "create" && payload.purchasePrice === undefined) {
    throw new Error("purchasePrice is required for SuperAdmin");
  }

  if (payload.purchasePrice !== undefined && payload.purchasePrice <= 0) {
    throw new Error("purchasePrice must be greater than 0");
  }
}

function maskProductByRole(product: ProductViewModel, role: UserRole): ProductViewModel {
  if (role === "SuperAdmin") {
    return { ...product };
  }

  const { purchasePrice: _ignored, ...safeProduct } = product;
  return safeProduct;
}

function computeDashboard(): DashboardViewModel {
  const totalSales = mockTransactions
    .filter((transaction) => transaction.type === "Sale")
    .reduce((sum, transaction) => sum + transaction.amount, 0);

  const totalExpenses = mockTransactions
    .filter((transaction) => transaction.type === "Expense" || transaction.type === "Withdrawal")
    .reduce((sum, transaction) => sum + transaction.amount, 0);

  const lowStockProducts = products.filter((product) => product.stock <= 5).length;

  return {
    ...mockDashboard,
    totalSales,
    totalExpenses,
    netProfit: totalSales - totalExpenses,
    lowStockProducts,
  };
}

export const mockApiClient: DataClient = {
  async login(email: string, password: string, shopID: string): Promise<AuthSession> {
    await delay();

    if (!email.trim() || !password.trim() || !shopID.trim()) {
      throw new Error("invalid credentials");
    }

    const role = roleFromEmail(email);

    return {
      token: createMockToken(role, shopID, email),
      user: {
        id: "mock-user-id",
        name: role === "SuperAdmin" ? "Owner Demo" : "Admin Demo",
        email,
        role,
        shopID,
        shopName: "Demo Shop",
        whatsAppNumber: "+212600000000",
        createdAt: new Date().toISOString(),
      },
    };
  },

  async getDashboard(_token: string): Promise<DashboardViewModel> {
    await delay();
    return computeDashboard();
  },

  async listCategories(_token: string): Promise<CategoryViewModel[]> {
    await delay();
    return [
      { id: "22222222-2222-2222-2222-222222222222", name: "Laptops" },
      { id: "33333333-3333-3333-3333-333333333333", name: "Smartphones" },
      { id: "44444444-4444-4444-4444-444444444444", name: "Audio" },
    ];
  },

  async listProducts(_token: string, role: UserRole): Promise<ProductViewModel[]> {
    await delay();
    return products.map((product) => maskProductByRole(product, role));
  },

  async getProductByID(_token: string, productID: string, role: UserRole): Promise<ProductViewModel> {
    await delay();

    const product = products.find((candidate) => candidate.id === productID);
    if (!product) {
      throw new Error("product not found");
    }

    return maskProductByRole(product, role);
  },

  async createProduct(
    _token: string,
    role: UserRole,
    payload: ProductCreatePayload,
  ): Promise<ProductViewModel> {
    await delay();

    validateProductPayload(role, payload, "create");

    const created: ProductViewModel = {
      id: `p-${Date.now()}`,
      name: payload.name.trim(),
      description: payload.description.trim(),
      category: "Mocked Category",
      categoryID: payload.categoryID.trim(),
      purchasePrice: role === "SuperAdmin" ? payload.purchasePrice : 0,
      sellingPrice: payload.sellingPrice,
      stock: payload.stock,
      imageURL: payload.imageURL.trim(),
      shopID: DEFAULT_SHOP_ID,
      createdAt: new Date().toISOString(),
    };

    products = [created, ...products];
    return maskProductByRole(created, role);
  },

  async updateProduct(
    _token: string,
    productID: string,
    role: UserRole,
    payload: ProductUpdatePayload,
  ): Promise<ProductViewModel> {
    await delay();

    validateProductPayload(role, payload, "update");

    const index = products.findIndex((candidate) => candidate.id === productID);
    if (index === -1) {
      throw new Error("product not found");
    }

    const existing = products[index];
    const updated: ProductViewModel = {
      ...existing,
      name: payload.name.trim(),
      description: payload.description.trim(),
      categoryID: payload.categoryID.trim(),
      sellingPrice: payload.sellingPrice,
      stock: payload.stock,
      imageURL: payload.imageURL.trim(),
      purchasePrice:
        role === "SuperAdmin"
          ? payload.purchasePrice !== undefined
            ? payload.purchasePrice
            : existing.purchasePrice
          : 0,
    };

    products[index] = updated;
    return maskProductByRole(updated, role);
  },

  async deleteProduct(_token: string, productID: string): Promise<void> {
    await delay();

    const nextProducts = products.filter((candidate) => candidate.id !== productID);
    if (nextProducts.length === products.length) {
      throw new Error("product not found");
    }

    products = nextProducts;
  },

  async createTransaction(_token: string, payload: TransactionPayload): Promise<void> {
    await delay();

    if (payload.amount <= 0) {
      throw new Error("amount must be greater than 0");
    }

    if (payload.type === "Sale") {
      if (!payload.productID) {
        throw new Error("productID is required for sale");
      }

      const quantity = payload.quantity ?? 0;
      if (quantity <= 0) {
        throw new Error("quantity must be greater than 0 for sale");
      }

      const index = products.findIndex((candidate) => candidate.id === payload.productID);
      if (index === -1) {
        throw new Error("product not found");
      }

      const product = products[index];
      if (product.stock < quantity) {
        throw new Error("insufficient stock");
      }

      products[index] = { ...product, stock: product.stock - quantity };
    }

    mockTransactions.push(payload);
  },

  async listPublicProducts(shopID: string): Promise<PublicProductViewModel[]> {
    await delay();
    return products
      .filter((p) => p.shopID === shopID)
      .map((product) => ({
        id: product.id,
        name: product.name,
        description: product.description,
        category: product.category,
        sellingPrice: product.sellingPrice,
        stock: product.stock,
        imageURL: product.imageURL,
        whatsappLink: `https://wa.me/212600000000?text=Bonjour%20je%20veux%20plus%20d%27information%20sur%20${encodeURIComponent(product.name)}`,
      }));
  },

  async uploadImage(_token: string, file: File): Promise<{ url: string }> {
    await delay(1500); // Simulate upload latency
    // Return a fake URL based loosely on the file name for demonstration
    return { url: `https://via.placeholder.com/800x600?text=${encodeURIComponent(file.name)}` };
  },

  async updateShop(_token: string, payload: { name: string; whatsAppNumber?: string }): Promise<void> {
    await delay();
    return new Promise((resolve) => Object.assign({}, payload, resolve()));
  },
};
