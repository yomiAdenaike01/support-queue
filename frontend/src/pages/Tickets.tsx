import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { TicketFilters } from "@/components/tickets/TicketFilters";
import { TicketTable } from "@/components/tickets/TicketTable";
import { useTickets } from "@/hooks/useTickets";
import { useUiStore } from "@/store/ui";

export function Tickets() {
  const { search, status, priority, category } = useUiStore();
  const query = useTickets({ search, status, priority, category, limit: 20 });
  const pageTickets = (query.data ?? []).slice(0, 20);
  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Tickets</h1>
          <p className="text-sm text-slate-400">Filter, review, and route support work.</p>
        </div>
        <Button variant="secondary" onClick={() => void query.refetch()} loading={query.isFetching}>
          <RefreshCw className="h-4 w-4" /> Refresh
        </Button>
      </div>
      <TicketFilters />
      {query.isError ? <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-3 text-red-200">Unable to load tickets.</div> : null}
      <TicketTable tickets={pageTickets} loading={query.isLoading} />
      <div className="text-sm text-slate-400">Showing {pageTickets.length} of {query.data?.length ?? 0} tickets · 20 per page</div>
    </div>
  );
}
