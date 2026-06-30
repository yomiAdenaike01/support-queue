import { api, mockDelay, USE_MOCK } from "@/api/client";
import { mockEvents } from "@/mocks/events";
import { mockTickets } from "@/mocks/tickets";
import type { Message, Ticket, TicketCategory, TicketEvent, TicketFilters, TicketPriority } from "@/types";

export interface CreateTicketInput {
  customerEmail: string;
  subject: string;
  message: string;
}

export interface AddTicketMessageInput {
  ticketId: string;
  customerEmail: string;
  content: string;
  role: Message["role"];
}

export interface UpdateTicketClassificationInput {
  ticketId: string;
  priority: TicketPriority | null;
  category: TicketCategory | null;
  assignedTeam: string | null;
  notifyTeam: boolean;
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
      events: [],
    };
    mockTickets.unshift(ticket);
    return mockDelay(ticket);
  }
  const { data } = await api.post<Ticket>("/tickets", input);
  return data;
}

export async function addTicketMessage(input: AddTicketMessageInput): Promise<Message> {
  if (USE_MOCK) {
    const ticket = await getTicket(input.ticketId);
    const message: Message = {
      id: `msg-${Date.now()}`,
      content: input.content,
      role: input.role,
      createdAt: new Date().toISOString(),
    };
    ticket.messages.push(message);
    ticket.updatedAt = message.createdAt;
    return mockDelay(message);
  }
  const { data } = await api.post<Message>(`/tickets/${input.ticketId}/message`, {
    content: input.content,
    customerEmail: input.customerEmail,
    role: input.role,
  });
  return data;
}

export async function reprocessTicket(id: string): Promise<Ticket> {
  if (USE_MOCK) {
    const ticket = await getTicket(id);
    ticket.status = "PROCESSING";
    ticket.retryCount += 1;
    ticket.updatedAt = new Date().toISOString();
    return mockDelay(ticket);
  }
  await api.post(`/tickets/${id}/reprocess`);
  return getTicket(id);
}

export async function rerunTicketPipeline(id: string): Promise<void> {
  if (USE_MOCK) {
    const ticket = await getTicket(id);
    ticket.status = "PROCESSING";
    ticket.retryCount += 1;
    ticket.updatedAt = new Date().toISOString();
    await mockDelay(undefined);
    return;
  }
  await api.post(`/tickets/${id}/pipeline`);
}

export async function updateTicketClassification(input: UpdateTicketClassificationInput): Promise<Ticket> {
  if (USE_MOCK) {
    const ticket = await getTicket(input.ticketId);
    ticket.priority = input.priority;
    ticket.category = input.category;
    ticket.assignedTeam = input.assignedTeam;
    ticket.updatedAt = new Date().toISOString();
    return mockDelay(ticket);
  }
  await api.patch(`/tickets/${input.ticketId}`, {
    priority: input.priority,
    category: input.category,
    assigned_teams: input.assignedTeam,
    notify_team: input.notifyTeam,
  });
  return getTicket(input.ticketId);
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
