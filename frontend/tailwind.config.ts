import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/app/**/*.{js,ts,jsx,tsx}",
    "./src/components/**/*.{js,ts,jsx,tsx}",
    "./src/lib/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          primary: "var(--primary)",
          accent: "var(--accent)",
          secondary: "var(--secondary)",
          bg: "var(--background)",
          fg: "var(--foreground)",
          muted: "var(--muted)",
          border: "var(--border)",
          darkPrimary: "var(--dark-primary)",
          darkAccent: "var(--dark-accent)",
        },
      },
      borderRadius: {
        xl: "var(--radius)",
      },
      boxShadow: {
        soft: "0 18px 45px -24px rgba(34, 18, 73, 0.45)",
      },
    },
  },
  plugins: [],
};

export default config;
