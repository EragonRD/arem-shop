import type { ReactNode } from "react";

import { Sidebar } from "@/components/layout/Sidebar";
import { TopNav } from "@/components/layout/TopNav";

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-brand-bg">
      <div className="pointer-events-none fixed inset-0 -z-10 bg-[radial-gradient(circle_at_20%_10%,rgba(97,87,234,0.18),transparent_35%),radial-gradient(circle_at_80%_30%,rgba(243,98,200,0.16),transparent_40%),radial-gradient(circle_at_60%_90%,rgba(97,87,234,0.1),transparent_30%)]" />
      <div className="flex min-h-screen">
        <Sidebar />
        <div className="flex min-h-screen w-full flex-col">
          <TopNav />
          <main className="w-full flex-1 px-4 py-6 lg:px-8">{children}</main>
        </div>
      </div>
    </div>
  );
}
