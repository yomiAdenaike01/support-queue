import { AlertTriangle, CheckCircle, Clock, Inbox, Loader, MailWarning } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { Skeleton } from "@/components/ui/Skeleton";
import type { Metrics } from "@/types";
import { formatNumber } from "@/utils/format";

export function MetricsRow({ metrics, loading }: { metrics?: Metrics; loading: boolean }) {
  const items = metrics
    ? [
        { label: "Total Tickets", value: metrics.totalTickets, icon: Inbox, tone: "text-slate-200" },
        { label: "Pending", value: metrics.pendingTickets, icon: Clock, tone: "text-amber-300" },
        { label: "Processing", value: metrics.processingTickets, icon: Loader, tone: "text-blue-300" },
        { label: "Processed", value: metrics.processedTickets, icon: CheckCircle, tone: "text-emerald-300" },
        { label: "Failed", value: metrics.failedTickets, icon: AlertTriangle, tone: "text-red-300" },
        { label: "Dead Letter", value: metrics.deadLetterCount, icon: MailWarning, tone: "text-red-300" },
      ]
    : [];
  return (
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-6">
      {loading
        ? Array.from({ length: 6 }, (_, index) => <Skeleton key={index} className="h-28" />)
        : items.map((item) => {
            const Icon = item.icon;
            return (
              <Card key={item.label} className="p-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-slate-400">{item.label}</span>
                  <Icon className={`h-5 w-5 ${item.tone} ${item.label === "Processing" && item.value > 0 ? "animate-spin" : ""}`} />
                </div>
                <div className="mt-4 font-mono text-3xl font-semibold">{formatNumber(item.value)}</div>
                <div className="mt-1 text-xs text-emerald-300">+8.4% vs yesterday</div>
              </Card>
            );
          })}
    </div>
  );
}
