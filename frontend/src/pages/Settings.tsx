import { Card } from "@/components/ui/Card";
import { StatusBadge } from "@/components/ui/Badge";
import { USE_MOCK, BASE_URL } from "@/api/client";
import { useWorkers } from "@/hooks/useWorkers";
import { useStream } from "@/hooks/useStream";

export function Settings() {
  const workers = useWorkers();
  const stream = useStream();
  const services = [
    ["Go API", "ONLINE", "v1.0.0", "response: 12ms"],
    ["Python Worker", "ONLINE", "running", `processed: ${workers.data?.[0]?.processedCount ?? 0}`],
    ["Redis", "ONLINE", "connected", `pending: ${stream.data?.pendingMessages ?? 0}`],
    ["Postgres", "ONLINE", "connected", "size: 24MB"],
    ["Ollama", "DEGRADED", "slow", "model: llama3.2"],
  ] as const;
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Settings</h1>
      <Card>
        <h2 className="mb-4 text-lg font-semibold">Service Health</h2>
        <div className="space-y-3">
          {services.map(([name, status, meta, detail]) => (
            <div key={name} className="grid gap-3 rounded-lg bg-surface-2 p-3 text-sm md:grid-cols-4">
              <span className="font-medium">{name}</span>
              <StatusBadge status={status} />
              <span className="text-slate-300">{meta}</span>
              <span className="text-slate-400">{detail}</span>
            </div>
          ))}
        </div>
      </Card>
      <div className="grid gap-6 xl:grid-cols-2">
        <Card>
          <h2 className="mb-4 text-lg font-semibold">Stream Configuration</h2>
          {[
            ["Stream key", "tickets"],
            ["Consumer group", "workers"],
            ["Dead letter stream", "tickets:dead-letter"],
            ["Max retries", "3"],
          ].map(([label, value]) => <Row key={label} label={label} value={value} />)}
        </Card>
        <Card>
          <h2 className="mb-4 text-lg font-semibold">API Configuration</h2>
          <Row label="API base URL" value={BASE_URL} />
          <Row label="Mock mode" value={USE_MOCK ? "Enabled" : "Disabled"} />
          <Row label="Refresh interval" value={`${Number(import.meta.env.VITE_REFRESH_INTERVAL ?? 30000) / 1000}s`} />
        </Card>
      </div>
      <Card>
        <h2 className="mb-4 text-lg font-semibold">Worker Configuration</h2>
        <div className="grid gap-3 md:grid-cols-2">
          {(workers.data ?? []).map((worker) => (
            <div key={worker.name} className="rounded-lg bg-surface-2 p-3 text-sm">
              <div className="font-medium">{worker.name}</div>
              <div className="mt-1 text-slate-400">{worker.language} · {worker.processedCount} processed</div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return <div className="flex items-center justify-between border-b border-slate-800 py-3 text-sm last:border-b-0"><span className="text-slate-400">{label}</span><span className="font-mono">{value}</span></div>;
}
