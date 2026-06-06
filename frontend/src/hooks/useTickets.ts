import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createTicket, getTicket, getTicketEvents, getTickets, reprocessTicket, resolveTicket } from "@/api/tickets";
import type { Ticket, TicketFilters } from "@/types";

export function useTickets(filters: TicketFilters = {}) {
  return useQuery({ queryKey: ["tickets", filters], queryFn: () => getTickets(filters), staleTime: 15_000 });
}

export function useTicket(id: string) {
  return useQuery({ queryKey: ["ticket", id], queryFn: () => getTicket(id), enabled: Boolean(id) });
}

export function useTicketEvents(id: string) {
  const client = useQueryClient()
  const ticketData = client.getQueryData<Ticket>(['ticket', id])
  return { data: ticketData?.events || [] }
}

export function useTicketMutations() {
  const queryClient = useQueryClient();
  return {
    create: useMutation({ mutationFn: createTicket }),
    reprocess: useMutation({
      mutationFn: reprocessTicket,
      onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tickets"] }),
    }),
    resolve: useMutation({
      mutationFn: resolveTicket,
      onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tickets"] }),
    }),
  };
}
