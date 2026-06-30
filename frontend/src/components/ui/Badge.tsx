import type { ReactNode } from "react";
import type {
  SentimentLabel,
  TicketCategory,
  TicketPriority,
  TicketStatus,
  WorkerStatus,
} from "@/types";
import { priorityClasses, PRIORITY_CONFIG } from "@/utils/priority";
import { statusClasses, STATUS_CONFIG } from "@/utils/status";

const categoryClasses: Record<TicketCategory, string> = {
  BILLING: "bg-cyan-500/15 text-cyan-300 border-cyan-500/30",
  ACCOUNT_ACCESS: "bg-indigo-500/15 text-indigo-300 border-indigo-500/30",
  TECHNICAL: "bg-blue-500/15 text-blue-300 border-blue-500/30",
  DELIVERY: "bg-emerald-500/15 text-emerald-300 border-emerald-500/30",
  CANCELLATION: "bg-red-500/15 text-red-300 border-red-500/30",
  SUBSCRIPTIONS: "bg-purple-500/15 text-purple-300 border-purple-500/30",
  GENERAL: "bg-slate-500/15 text-slate-300 border-slate-500/30",
};

const sentimentClasses: Record<SentimentLabel, string> = {
  HIGH: "bg-red-500/15 text-red-300 border-red-500/30",
  MEDIUM: "bg-amber-500/15 text-amber-300 border-amber-500/30",
  LOW: "bg-emerald-500/15 text-emerald-300 border-emerald-500/30",
};

interface BadgeProps {
  children: ReactNode;
  className?: string;
}

export function Badge({ children, className = "" }: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs font-semibold ${className}`}
    >
      {children}
    </span>
  );
}

export function StatusBadge({
  status,
}: {
  status: TicketStatus | WorkerStatus;
}) {
  console.log({ status });
  const label =
    status in STATUS_CONFIG
      ? STATUS_CONFIG[status as TicketStatus].label
      : status;
  return <Badge className={statusClasses[status]}>{label}</Badge>;
}

export function PriorityBadge({
  priority,
}: {
  priority: TicketPriority | null;
}) {
  if (!priority)
    return (
      <Badge className="border-slate-600 bg-slate-800 text-slate-400">
        Unscored
      </Badge>
    );
  return (
    <Badge className={priorityClasses[priority]}>
      {PRIORITY_CONFIG[priority].label}
    </Badge>
  );
}

export function CategoryBadge({ category }: { category: TicketCategory }) {
  return (
    <Badge className={categoryClasses[category]}>
      {category.replace("_", " ")}
    </Badge>
  );
}

export function SentimentBadge({
  sentiment,
}: {
  sentiment: SentimentLabel | null;
}) {
  if (!sentiment)
    return (
      <Badge className="border-slate-600 bg-slate-800 text-slate-400">
        Pending
      </Badge>
    );
  return <Badge className={sentimentClasses[sentiment]}>{sentiment}</Badge>;
}
