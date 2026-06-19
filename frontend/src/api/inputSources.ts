import { api, mockDelay, USE_MOCK } from "@/api/client";
import { mockInputSources } from "@/mocks/inputSources";
import type { TicketInputSourceConfig } from "@/types";

export async function getInputSources(): Promise<TicketInputSourceConfig[]> {
  // if (USE_MOCK) return
  return mockDelay(mockInputSources);
  // return []
  // const { data } = await api.get<TicketInputSourceConfig[]>("/input-sources");
  // return data;
}

export async function updateInputSource(input: TicketInputSourceConfig): Promise<TicketInputSourceConfig> {
  if (USE_MOCK) {
    const index = mockInputSources.findIndex((source) => source.id === input.id);
    if (index >= 0) mockInputSources[index] = input;
    return mockDelay(input);
  }
  const { data } = await api.patch<TicketInputSourceConfig>(`/input-sources/${input.id}`, input);
  return data;
}
