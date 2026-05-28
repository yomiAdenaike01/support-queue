import { Eye, TriangleAlert } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { CategoryBadge, PriorityBadge, SentimentBadge, StatusBadge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import type { Ticket } from "@/types";
import { formatDate, relativeTime } from "@/utils/format";

export function TicketRow({ ticket, index }: { ticket: Ticket; index: number }) {
  const navigate = useNavigate();
  return (
    <tr onClick={() => navigate(`/tickets/${ticket.id}`)} className="cursor-pointer odd:bg-surface even:bg-surface-2/40 hover:bg-slate-700/40">
      <td className="px-4 py-3 font-mono text-slate-400">{index + 1}</td>
      <td className="px-4 py-3">{ticket.customerEmail}</td>
      <td className="px-4 py-3">
        <div className="flex items-center gap-2">
          {ticket.urgencyFlag ? <TriangleAlert className="h-4 w-4 text-red-300" /> : null}
          <span className="max-w-sm truncate">{ticket.subject}</span>
        </div>
      </td>
      <td className="px-4 py-3"><StatusBadge status={ticket.status} /></td>
      <td className="px-4 py-3"><PriorityBadge priority={ticket.priority} /></td>
      <td className="px-4 py-3">{ticket.category ? <CategoryBadge category={ticket.category} /> : <span className="text-slate-500">Pending</span>}</td>
      <td className="px-4 py-3">{ticket.assignedTeam ?? "Unassigned"}</td>
      <td className="px-4 py-3"><SentimentBadge sentiment={ticket.sentimentLabel} /></td>
      <td className="px-4 py-3" title={formatDate(ticket.createdAt)}>{relativeTime(ticket.createdAt)}</td>
      <td className="px-4 py-3">
        <Button variant="secondary" className="h-8 px-3" onClick={(event) => { event.stopPropagation(); navigate(`/tickets/${ticket.id}`); }}>
          <Eye className="h-4 w-4" /> View
        </Button>
      </td>
    </tr>
  );
}
