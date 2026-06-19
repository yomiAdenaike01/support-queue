import {
  LayoutDashboard,
  ListPlus,
  Play,
  Plus,
  Settings,
  Ticket,
  Users,
} from "lucide-react";
import { NavLink } from "react-router-dom";
import { useMetrics } from "@/hooks/useMetrics";

const navItems = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/tickets", label: "Tickets", icon: Ticket },
  { to: "/tickets/new", label: "Create Ticket", icon: Plus },
  { to: "/teams", label: "Teams", icon: Users },
  { to: "/input-sources", label: "Input Sources", icon: ListPlus },
  { to: "/simulator", label: "Simulator", icon: Play },
  { to: "/settings", label: "Settings", icon: Settings },
];

export function Sidebar() {
  const { data } = useMetrics();
  return (
    <aside className="fixed inset-y-0 left-0 z-40 flex w-20 flex-col border-r border-slate-800 bg-surface px-3 py-5 lg:w-64">
      <div className="mb-8 flex items-center gap-3 px-2">
        <div className="grid h-10 w-10 place-items-center rounded-lg bg-accent font-mono font-bold">
          SO
        </div>
        <span className="hidden text-lg font-semibold lg:block">
          SupportOps
        </span>
      </div>
      <nav className="space-y-1">
        {navItems.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === "/"}
              className={({ isActive }) =>
                `flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition ${
                  isActive
                    ? "bg-accent text-white"
                    : "text-slate-400 hover:bg-surface-2 hover:text-white"
                }`
              }
            >
              <Icon className="h-5 w-5 shrink-0" />
              <span className="hidden lg:block">{item.label}</span>
              {item.to === "/tickets" && data ? (
                <span className="ml-auto hidden rounded-full bg-amber-500/20 px-2 py-0.5 text-xs text-amber-200 lg:block">
                  {data.pendingTickets.today}
                </span>
              ) : null}
            </NavLink>
          );
        })}
      </nav>
      <div className="mt-auto hidden rounded-lg border border-slate-800 bg-surface-2 p-3 lg:block">
        <div className="flex items-center gap-2 text-sm">
          <span className="h-2.5 w-2.5 rounded-full bg-emerald-400" />
          <span className="text-slate-300">System Status</span>
        </div>
      </div>
    </aside>
  );
}
