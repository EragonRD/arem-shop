import { apiClient } from "@/lib/services/api";
import { mockApiClient } from "@/lib/services/mockApi";

const requestedMode = (process.env.NEXT_PUBLIC_DATA_MODE ?? "api").toLowerCase();
export const dataMode = requestedMode === "mock" ? "mock" : "api";
export const dataClient = dataMode === "api" ? apiClient : mockApiClient;
