import {
  AlertTriangle,
  CheckCircle,
  Clock,
  Inbox,
  Loader,
  MailWarning,
} from "lucide-react";
import { Card } from "@/components/ui/Card";
import { Skeleton } from "@/components/ui/Skeleton";
import type { Metrics } from "@/types";
import { formatNumber } from "@/utils/format";

export function MetricsRow({
  metrics,
  loading,
}: {
  metrics?: Metrics;
  loading: boolean;
}) {
  const items = metrics
    ? [
        {
          label: "Total Tickets",
          valueKey: "totalTickets",
          value: metrics.totalTickets.today,
          icon: Inbox,
          tone: "text-slate-200",
        },
        {
          label: "Pending",
          valueKey: "pendingTickets",
          value: metrics.pendingTickets.today,
          icon: Clock,
          tone: "text-amber-300",
        },
        {
          label: "Processing",
          valueKey: "processingTickets",
          value: metrics.processingTickets.today,
          icon: Loader,
          tone: "text-blue-300",
        },
        {
          label: "Processed",
          valueKey: "processedTickets",
          value: metrics.processedTickets.today,
          icon: CheckCircle,
          tone: "text-emerald-300",
        },
        {
          label: "Failed",
          valueKey: "failedTickets",
          value: metrics.failedTickets.today,
          icon: AlertTriangle,
          tone: "text-red-300",
        },
        {
          label: "Dead Letter",
          valueKey: "failedTickets",
          value: metrics?.deadLetterCount?.today || 0,
          icon: MailWarning,
          tone: "text-red-300",
        },
      ]
    : [];
  return (
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-6">
      {loading
        ? Array.from({ length: 6 }, (_, index) => (
            <Skeleton key={index} className="h-28" />
          ))
        : items.map((item) => {
            const Icon = item.icon;
            const percentageDifference =
              metrics[item.valueKey]?.percentageDifference ?? 0;
            return (
              <Card key={item.label} className="p-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-slate-400">{item.label}</span>
                  <Icon
                    className={`h-5 w-5 ${item.tone} ${item.label === "Processing" && item.value > 0 ? "animate-spin" : ""}`}
                  />
                </div>
                <div className="mt-4 font-mono text-3xl font-semibold">
                  {formatNumber(item.value)}
                </div>
                <div className="mt-1 text-xs text-emerald-300">
                  {percentageDifference > 0
                    ? `+${formatNumber(percentageDifference)}`
                    : formatNumber(percentageDifference)}
                  % vs yesterday
                </div>
              </Card>
            );
          })}
    </div>
  );
}
