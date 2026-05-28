import { api, mockDelay, USE_MOCK } from "@/api/client";
import { mockStreamStats } from "@/mocks/stream";
import type { StreamStats } from "@/types";

export async function getStreamStats(): Promise<StreamStats> {
  if (USE_MOCK) return mockDelay(mockStreamStats);
  const { data } = await api.get<StreamStats>("/api/stream/stats");
  return data;
}
