import { Button } from "@/components/ui/Button";
import { useUiStore } from "@/store/ui";
import type { TicketCategory, TicketPriority, TicketStatus } from "@/types";

const statuses: (TicketStatus | "ALL")[] = ["ALL", "PENDING", "PROCESSING", "PROCESSED", "FAILED", "RESOLVED"];
const priorities: (TicketPriority | "ALL")[] = ["ALL", "URGENT", "HIGH", "MEDIUM", "LOW"];
const categories: (TicketCategory | "ALL")[] = ["ALL", "BILLING", "TECHNICAL", "DELIVERY", "SUBSCRIPTIONS", "GENERAL", "ACCOUNT_ACCESS", "CANCELLATION"];

export function TicketFilters() {
  const { search, status, priority, category, setSearch, setStatus, setPriority, setCategory, clearFilters } = useUiStore();
  return (
    <div className="flex flex-wrap gap-3 rounded-xl border border-slate-800 bg-surface p-4">
      <input
        value={search}
        onChange={(event) => setSearch(event.target.value)}
        placeholder="Search subject or email"
        className="min-w-64 flex-1 rounded-lg border border-slate-700 bg-surface-2 px-3 py-2 text-sm outline-none focus:border-accent"
      />
      <Select value={status} onChange={(value) => setStatus(value as TicketStatus | "ALL")} options={statuses} />
      <Select value={priority} onChange={(value) => setPriority(value as TicketPriority | "ALL")} options={priorities} />
      <Select value={category} onChange={(value) => setCategory(value as TicketCategory | "ALL")} options={categories} />
      <Button variant="secondary" onClick={clearFilters}>Clear</Button>
    </div>
  );
}

function Select({ value, onChange, options }: { value: string; onChange: (value: string) => void; options: string[] }) {
  return (
    <select
      value={value}
      onChange={(event) => onChange(event.target.value)}
      className="rounded-lg border border-slate-700 bg-surface-2 px-3 py-2 text-sm outline-none focus:border-accent"
    >
      {options.map((item) => <option key={item}>{item}</option>)}
    </select>
  );
}
