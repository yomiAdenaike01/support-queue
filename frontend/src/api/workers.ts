import { api, mockDelay, USE_MOCK } from "@/api/client";
import { mockWorkers } from "@/mocks/workers";
import type { Worker } from "@/types";

export async function getWorkers(): Promise<Worker[]> {
  if (USE_MOCK) return mockDelay(mockWorkers.map((worker) => ({ ...worker, lastHeartbeat: new Date().toISOString() })));
  const { data } = await api.get<Worker[]>("/api/workers");
  return data;
}
