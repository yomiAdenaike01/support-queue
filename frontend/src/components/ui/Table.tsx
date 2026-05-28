import type { PropsWithChildren } from "react";

export function Table({ children }: PropsWithChildren) {
  return (
    <div className="overflow-x-auto rounded-xl border border-[--color-border]">
      <table className="min-w-full divide-y divide-slate-800 text-sm">{children}</table>
    </div>
  );
}
