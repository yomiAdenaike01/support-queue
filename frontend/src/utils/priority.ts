import type { TicketPriority } from "@/types";

export const PRIORITY_CONFIG: Record<TicketPriority, { label: string; colour: string; emoji: string }> = {
  URGENT: { label: "Urgent", colour: "red", emoji: "R" },
  HIGH: { label: "High", colour: "amber", emoji: "H" },
  MEDIUM: { label: "Medium", colour: "blue", emoji: "M" },
  LOW: { label: "Low", colour: "grey", emoji: "L" },
};

export const priorityClasses: Record<TicketPriority, string> = {
  URGENT: "bg-red-500/15 text-red-300 border-red-500/30",
  HIGH: "bg-amber-500/15 text-amber-300 border-amber-500/30",
  MEDIUM: "bg-blue-500/15 text-blue-300 border-blue-500/30",
  LOW: "bg-slate-500/15 text-slate-300 border-slate-500/30",
};
