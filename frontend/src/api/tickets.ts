import { api, mockDelay, USE_MOCK } from "@/api/client";
import { mockEvents } from "@/mocks/events";
import { mockTickets } from "@/mocks/tickets";
import type { Ticket, TicketEvent, TicketFilters } from "@/types";

export interface CreateTicketInput {
  customerEmail: string;
  subject: string;
  message: string;
}

export async function getTickets(filters: TicketFilters = {}): Promise<Ticket[]> {
  // if (USE_MOCK) {
  //   const search = filters.search?.toLowerCase().trim();
  //   return mockDelay(
  //     mockTickets.filter((ticket) => {
  //       const matchesStatus = !filters.status || filters.status === "ALL" || ticket.status === filters.status;
  //       const matchesPriority = !filters.priority || filters.priority === "ALL" || ticket.priority === filters.priority;
  //       const matchesCategory = !filters.category || filters.category === "ALL" || ticket.category === filters.category;
  //       const matchesSearch =
  //         !search ||
  //         ticket.subject.toLowerCase().includes(search) ||
  //         ticket.customerEmail.toLowerCase().includes(search) ||
  //         ticket.id.toLowerCase().includes(search);
  //       return matchesStatus && matchesPriority && matchesCategory && matchesSearch;
  //     }),
  //   );
  // }
  const { data } = await api.get<Ticket[]>("/tickets", { params: filters });
  return data;
}

export async function getTicket(id: string): Promise<Ticket> {
  if (USE_MOCK) {
    const ticket = mockTickets.find((item) => item.id === id);
    if (!ticket) throw new Error("Ticket not found");
    return mockDelay(ticket);
  }
  const { data } = await api.get<Ticket>(`/tickets/${id}`);
  return data;
}

export async function getTicketEvents(id: string): Promise<TicketEvent[]> {
  if (USE_MOCK) return mockDelay(mockEvents.filter((event) => event.ticketId === id));
  const { data } = await api.get<TicketEvent[]>(`/tickets/${id}/events`);
  return data;
}

export async function createTicket(input: CreateTicketInput): Promise<Ticket> {
  if (USE_MOCK) {
    const now = new Date().toISOString();
    const ticket: Ticket = {
      id: `TCK-${String(mockTickets.length + 1).padStart(5, "0")}`,
      customerEmail: input.customerEmail,
      subject: input.subject,
      status: "PENDING",
      priority: null,
      category: null,
      assignedTeam: null,
      suggestedResponse: null,
      retryCount: 0,
      errorMessage: null,
      urgencyFlag: false,
      urgencyReason: null,
      sentimentScore: null,
      sentimentLabel: null,
      messages: [{ id: `msg-${Date.now()}`, role: "customer", content: input.message, createdAt: now }],
      createdAt: now,
      updatedAt: null,
      processedAt: null,
    };
    mockTickets.unshift(ticket);
    return mockDelay(ticket);
  }
  const { data } = await api.post<Ticket>("/tickets", input);
  return data;
}

export async function reprocessTicket(id: string): Promise<Ticket> {
  if (USE_MOCK) {
    const ticket = await getTicket(id);
    ticket.status = "PROCESSING";
    ticket.retryCount += 1;
    return mockDelay(ticket);
  }
  const { data } = await api.post<Ticket>(`/tickets/${id}/reprocess`);
  return data;
}

export async function resolveTicket(id: string): Promise<Ticket> {
  if (USE_MOCK) {
    const ticket = await getTicket(id);
    ticket.status = "RESOLVED";
    ticket.updatedAt = new Date().toISOString();
    return mockDelay(ticket);
  }
  const { data } = await api.post<Ticket>(`/tickets/${id}/resolve`);
  return data;
}
