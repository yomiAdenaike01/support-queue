import { api } from "@/api/client";
import type { Metrics } from "@/types";

export async function getMetrics(): Promise<Metrics> {
  // if (USE_MOCK) return mockDelay(mockMetrics);
  const { data } = await api.get<Metrics>("/metrics/summary");
  return data;
}
const getStartOfToday = (): Date => {
  const now = new Date()
  return new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0)
}
export async function getTicketsOvertime(start: Date = getStartOfToday(), end: Date = new Date(Date.now() + 1000 * 60 * 60 * 24)): Promise<Array<{ hour: string, count: number }>> {
  // if (USE_MOCK) return mockDelay(mockMetrics);
  console.log({ start, end })
  const { data } = await api.get<Array<{ hour: string, count: number }>>(`/metrics/series/tickets-count?start_datetime=${start.toISOString()}&end_datetime=${end.toISOString()}`);
  return data;
}
