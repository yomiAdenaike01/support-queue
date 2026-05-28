import { Bell, Search } from "lucide-react";
import { useUiStore } from "@/store/ui";

export function Header() {
  const search = useUiStore((state) => state.search);
  const setSearch = useUiStore((state) => state.setSearch);
  return (
    <header className="sticky top-0 z-30 border-b border-slate-800 bg-[--color-bg]/90 px-4 py-4 backdrop-blur lg:px-8">
      <div className="flex items-center gap-4">
        <div className="relative max-w-xl flex-1">
          <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500" />
          <input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Search tickets and customers"
            className="w-full rounded-lg border border-slate-800 bg-surface py-2 pl-10 pr-3 text-sm outline-none focus:border-accent"
          />
        </div>
        <div className="hidden items-center gap-2 rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-200 sm:flex">
          <span className="h-2 w-2 rounded-full bg-emerald-400" />
          API Online
        </div>
        <button className="rounded-lg bg-surface-2 p-2 text-slate-300 hover:text-white" aria-label="Notifications">
          <Bell className="h-5 w-5" />
        </button>
      </div>
    </header>
  );
}
