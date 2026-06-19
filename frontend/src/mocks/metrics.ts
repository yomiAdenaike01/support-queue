import type { Metrics } from "@/types";

export const mockMetrics: Metrics = {
  totalTickets: { today: 1240, yesterday: 1160, percentageDifference: 6.9 },
  pendingTickets: { today: 23, yesterday: 31, percentageDifference: -25.8 },
  processingTickets: { today: 4, yesterday: 5, percentageDifference: -20 },
  processedTickets: { today: 1180, yesterday: 1092, percentageDifference: 8.1 },
  failedTickets: { today: 33, yesterday: 41, percentageDifference: -19.5 },
  resolvedTickets: { today: 156, yesterday: 132, percentageDifference: 18.2 },
  averageProcessingTimeMs: { today: 4200, yesterday: 4700, percentageDifference: -10.6 },
  streamPendingMessages: { today: 4, yesterday: 9, percentageDifference: -55.6 },
  deadLetterCount: { today: 7, yesterday: 12, percentageDifference: -41.7 },
};
