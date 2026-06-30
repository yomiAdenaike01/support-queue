import { createTicket } from "@/api/tickets";
import { api, mockDelay, USE_MOCK } from "@/api/client";
import type { Ticket } from "@/types";

const randomSubjects = [
  "Charged twice for my subscription",
  "Cannot access my account after password reset",
  "Delivery tracking has not updated",
  "Need to cancel before renewal",
  "App crashes when I upload a receipt",
  "Refund has not arrived after three weeks",
  "Received the wrong item",
  "Checkout shows a 500 error",
  "Need an invoice for finance",
  "Support button in the app is not working",
];

const randomMessages = [
  "I need help with this as soon as possible. I have already tried the help centre steps and still cannot resolve it.",
  "Can someone check my account and tell me what happened? This is blocking my team from moving forward today.",
  "I am following up because the issue is still happening and I have not received a clear answer yet.",
  "Please create a ticket for this and route it to the right team. I can provide more details if needed.",
  "This looks urgent from our side because it affects billing and customer access.",
];

const domains = ["example.com", "acme.example", "customer.test", "mail.test"];

export async function seedTickets(count: number): Promise<{ count: number }> {
  if (USE_MOCK) return mockDelay({ count });
  const { data } = await api.post<{ count: number }>("/simulator/seed", { count });
  return data;
}

export async function startSimulator(): Promise<{ started: boolean }> {
  await Promise.all(Array.from({ length: 20 }).map(() => {
    return randomTicket()
  }))
  return {
    started: true
  }
  // const { data } = await api.post<{ started: boolean }>("/simulator/start");
  // return data;
}

export async function stopSimulator(): Promise<{ stopped: boolean }> {
  if (USE_MOCK) return mockDelay({ stopped: true });
  const { data } = await api.post<{ stopped: boolean }>("/simulator/stop");
  return data;
}

export async function randomTicket(): Promise<Ticket> {
  return createTicket(generateRandomTicketInput());
}

export async function forceFailure(): Promise<Ticket> {
  if (USE_MOCK) {
    const ticket = await randomTicket();
    ticket.status = "FAILED";
    ticket.errorMessage = "Forced simulator failure";
    return ticket;
  }
  const { data } = await api.post<Ticket>("/simulator/force-failure");
  return data;
}

function generateRandomTicketInput() {
  const now = Date.now();
  const subject = sample(randomSubjects);
  const message = sample(randomMessages);
  const name = subject.toLowerCase().replace(/[^a-z]+/g, ".").replace(/^\.+|\.+$/g, "").slice(0, 18) || "customer";
  return {
    customerEmail: `${name}.${String(now).slice(-5)}@${sample(domains)}`,
    subject,
    message,
  };
}

function sample<T>(items: T[]): T {
  return items[Math.floor(Math.random() * items.length)];
}
