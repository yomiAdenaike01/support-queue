import { useEffect, useState } from "react";
import { ChevronLeft, ChevronRight, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { TicketFilters } from "@/components/tickets/TicketFilters";
import { TicketTable } from "@/components/tickets/TicketTable";
import { useTickets } from "@/hooks/useTickets";
import { useUiStore } from "@/store/ui";

const PAGE_LIMIT = 20;
const INIT_PAGE = 0;

export function Tickets() {
  const { search, status, priority, category } = useUiStore();
  const [currentPage, setCurrentPage] = useState<number>(INIT_PAGE);
  useEffect(() => {
    setCurrentPage((s) => (s === INIT_PAGE ? s : INIT_PAGE));
  }, [search, status, priority, category]);

  const query = useTickets({
    search,
    status,
    priority,
    category,
    limit: PAGE_LIMIT,
  });
  const pageTickets = (query.data ?? []).slice(
    currentPage * PAGE_LIMIT,
    (currentPage + 1) * PAGE_LIMIT,
  );
  const totalTickets = query.data?.length ?? 0;
  const hasPreviousPage = currentPage > INIT_PAGE;
  const hasNextPage = (currentPage + 1) * PAGE_LIMIT < totalTickets;

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Tickets</h1>
          <p className="text-sm text-slate-400">
            Filter, review, and route support work.
          </p>
        </div>
        <Button
          variant="secondary"
          onClick={() => void query.refetch()}
          loading={query.isFetching}
        >
          <RefreshCw className="h-4 w-4" /> Refresh
        </Button>
      </div>
      <TicketFilters />
      {query.isError ? (
        <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-3 text-red-200">
          Unable to load tickets.
        </div>
      ) : null}
      <TicketTable tickets={pageTickets} loading={query.isLoading} />
      <div className="flex flex-col gap-3 text-sm text-slate-400 sm:flex-row sm:items-center sm:justify-between">
        <div>
          Showing {pageTickets.length} of {totalTickets} tickets · {PAGE_LIMIT}
          per page
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="secondary"
            onClick={() =>
              setCurrentPage((page) => Math.max(INIT_PAGE, page - 1))
            }
            disabled={!hasPreviousPage || query.isLoading}
          >
            <ChevronLeft className="h-4 w-4" /> Previous
          </Button>
          <span className="min-w-16 text-center">Page {currentPage + 1}</span>
          <Button
            variant="secondary"
            onClick={() => setCurrentPage((page) => page + 1)}
            disabled={!hasNextPage || query.isLoading}
          >
            Next <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
