import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import type { TicketEvent } from "@/types";
import { formatDate } from "@/utils/format";

const classes: Record<TicketEvent["eventType"], string> = {
  TICKET_CREATED: "bg-blue-500/15 text-blue-300 border-blue-500/30",
  TICKET_PROCESSING: "bg-amber-500/15 text-amber-300 border-amber-500/30",
  TICKET_CLASSIFIED: "bg-purple-500/15 text-purple-300 border-purple-500/30",
  TICKET_PROCESSED: "bg-emerald-500/15 text-emerald-300 border-emerald-500/30",
  TICKET_FAILED: "bg-red-500/15 text-red-300 border-red-500/30",
  TICKET_REPROCESSED: "bg-amber-500/15 text-amber-300 border-amber-500/30",
  TICKET_RESOLVED: "bg-emerald-500/15 text-emerald-300 border-emerald-500/30",
};

export function TicketTimeline({ events }: { events: TicketEvent[] }) {
  console.log({ events });
  return (
    <Card>
      <h2 className="mb-5 text-lg font-semibold">Event Timeline</h2>
      <div className="space-y-4">
        {events
          .slice()
          .sort(
            (a, b) =>
              new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
          )
          .map((event) => (
            <details
              key={event.eventType}
              className="border-l border-slate-700 pl-5"
            >
              <summary className="cursor-pointer list-none">
                <div className="flex flex-wrap items-center gap-3">
                  <Badge className={classes[event.eventType]}>
                    {String(event.eventType).replace(/_/g, " ")}
                  </Badge>
                  <span className="text-sm text-slate-400">
                    {formatDate(event.createdAt)}
                  </span>
                </div>
              </summary>
              <pre className="mt-3 overflow-x-auto rounded-lg bg-surface-2 p-3 text-xs text-slate-300">
                {JSON.stringify(event.payload ?? {}, null, 2)}
              </pre>
            </details>
          ))}
      </div>
    </Card>
  );
}
