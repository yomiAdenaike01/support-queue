import { ActivityFeed } from "@/components/dashboard/ActivityFeed";
import { MetricsRow } from "@/components/dashboard/MetricsRow";
import { SentimentChart } from "@/components/dashboard/SentimentChart";
import { StreamStats } from "@/components/dashboard/StreamStats";
import { TicketsByCategory } from "@/components/dashboard/TicketsByCategory";
import { TicketsByPriority } from "@/components/dashboard/TicketsByPriority";
import { TicketsOverTime } from "@/components/dashboard/TicketsOverTime";
import { WorkerHealth } from "@/components/dashboard/WorkerHealth";
import { SkeletonRows } from "@/components/ui/Skeleton";
import { useMetrics } from "@/hooks/useMetrics";
import { useStream } from "@/hooks/useStream";
import { useTickets } from "@/hooks/useTickets";
import { useWorkers } from "@/hooks/useWorkers";

export function Dashboard() {
  const metrics = useMetrics();
  const workers = useWorkers();
  const stream = useStream();
  const tickets = useTickets();

  return (
    <div className="space-y-6">
      <MetricsRow metrics={metrics.data} loading={metrics.isLoading} />
      {metrics.isError ? <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-3 text-red-200">Unable to load metrics.</div> : null}
      <div className="grid gap-6 xl:grid-cols-2">
        {workers.data ? <WorkerHealth workers={workers.data} /> : <SkeletonRows rows={3} />}
        {stream.data ? <StreamStats stats={stream.data} metrics={metrics.data} /> : <SkeletonRows rows={3} />}
      </div>
      <TicketsOverTime />
      <div className="grid gap-6 xl:grid-cols-3">
        <TicketsByCategory tickets={tickets.data ?? []} />
        <TicketsByPriority tickets={tickets.data ?? []} />
        <SentimentChart tickets={tickets.data ?? []} />
      </div>
      <ActivityFeed />
    </div>
  );
}
