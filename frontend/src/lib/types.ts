export type UserRole = "SuperAdmin" | "Admin";
export type Locale = "fr" | "en";

export interface AuthUser {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  shopID: string;
  shopName: string;
  whatsAppNumber: string;
  createdAt: string;
}

export interface AuthSession {
  token: string;
  user: AuthUser;
}

export interface JwtClaims {
  userID: string;
  email: string;
  role: UserRole;
  shopID: string;
  exp?: number;
  nbf?: number;
  iat?: number;
}

export interface CategoryViewModel {
  id: string;
  name: string;
}

export interface ProductViewModel {
  id: string;
  name: string;
  description: string;
  categoryID: string;
  category: string;
  purchasePrice?: number;
  sellingPrice: number;
  stock: number;
  imageURL: string;
  shopID: string;
  createdAt: string;
}

export interface ProductCreatePayload {
  name: string;
  description: string;
  categoryID: string;
  purchasePrice?: number;
  sellingPrice: number;
  stock: number;
  imageURL: string;
}

export interface ProductUpdatePayload extends ProductCreatePayload { }

export interface DashboardViewModel {
  totalSales: number;
  totalExpenses: number;
  netProfit: number;
  lowStockProducts: number;
  shopID: string;
}

export type TransactionType = "Sale" | "Expense" | "Withdrawal";

export interface TransactionPayload {
  type: TransactionType;
  productID?: string;
  quantity?: number;
  amount: number;
}

export interface PublicProductViewModel {
  id: string;
  name: string;
  description: string;
  category: string;
  sellingPrice: number;
  stock: number;
  imageURL: string;
  whatsappLink: string;
}
