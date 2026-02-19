import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";

import { I18nProvider } from "@/lib/i18n/I18nProvider";
import PublicShopPage from "./page";

const listPublicProductsMock = vi.fn();

vi.mock("@/lib/services/dataClient", () => {
  return {
    dataClient: {
      listPublicProducts: (...args: unknown[]) => listPublicProductsMock(...args),
    },
  };
});

describe("PublicShopPage", () => {
  it("renders public catalog without auth", async () => {
    listPublicProductsMock.mockResolvedValueOnce([
      {
        id: "p-1",
        name: "Public Laptop",
        description: "Demo",
        category: "Laptops",
        sellingPrice: 1200,
        stock: 7,
        imageURL: "",
        whatsappLink: "https://wa.me/212600000000",
      },
    ]);

    render(
      <I18nProvider>
        <PublicShopPage params={{ shopID: "11111111-1111-1111-1111-111111111111" }} />
      </I18nProvider>,
    );

    await waitFor(() => {
      expect(screen.getByText("Public Laptop")).toBeInTheDocument();
    });

    expect(screen.getByText(/Catalogue public|Public catalog/i)).toBeInTheDocument();
  });
});
