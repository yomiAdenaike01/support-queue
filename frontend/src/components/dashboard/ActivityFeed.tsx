import { Card } from "@/components/ui/Card";
import { mockEvents } from "@/mocks/events";
import { relativeTime } from "@/utils/format";

const tone: Record<string, string> = {
  TICKET_CREATED: "text-blue-300",
  TICKET_PROCESSING: "text-amber-300",
  TICKET_CLASSIFIED: "text-purple-300",
  TICKET_PROCESSED: "text-emerald-300",
  TICKET_FAILED: "text-red-300",
  TICKET_REPROCESSED: "text-amber-300",
  TICKET_RESOLVED: "text-emerald-300",
};

export function ActivityFeed() {
  return (
    <Card>
      <h2 className="mb-4 text-lg font-semibold">Recent Activity</h2>
      <div className="space-y-2">
        {mockEvents
          .slice()
          .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
          .slice(0, 20)
          .map((event) => (
            <div key={event.id} className="flex items-center justify-between rounded-lg bg-surface-2 px-3 py-2 text-sm">
              <span className={tone[event.eventType]}>{event.eventType.replace(/_/g, " ")}</span>
              <span className="text-slate-400">{event.ticketId} · {relativeTime(event.createdAt)}</span>
            </div>
          ))}
      </div>
    </Card>
  );
}
