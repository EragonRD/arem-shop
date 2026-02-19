import type { ReactNode } from "react";

interface StatCardProps {
  label: string;
  value: string;
  icon?: ReactNode;
}

export function StatCard({ label, value, icon }: StatCardProps) {
  return (
    <article className="group rounded-2xl border border-brand-border/70 bg-white/85 p-5 shadow-soft backdrop-blur transition hover:-translate-y-0.5">
      <div className="mb-4 flex items-center justify-between">
        <span className="text-sm font-medium text-brand-fg/75">{label}</span>
        {icon ? <span className="text-brand-primary">{icon}</span> : null}
      </div>
      <p className="font-display text-2xl font-semibold text-brand-fg">{value}</p>
    </article>
  );
}
