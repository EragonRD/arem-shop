import React from "react";
import { render, screen } from "@testing-library/react";

import { ProductCard } from "@/components/products/ProductCard";
import type { ProductViewModel } from "@/lib/types";

const productFixture: ProductViewModel = {
  id: "p-1",
  name: "Arem Laptop",
  description: "Demo",
  category: "Laptops",
  purchasePrice: 800,
  sellingPrice: 999.99,
  stock: 3,
  imageURL: "",
  shopID: "11111111-1111-1111-1111-111111111111",
  createdAt: "2026-02-19T00:00:00Z",
};

describe("ProductCard", () => {
  it("hides purchase price for Admin", () => {
    render(<ProductCard product={productFixture} role="Admin" />);

    expect(screen.queryByText(/Buy:/i)).not.toBeInTheDocument();
  });

  it("shows purchase price for SuperAdmin", () => {
    render(<ProductCard product={productFixture} role="SuperAdmin" />);

    expect(screen.getByText(/Buy:/i)).toBeInTheDocument();
  });
});
