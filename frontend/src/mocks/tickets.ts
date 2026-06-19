import type { Ticket, TicketCategory, TicketPriority, TicketStatus } from "@/types";

const subjects = [
  "Charged twice for order #4821",
  "Can't log into my account",
  "Parcel not updated in 4 days",
  "Cancel my subscription",
  "App keeps crashing",
  "Refund after 3 weeks",
  "Wrong item received",
  "500 error at checkout",
  "Need invoice for last month",
  "Password reset email never arrived",
];

const statuses: TicketStatus[] = ["PENDING", "PROCESSING", "PROCESSED", "FAILED", "RESOLVED"];
const priorities: TicketPriority[] = ["LOW", "MEDIUM", "HIGH", "URGENT"];
const categories: TicketCategory[] = [
  "BILLING",
  "ACCOUNT_ACCESS",
  "TECHNICAL",
  "DELIVERY",
  "CANCELLATION",
  "SUBSCRIPTIONS",
  "GENERAL",
];

export const mockTickets: Ticket[] = Array.from({ length: 30 }, (_, index) => {
  const status = statuses[index % statuses.length];
  const priority = status === "PENDING" ? null : priorities[index % priorities.length];
  const category = status === "PENDING" ? null : categories[index % categories.length];
  const created = new Date(Date.now() - (index + 1) * 37 * 60 * 1000).toISOString();
  const urgent = priority === "URGENT" || index % 11 === 0;
  return {
    id: `TCK-${String(index + 1).padStart(5, "0")}`,
    customerEmail: `customer${index + 1}@example.com`,
    subject: subjects[index % subjects.length],
    status,
    priority,
    category,
    assignedTeam: category ? `${category.toLowerCase().replace("_", " ")} team` : null,
    suggestedResponse:
      status === "PROCESSED" || status === "RESOLVED"
        ? "Thanks for contacting SupportOps. We reviewed the issue and prepared the next steps for your account. Please reply if anything still looks incorrect."
        : null,
    retryCount: status === "FAILED" ? (index % 3) + 1 : index % 9 === 0 ? 1 : 0,
    errorMessage: status === "FAILED" ? "Worker timed out while waiting for model response" : null,
    urgencyFlag: urgent,
    urgencyReason: urgent ? "Customer mentions financial impact or blocked account access." : null,
    sentimentScore: status === "PENDING" ? null : Math.max(0.12, 0.92 - index * 0.021),
    sentimentLabel: status === "PENDING" ? null : index % 3 === 0 ? "HIGH" : index % 3 === 1 ? "MEDIUM" : "LOW",
    messages: [
      {
        id: `msg-${index}-1`,
        role: "customer",
        content: `Hi, I need help with this issue: ${subjects[index % subjects.length]}.`,
        createdAt: created,
      },
      {
        id: `msg-${index}-2`,
        role: "assistant",
        content: "Thanks for the context. I am checking the account and order history now.",
        createdAt: new Date(new Date(created).getTime() + 8 * 60 * 1000).toISOString(),
      },
      {
        id: `msg-${index}-3`,
        role: "customer",
        content: "Please treat this as urgent. I have already waited several days for a response.",
        createdAt: new Date(new Date(created).getTime() + 16 * 60 * 1000).toISOString(),
      },
      {
        id: `msg-${index}-4`,
        role: "assistant",
        content: "I have added the details to the case and the specialist team will follow up with the resolution.",
        createdAt: new Date(new Date(created).getTime() + 24 * 60 * 1000).toISOString(),
      },
    ],
    createdAt: created,
    updatedAt: new Date(new Date(created).getTime() + 31 * 60 * 1000).toISOString(),
    processedAt:
      status === "PROCESSED" || status === "RESOLVED"
        ? new Date(new Date(created).getTime() + 18 * 60 * 1000).toISOString()
        : null,
    events: [],
  };
});
