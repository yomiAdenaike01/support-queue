import { Inbox } from "lucide-react";

export function EmptyState({ title, message }: { title: string; message: string }) {
  return (
    <div className="flex min-h-48 flex-col items-center justify-center rounded-xl border border-dashed border-slate-700 p-8 text-center">
      <Inbox className="mb-3 h-8 w-8 text-slate-500" />
      <h3 className="text-base font-semibold text-slate-100">{title}</h3>
      <p className="mt-1 max-w-md text-sm text-slate-400">{message}</p>
    </div>
  );
}
