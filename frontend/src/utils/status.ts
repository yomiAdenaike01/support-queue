import type { TicketStatus, WorkerStatus } from "@/types";

export const STATUS_CONFIG: Record<TicketStatus, { label: string; colour: string; icon: string }> = {
  PENDING: { label: "Pending", colour: "amber", icon: "Clock" },
  PROCESSING: { label: "Processing", colour: "blue", icon: "Loader" },
  PROCESSED: { label: "Processed", colour: "green", icon: "CheckCircle" },
  FAILED: { label: "Failed", colour: "red", icon: "XCircle" },
  RESOLVED: { label: "Resolved", colour: "purple", icon: "CheckCheck" },
};

export const statusClasses: Record<TicketStatus | WorkerStatus, string> = {
  PENDING: "bg-amber-500/15 text-amber-300 border-amber-500/30",
  PROCESSING: "bg-blue-500/15 text-blue-300 border-blue-500/30",
  PROCESSED: "bg-emerald-500/15 text-emerald-300 border-emerald-500/30",
  FAILED: "bg-red-500/15 text-red-300 border-red-500/30",
  RESOLVED: "bg-purple-500/15 text-purple-300 border-purple-500/30",
  ONLINE: "bg-emerald-500/15 text-emerald-300 border-emerald-500/30",
  OFFLINE: "bg-red-500/15 text-red-300 border-red-500/30",
  DEGRADED: "bg-amber-500/15 text-amber-300 border-amber-500/30",
};
