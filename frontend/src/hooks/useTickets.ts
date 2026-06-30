import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { addTicketMessage, createTicket, getTicket, getTicketEvents, getTickets, reprocessTicket, rerunTicketPipeline, resolveTicket, updateTicketClassification } from "@/api/tickets";
import type { Ticket, TicketFilters } from "@/types";

export const TICKET_DETAIL_REFRESH_INTERVAL_MS = 120_000;

export function useTickets(filters: TicketFilters = {}) {
  return useQuery({ queryKey: ["tickets", filters], queryFn: () => getTickets(filters), staleTime: 15_000 });
}

export function useTicket(id: string) {
  return useQuery({
    queryKey: ["ticket", id],
    queryFn: () => getTicket(id),
    enabled: Boolean(id),
    refetchInterval: TICKET_DETAIL_REFRESH_INTERVAL_MS,
    refetchIntervalInBackground: true,
  });
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
      onSuccess: (ticket, id) => {
        queryClient.setQueryData(["ticket", id], ticket);
        void queryClient.invalidateQueries({ queryKey: ["tickets"] });
      },
    }),
    resolve: useMutation({
      mutationFn: resolveTicket,
      onSuccess: (ticket, id) => {
        queryClient.setQueryData(["ticket", id], ticket);
        void queryClient.invalidateQueries({ queryKey: ["tickets"] });
      },
    }),
    rerunPipeline: useMutation({
      mutationFn: rerunTicketPipeline,
      onSuccess: (_, id) => {
        void queryClient.invalidateQueries({ queryKey: ["ticket", id] });
        void queryClient.invalidateQueries({ queryKey: ["tickets"] });
      },
    }),
    updateClassification: useMutation({
      mutationFn: updateTicketClassification,
      onSuccess: (ticket, input) => {
        queryClient.setQueryData(["ticket", input.ticketId], ticket);
        void queryClient.invalidateQueries({ queryKey: ["tickets"] });
      },
    }),
    addMessage: useMutation({
      mutationFn: addTicketMessage,
      onSuccess: (_, input) => {
        void queryClient.invalidateQueries({ queryKey: ["ticket", input.ticketId] });
        void queryClient.invalidateQueries({ queryKey: ["tickets"] });
      },
    }),
  };
}
