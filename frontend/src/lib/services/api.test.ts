import { beforeEach, describe, expect, it, vi } from "vitest";

const { handleUnauthorizedMock } = vi.hoisted(() => {
  return {
    handleUnauthorizedMock: vi.fn(),
  };
});

vi.mock("@/lib/auth/guards", () => {
  return {
    handleUnauthorized: handleUnauthorizedMock,
  };
});

import { UnauthorizedError, apiClient } from "@/lib/services/api";

describe("apiClient", () => {
  beforeEach(() => {
    handleUnauthorizedMock.mockReset();
    vi.restoreAllMocks();
  });

  it("maps product envelope payload correctly", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce(
      new Response(
        JSON.stringify({
          success: true,
          data: [
            {
              id: "p-1",
              name: "Laptop",
              description: "Demo",
              category: "Laptops",
              purchasePrice: "800.00",
              sellingPrice: "999.99",
              stock: 10,
              imageURL: "https://example.com/laptop.jpg",
              shopID: "11111111-1111-1111-1111-111111111111",
              createdAt: "2026-02-19T00:00:00Z",
            },
          ],
        }),
        { status: 200 },
      ),
    );

    const products = await apiClient.listProducts("token", "SuperAdmin");

    expect(products).toHaveLength(1);
    expect(products[0].sellingPrice).toBe(999.99);
    expect(products[0].purchasePrice).toBe(800);
  });

  it("triggers unauthorized handler on 401", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValueOnce(
      new Response(JSON.stringify({ error: "invalid token" }), { status: 401 }),
    );

    await expect(apiClient.listProducts("token", "SuperAdmin")).rejects.toBeInstanceOf(UnauthorizedError);
    expect(handleUnauthorizedMock).toHaveBeenCalledTimes(1);
  });
});
