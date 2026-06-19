import { api } from "@/api/client";
import type { Metrics } from "@/types";

export async function getMetrics(): Promise<Metrics> {
  // if (USE_MOCK) return mockDelay(mockMetrics);
  const { data } = await api.get<Metrics>("/metrics/summary");
  return data;
}
