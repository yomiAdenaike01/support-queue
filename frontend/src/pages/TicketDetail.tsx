import { ArrowLeft, Copy } from "lucide-react";
import { useState } from "react";
import type { ReactNode } from "react";
import { Link, useParams } from "react-router-dom";
import { CategoryBadge, PriorityBadge, SentimentBadge, StatusBadge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Modal } from "@/components/ui/Modal";
import { SkeletonRows } from "@/components/ui/Skeleton";
import { useToast } from "@/components/ui/Toast";
import { MessageThread } from "@/components/tickets/MessageThread";
import { MessageComposer } from "@/components/tickets/MessageComposer";
import { TicketTimeline } from "@/components/tickets/TicketTimeline";
import { useTicket, useTicketEvents, useTicketMutations } from "@/hooks/useTickets";
import { formatDate } from "@/utils/format";

export function TicketDetail() {
  const { id = "" } = useParams();
  const ticket = useTicket(id);
  const events = useTicketEvents(id);
  const mutations = useTicketMutations();
  const [confirm, setConfirm] = useState<"reprocess" | "resolve" | null>(null);
  const toast = useToast();

  if (ticket.isLoading) return <SkeletonRows rows={8} />;
  if (!ticket.data) return <div className="text-red-200">Ticket not found.</div>;
  const item = ticket.data;

  const runAction = async () => {
    if (confirm === "reprocess") await mutations.reprocess.mutateAsync(item.id);
    if (confirm === "resolve") await mutations.resolve.mutateAsync(item.id);
    toast.push("Ticket updated", "success");
    setConfirm(null);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <Link to="/tickets" className="inline-flex items-center gap-2 text-sm text-slate-300 hover:text-white"><ArrowLeft className="h-4 w-4" /> Back</Link>
        <h1 className="font-mono text-xl font-semibold">{item.id}</h1>
        <StatusBadge status={item.status} />
        <PriorityBadge priority={item.priority} />
        <div className="ml-auto flex gap-2">
          {item.status === "FAILED" ? <Button variant="secondary" onClick={() => setConfirm("reprocess")}>Reprocess</Button> : null}
          {item.status !== "RESOLVED" ? <Button onClick={() => setConfirm("resolve")}>Mark Resolved</Button> : null}
        </div>
      </div>
      <div className="grid gap-6 xl:grid-cols-[3fr_2fr]">
        <Card>
          <h2 className="text-xl font-semibold">{item.subject}</h2>
          <div className="mt-2 text-sm text-slate-400">{item.customerEmail} · Created {formatDate(item.createdAt)}</div>
          <div className="mt-6"><MessageThread messages={item.messages} /></div>
          <MessageComposer
            loading={mutations.addMessage.isPending}
            onSubmit={async (message) => {
              await mutations.addMessage.mutateAsync({
                ticketId: item.id,
                customerEmail: item.customerEmail,
                content: message.content,
                role: message.role,
              });
              toast.push("Message added", "success");
            }}
          />
        </Card>
        <Card>
          <h2 className="mb-4 text-lg font-semibold">Classification</h2>
          <div className="space-y-3 text-sm">
            <Row label="Category" value={item.category ? <CategoryBadge category={item.category} /> : "Pending"} />
            <Row label="Priority" value={<PriorityBadge priority={item.priority} />} />
            <Row label="Assigned team" value={item.assignedTeam ?? "Unassigned"} />
            <Row label="Sentiment" value={<SentimentBadge sentiment={item.sentimentLabel} />} />
            <Row label="Score" value={item.sentimentScore?.toFixed(2) ?? "Pending"} />
          </div>
          {item.urgencyFlag ? <div className="mt-5 rounded-lg border border-red-500/30 bg-red-500/10 p-3 text-sm text-red-100">{item.urgencyReason}</div> : null}
          {item.suggestedResponse ? (
            <div className="mt-5 rounded-lg border border-slate-700 bg-surface-2 p-4">
              <div className="mb-2 flex items-center justify-between">
                <span className="text-sm font-semibold">Suggested response</span>
                <Button variant="ghost" className="h-8 px-2" onClick={() => { void navigator.clipboard.writeText(item.suggestedResponse ?? ""); toast.push("Copied response", "success"); }}>
                  <Copy className="h-4 w-4" />
                </Button>
              </div>
              <p className="text-sm leading-6 text-slate-300">{item.suggestedResponse}</p>
            </div>
          ) : null}
        </Card>
      </div>
      <TicketTimeline events={events.data ?? []} />
      <Modal title={confirm === "reprocess" ? "Reprocess ticket" : "Resolve ticket"} open={Boolean(confirm)} onClose={() => setConfirm(null)}>
        <p className="mb-5 text-sm text-slate-300">Confirm this ticket action for {item.id}.</p>
        <div className="flex justify-end gap-2">
          <Button variant="secondary" onClick={() => setConfirm(null)}>Cancel</Button>
          <Button onClick={() => void runAction()} loading={mutations.reprocess.isPending || mutations.resolve.isPending}>Confirm</Button>
        </div>
      </Modal>
    </div>
  );
}

function Row({ label, value }: { label: string; value: ReactNode }) {
  return <div className="flex items-center justify-between gap-4"><span className="text-slate-400">{label}</span><span>{value}</span></div>;
}
