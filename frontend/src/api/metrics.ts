import { api, mockDelay, USE_MOCK } from "@/api/client";
import { mockMetrics } from "@/mocks/metrics";
import type { Metrics } from "@/types";

export async function getMetrics(): Promise<Metrics> {
  if (USE_MOCK) return mockDelay(mockMetrics);
  const { data } = await api.get<Metrics>("/api/metrics");
  return data;
}
