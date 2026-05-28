import { Card } from "@/components/ui/Card";
import { StatusBadge } from "@/components/ui/Badge";
import type { Worker } from "@/types";
import { relativeTime } from "@/utils/format";

export function WorkerHealth({ workers }: { workers: Worker[] }) {
  return (
    <Card>
      <h2 className="mb-4 text-lg font-semibold">Worker Health</h2>
      <div className="space-y-3">
        {workers.map((worker) => {
          const heartbeatAge = Date.now() - new Date(worker.lastHeartbeat).getTime();
          const status = heartbeatAge > 60_000 ? "OFFLINE" : worker.status;
          return (
            <div key={worker.name} className="flex items-center justify-between gap-4 rounded-lg bg-surface-2 p-3">
              <div>
                <div className="font-medium">{worker.name}</div>
                <div className="text-xs text-slate-400">{worker.language} · {relativeTime(worker.lastHeartbeat)}</div>
              </div>
              <div className="text-right text-xs text-slate-400">
                <StatusBadge status={status} />
                <div className="mt-1">{worker.processedCount} ok · {worker.failedCount} failed</div>
              </div>
            </div>
          );
        })}
      </div>
    </Card>
  );
}
