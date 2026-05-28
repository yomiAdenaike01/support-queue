import type { Metrics } from "@/types";

export const mockMetrics: Metrics = {
  totalTickets: 1240,
  pendingTickets: 23,
  processingTickets: 4,
  processedTickets: 1180,
  failedTickets: 33,
  resolvedTickets: 156,
  averageProcessingTimeMs: 4200,
  streamPendingMessages: 4,
  deadLetterCount: 7,
};
