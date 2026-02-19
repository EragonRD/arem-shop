import { apiClient } from "@/lib/services/api";
import { mockApiClient } from "@/lib/services/mockApi";

export const dataMode = process.env.NEXT_PUBLIC_DATA_MODE === "api" ? "api" : "mock";
export const dataClient = dataMode === "api" ? apiClient : mockApiClient;
