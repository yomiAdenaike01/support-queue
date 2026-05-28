import { Card } from "@/components/ui/Card";
import type { Metrics, StreamStats as StreamStatsType } from "@/types";
import { formatDuration, formatNumber } from "@/utils/format";

export function StreamStats({ stats, metrics }: { stats: StreamStatsType; metrics?: Metrics }) {
  return (
    <Card>
      <h2 className="mb-4 text-lg font-semibold">Stream Stats</h2>
      <div className="grid grid-cols-2 gap-3">
        <Stat label="Pending" value={stats.pendingMessages} />
        <Stat label="Dead Letter" value={stats.deadLetterCount} />
        <Stat label="Processed" value={stats.totalProcessed} />
        <Stat label="Avg Time" value={formatDuration(metrics?.averageProcessingTimeMs ?? 0)} />
      </div>
    </Card>
  );
}

function Stat({ label, value }: { label: string; value: number | string }) {
  return (
    <div className="rounded-lg bg-surface-2 p-3">
      <div className="text-xs text-slate-400">{label}</div>
      <div className="mt-1 font-mono text-xl font-semibold">{typeof value === "number" ? formatNumber(value) : value}</div>
    </div>
  );
}
