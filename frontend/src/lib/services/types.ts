import type {
  DashboardViewModel,
  ProductCreatePayload,
  ProductUpdatePayload,
  ProductViewModel,
  PublicProductViewModel,
  TransactionPayload,
  UserRole,
  AuthSession,
} from "@/lib/types";

export interface DataClient {
  login: (email: string, password: string, shopID: string) => Promise<AuthSession>;
  getDashboard: (token: string) => Promise<DashboardViewModel>;
  listProducts: (token: string, role: UserRole) => Promise<ProductViewModel[]>;
  getProductByID: (token: string, productID: string, role: UserRole) => Promise<ProductViewModel>;
  createProduct: (token: string, role: UserRole, payload: ProductCreatePayload) => Promise<ProductViewModel>;
  updateProduct: (
    token: string,
    productID: string,
    role: UserRole,
    payload: ProductUpdatePayload,
  ) => Promise<ProductViewModel>;
  deleteProduct: (token: string, productID: string) => Promise<void>;
  createTransaction: (token: string, payload: TransactionPayload) => Promise<void>;
  listPublicProducts: (shopID: string) => Promise<PublicProductViewModel[]>;
}
