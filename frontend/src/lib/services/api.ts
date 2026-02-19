import { handleUnauthorized } from "@/lib/auth/guards";
import {
  toDashboardViewModel,
  toProductPayload,
  toProductViewModel,
  toPublicProductViewModel,
  toTransactionPayload,
} from "@/lib/services/adapters";
import type { DataClient } from "@/lib/services/types";
import type {
  AuthSession,
  DashboardViewModel,
  ProductCreatePayload,
  ProductUpdatePayload,
  ProductViewModel,
  PublicProductViewModel,
  TransactionPayload,
  UserRole,
} from "@/lib/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type ProductEnvelope<T> = {
  success: boolean;
  data?: T;
  error?: string;
};

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export class UnauthorizedError extends ApiError {
  constructor(message = "unauthorized") {
    super(401, message);
  }
}

function extractErrorMessage(payload: unknown, fallback: string): string {
  if (!payload || typeof payload !== "object") {
    return fallback;
  }

  const record = payload as Record<string, unknown>;
  if (typeof record.error === "string" && record.error.trim() !== "") {
    return record.error;
  }

  return fallback;
}

async function request<T>(path: string, init: RequestInit = {}, token?: string): Promise<T> {
  const headers = new Headers(init.headers);
  if (!headers.has("Content-Type") && init.body) {
    headers.set("Content-Type", "application/json");
  }
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers,
    cache: "no-store",
  });

  const text = await response.text();
  let parsed: unknown = null;
  if (text) {
    try {
      parsed = JSON.parse(text) as unknown;
    } catch {
      parsed = null;
    }
  }

  if (response.status === 401) {
    handleUnauthorized();
    if (typeof window !== "undefined") {
      window.dispatchEvent(new Event("arem:unauthorized"));
    }
    throw new UnauthorizedError(extractErrorMessage(parsed, "unauthorized"));
  }

  if (!response.ok) {
    throw new ApiError(response.status, extractErrorMessage(parsed, `request failed with status ${response.status}`));
  }

  return parsed as T;
}

function parseEnvelope<T>(payload: ProductEnvelope<T>): T {
  if (!payload.success || payload.data === undefined) {
    throw new ApiError(500, payload.error ?? "invalid envelope response");
  }

  return payload.data;
}

export const apiClient: DataClient = {
  async login(email: string, password: string, shopID: string): Promise<AuthSession> {
    return request<AuthSession>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password, shopID }),
    });
  },

  async getDashboard(token: string): Promise<DashboardViewModel> {
    const payload = await request<DashboardViewModel>("/reports/dashboard", { method: "GET" }, token);
    return toDashboardViewModel(payload);
  },

  async listProducts(token: string, _role: UserRole): Promise<ProductViewModel[]> {
    const payload = await request<ProductEnvelope<unknown[]>>("/products", { method: "GET" }, token);
    const data = parseEnvelope(payload);
    return data.map((product) => toProductViewModel(product as never));
  },

  async getProductByID(token: string, productID: string, _role: UserRole): Promise<ProductViewModel> {
    const payload = await request<ProductEnvelope<unknown>>(`/products/${productID}`, { method: "GET" }, token);
    const data = parseEnvelope(payload);
    return toProductViewModel(data as never);
  },

  async createProduct(
    token: string,
    role: UserRole,
    input: ProductCreatePayload,
  ): Promise<ProductViewModel> {
    const payload = await request<ProductEnvelope<unknown>>(
      "/products",
      {
        method: "POST",
        body: JSON.stringify(toProductPayload(role, input)),
      },
      token,
    );

    return toProductViewModel(parseEnvelope(payload) as never);
  },

  async updateProduct(
    token: string,
    productID: string,
    role: UserRole,
    input: ProductUpdatePayload,
  ): Promise<ProductViewModel> {
    const payload = await request<ProductEnvelope<unknown>>(
      `/products/${productID}`,
      {
        method: "PUT",
        body: JSON.stringify(toProductPayload(role, input)),
      },
      token,
    );

    return toProductViewModel(parseEnvelope(payload) as never);
  },

  async deleteProduct(token: string, productID: string): Promise<void> {
    const payload = await request<ProductEnvelope<unknown>>(`/products/${productID}`, { method: "DELETE" }, token);
    parseEnvelope(payload);
  },

  async createTransaction(token: string, payload: TransactionPayload): Promise<void> {
    await request<unknown>(
      "/transactions",
      {
        method: "POST",
        body: JSON.stringify(toTransactionPayload(payload)),
      },
      token,
    );
  },

  async listPublicProducts(shopID: string): Promise<PublicProductViewModel[]> {
    const payload = await request<unknown[]>(`/public/${shopID}/products`, { method: "GET" });
    return payload.map((product) => toPublicProductViewModel(product as never));
  },
};
