import { EmptyState } from "@/components/ui/EmptyState";
import { SkeletonRows } from "@/components/ui/Skeleton";
import { Table } from "@/components/ui/Table";
import { TicketRow } from "@/components/tickets/TicketRow";
import type { Ticket } from "@/types";

export function TicketTable({ tickets, loading }: { tickets: Ticket[]; loading: boolean }) {
  if (loading) return <SkeletonRows rows={8} />;
  if (tickets.length === 0) return <EmptyState title="No tickets found" message="Adjust filters or create a ticket to populate the queue." />;
  return (
    <Table>
      <thead className="bg-surface-2 text-left text-xs uppercase tracking-wide text-slate-400">
        <tr>
          {["#", "Customer", "Subject", "Status", "Priority", "Category", "Team", "Sentiment", "Created", "Actions"].map((header) => (
            <th key={header} className="px-4 py-3 font-medium">{header}</th>
          ))}
        </tr>
      </thead>
      <tbody className="divide-y divide-slate-800">
        {tickets.map((ticket, index) => <TicketRow key={ticket.id} ticket={ticket} index={index} />)}
      </tbody>
    </Table>
  );
}
