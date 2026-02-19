import type { Metadata } from "next";
import { Outfit, Space_Grotesk } from "next/font/google";

import { AppProviders } from "@/components/providers/AppProviders";

import "./globals.css";

const outfit = Outfit({
  subsets: ["latin"],
  variable: "--font-outfit",
});

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-space-grotesk",
});

export const metadata: Metadata = {
  title: "Arem Shop Frontend",
  description: "Arem Shop dashboard and public storefront",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="fr">
      <body className={`${outfit.variable} ${spaceGrotesk.variable} font-sans antialiased`}>
        <AppProviders>{children}</AppProviders>
      </body>
    </html>
  );
}
