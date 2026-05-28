import { mockTickets } from "@/mocks/tickets";
import type { TicketEvent } from "@/types";

export const mockEvents: TicketEvent[] = mockTickets.flatMap((ticket, index) => {
  const base = new Date(ticket.createdAt).getTime();
  const events: TicketEvent[] = [
    {
      id: `evt-${ticket.id}-created`,
      ticketId: ticket.id,
      eventType: "TICKET_CREATED",
      payload: { source: "simulator" },
      createdAt: ticket.createdAt,
    },
  ];
  if (ticket.status !== "PENDING") {
    events.push({
      id: `evt-${ticket.id}-processing`,
      ticketId: ticket.id,
      eventType: "TICKET_PROCESSING",
      payload: { worker: "python-worker-1" },
      createdAt: new Date(base + 5 * 60 * 1000).toISOString(),
    });
    events.push({
      id: `evt-${ticket.id}-classified`,
      ticketId: ticket.id,
      eventType: "TICKET_CLASSIFIED",
      payload: { category: ticket.category, priority: ticket.priority },
      createdAt: new Date(base + 9 * 60 * 1000).toISOString(),
    });
  }
  events.push({
    id: `evt-${ticket.id}-final-${index}`,
    ticketId: ticket.id,
    eventType:
      ticket.status === "FAILED" ? "TICKET_FAILED" : ticket.status === "RESOLVED" ? "TICKET_RESOLVED" : "TICKET_PROCESSED",
    payload: ticket.errorMessage ? { error: ticket.errorMessage } : { team: ticket.assignedTeam },
    createdAt: ticket.updatedAt ?? ticket.createdAt,
  });
  return events;
});
