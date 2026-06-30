import { ArrowLeft, Copy, RotateCcw } from "lucide-react";
import { useEffect, useState } from "react";
import type { ReactNode } from "react";
import { Link, useParams } from "react-router-dom";
import { Badge, CategoryBadge, PriorityBadge, SentimentBadge, StatusBadge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Modal } from "@/components/ui/Modal";
import { SkeletonRows } from "@/components/ui/Skeleton";
import { useToast } from "@/components/ui/Toast";
import { MessageThread } from "@/components/tickets/MessageThread";
import { MessageComposer } from "@/components/tickets/MessageComposer";
import { TicketTimeline } from "@/components/tickets/TicketTimeline";
import { useTeams } from "@/hooks/useTeams";
import { TICKET_DETAIL_REFRESH_INTERVAL_MS, useTicket, useTicketEvents, useTicketMutations } from "@/hooks/useTickets";
import { mockTeams } from "@/mocks/teams";
import type { TicketCategory, TicketPriority } from "@/types";
import { formatDate } from "@/utils/format";

const priorityOptions: TicketPriority[] = ["URGENT", "HIGH", "MEDIUM", "LOW"];
const categoryOptions: TicketCategory[] = ["BILLING", "ACCOUNT_ACCESS", "TECHNICAL", "DELIVERY", "CANCELLATION", "SUBSCRIPTIONS", "GENERAL"];

export function TicketDetail() {
  const { id = "" } = useParams();
  const ticket = useTicket(id);
  const events = useTicketEvents(id);
  const teams = useTeams();
  const mutations = useTicketMutations();
  const [confirm, setConfirm] = useState<"reprocess" | "resolve" | "rerunPipeline" | null>(null);
  const [priority, setPriority] = useState<TicketPriority | "">("");
  const [category, setCategory] = useState<TicketCategory | "">("");
  const [assignedTeam, setAssignedTeam] = useState("");
  const [editingClassificationField, setEditingClassificationField] = useState<"category" | "priority" | "assignedTeam" | null>(null);
  const [notifyTeam, setNotifyTeam] = useState(false);
  const [now, setNow] = useState(Date.now());
  const toast = useToast();

  useEffect(() => {
    const interval = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(interval);
  }, []);

  useEffect(() => {
    if (!ticket.data) return;
    setPriority(ticket.data.priority ?? "");
    setCategory(ticket.data.category ?? "");
    setAssignedTeam(ticket.data.assignedTeam ?? "");
    setNotifyTeam(false);
  }, [ticket.data?.id, ticket.data?.priority, ticket.data?.category, ticket.data?.assignedTeam]);

  if (ticket.isLoading) return <SkeletonRows rows={8} />;
  if (!ticket.data) return <div className="text-red-200">Ticket not found.</div>;
  const item = ticket.data;
  const teamOptions = uniqueTeamDepartments([
    ...(teams.data ?? []).map((team) => team.department),
    ...mockTeams.map((team) => team.department),
    assignedTeam,
  ]);

  const runAction = async () => {
    if (confirm === "reprocess") await mutations.reprocess.mutateAsync(item.id);
    if (confirm === "resolve") await mutations.resolve.mutateAsync(item.id);
    if (confirm === "rerunPipeline") await mutations.rerunPipeline.mutateAsync(item.id);
    toast.push(confirm === "rerunPipeline" ? "Classifier queued" : "Ticket updated", "success");
    setConfirm(null);
  };

  const saveClassification = async () => {
    await mutations.updateClassification.mutateAsync({
      ticketId: item.id,
      priority: priority || null,
      category: category || null,
      assignedTeam: assignedTeam || null,
      notifyTeam,
    });
    toast.push(notifyTeam ? "Classification saved and team notified" : "Classification saved", "success");
    setNotifyTeam(false);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <Link to="/tickets" className="inline-flex items-center gap-2 text-sm text-slate-300 hover:text-white"><ArrowLeft className="h-4 w-4" /> Back</Link>
        <h1 className="font-mono text-xl font-semibold">{item.id}</h1>
        <StatusBadge status={item.status} />
        <PriorityBadge priority={item.priority} />
        <RefreshIndicator dataUpdatedAt={ticket.dataUpdatedAt} isFetching={ticket.isFetching} now={now} />
        <div className="ml-auto flex gap-2">
          <Button
            variant="secondary"
            onClick={() => setConfirm("rerunPipeline")}
            title="Recalculate category, priority, urgency, sentiment, routing, and suggested response."
          >
            <RotateCcw className="h-4 w-4" /> Rerun Classifier
          </Button>
          {item.status === "FAILED" ? (
            <Button
              variant="secondary"
              onClick={() => setConfirm("reprocess")}
              title="Retry a failed ticket through automated processing."
            >
              Retry Processing
            </Button>
          ) : null}
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
          <div className="space-y-4 text-sm">
            <EditablePillSelect
              label="Category"
              field="category"
              value={category}
              editingField={editingClassificationField}
              setEditingField={setEditingClassificationField}
              onChange={(value) => {
                setCategory(value as TicketCategory | "");
                setEditingClassificationField(null);
              }}
              options={categoryOptions}
              placeholder="Unclassified"
              renderPill={(value) => value ? <CategoryBadge category={value as TicketCategory} /> : <EmptyPill>Unclassified</EmptyPill>}
            />
            <EditablePillSelect
              label="Priority"
              field="priority"
              value={priority}
              editingField={editingClassificationField}
              setEditingField={setEditingClassificationField}
              onChange={(value) => {
                setPriority(value as TicketPriority | "");
                setEditingClassificationField(null);
              }}
              options={priorityOptions}
              placeholder="Unscored"
              renderPill={(value) => <PriorityBadge priority={(value || null) as TicketPriority | null} />}
            />
            <EditablePillSelect
              label="Assigned team"
              field="assignedTeam"
              value={assignedTeam}
              editingField={editingClassificationField}
              setEditingField={setEditingClassificationField}
              onChange={(value) => {
                setAssignedTeam(value);
                setEditingClassificationField(null);
              }}
              options={teamOptions}
              placeholder="Unassigned"
              renderPill={(value) => value ? <Badge className="border-blue-500/30 bg-blue-500/15 text-blue-300">{value}</Badge> : <EmptyPill>Unassigned</EmptyPill>}
            />
            <label className="flex items-center justify-between gap-4 rounded-lg border border-slate-800 bg-surface-2 p-3">
              <span>
                <span className="block font-medium text-slate-200">Notify assigned team</span>
                <span className="mt-1 block text-xs text-slate-500">Send a notification when these changes are saved.</span>
              </span>
              <input
                type="checkbox"
                checked={notifyTeam}
                disabled={!assignedTeam}
                onChange={(event) => setNotifyTeam(event.target.checked)}
                className="h-4 w-4 rounded border-slate-600 bg-surface"
              />
            </label>
            <Button
              className="w-full"
              loading={mutations.updateClassification.isPending}
              onClick={() => void saveClassification()}
            >
              Save Classification
            </Button>
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
      <Modal title={modalTitle(confirm)} open={Boolean(confirm)} onClose={() => setConfirm(null)}>
        <p className="mb-5 text-sm text-slate-300">{modalDescription(confirm, item.id)}</p>
        <div className="flex justify-end gap-2">
          <Button variant="secondary" onClick={() => setConfirm(null)}>Cancel</Button>
          <Button onClick={() => void runAction()} loading={mutations.reprocess.isPending || mutations.resolve.isPending || mutations.rerunPipeline.isPending}>Confirm</Button>
        </div>
      </Modal>
    </div>
  );
}

function Row({ label, value }: { label: string; value: ReactNode }) {
  return <div className="flex items-center justify-between gap-4"><span className="text-slate-400">{label}</span><span>{value}</span></div>;
}

function EditablePillSelect({
  label,
  field,
  value,
  editingField,
  setEditingField,
  onChange,
  options,
  placeholder,
  renderPill,
}: {
  label: string;
  field: "category" | "priority" | "assignedTeam";
  value: string;
  editingField: "category" | "priority" | "assignedTeam" | null;
  setEditingField: (field: "category" | "priority" | "assignedTeam" | null) => void;
  onChange: (value: string) => void;
  options: string[];
  placeholder: string;
  renderPill: (value: string) => ReactNode;
}) {
  const isEditing = editingField === field;
  return (
    <div className="flex min-h-10 items-center justify-between gap-4">
      <span className="text-slate-400">{label}</span>
      {isEditing ? (
        <select
          autoFocus
          value={value}
          onBlur={() => setEditingField(null)}
          onChange={(event) => onChange(event.target.value)}
          className="min-w-44 rounded-full border border-slate-700 bg-surface-2 px-3 py-1.5 text-xs font-semibold outline-none focus:border-accent"
        >
          <option value="">{placeholder}</option>
          {options.map((option) => (
            <option key={option} value={option}>
              {option.replace(/_/g, " ")}
            </option>
          ))}
        </select>
      ) : (
        <button
          type="button"
          onClick={() => setEditingField(field)}
          className="rounded-full outline-none transition hover:scale-[1.02] focus:ring-2 focus:ring-accent/60"
          title={`Change ${label.toLowerCase()}`}
        >
          {renderPill(value)}
        </button>
      )}
    </div>
  );
}

function EmptyPill({ children }: { children: ReactNode }) {
  return <Badge className="border-slate-600 bg-slate-800 text-slate-400">{children}</Badge>;
}

function uniqueTeamDepartments(departments: string[]): string[] {
  return Array.from(new Set(departments.map((department) => department.trim()).filter(Boolean)));
}

function RefreshIndicator({ dataUpdatedAt, isFetching, now }: { dataUpdatedAt: number; isFetching: boolean; now: number }) {
  const nextRefreshAt = dataUpdatedAt + TICKET_DETAIL_REFRESH_INTERVAL_MS;
  const remainingMs = Math.max(0, nextRefreshAt - now);
  const lastRefresh = dataUpdatedAt > 0 ? new Date(dataUpdatedAt).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }) : "pending";

  return (
    <div className="inline-flex items-center gap-2 rounded-full border border-blue-500/30 bg-blue-500/10 px-3 py-1 text-xs font-medium text-blue-200">
      <span className={`h-2 w-2 rounded-full bg-blue-300 ${isFetching ? "animate-pulse" : ""}`} />
      <span>{isFetching ? "Refreshing ticket" : `Next refresh in ${formatCountdown(remainingMs)}`}</span>
      <span className="hidden text-blue-200/70 sm:inline">Last {lastRefresh}</span>
    </div>
  );
}

function formatCountdown(ms: number): string {
  const totalSeconds = Math.ceil(ms / 1000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}m ${seconds.toString().padStart(2, "0")}s`;
}

function modalTitle(action: "reprocess" | "resolve" | "rerunPipeline" | null): string {
  if (action === "reprocess") return "Retry processing";
  if (action === "rerunPipeline") return "Rerun classifier";
  return "Resolve ticket";
}

function modalDescription(action: "reprocess" | "resolve" | "rerunPipeline" | null, ticketId: string): string {
  if (action === "rerunPipeline") return `Queue ${ticketId} to recalculate classification, priority, urgency, routing, and suggested response.`;
  if (action === "reprocess") return `Retry automated processing for failed ticket ${ticketId}. This queues the ticket back onto the worker stream.`;
  return `Confirm this ticket action for ${ticketId}.`;
}
