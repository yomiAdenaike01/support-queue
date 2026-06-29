import { api, mockDelay, USE_MOCK } from "@/api/client";
import { mockInputSources } from "@/mocks/inputSources";
import type { TicketInputSourceConfig, TicketInputSourceStatus, TicketInputSourceType } from "@/types";

interface BackendInputSource {
  id: string;
  name: string;
  enabled: boolean;
  config: Record<string, string | boolean> | null;
  sourceType: string;
  teamId: string | null;
  department: string | null;
  connectionValue: string;
  status: string;
  createdAt: string;
}

interface SaveBackendInputSource {
  connectionValue: string;
  status: string;
  sourceType: TicketInputSourceType;
  name: string;
  config: Record<string, string | boolean>;
  enabled: boolean;
  teamId: string | null;
}

const uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

export async function getInputSources(): Promise<TicketInputSourceConfig[]> {
  if (USE_MOCK) return mockDelay(mockInputSources);
  const { data } = await api.get<BackendInputSource[]>("/input-sources");
  return mergeBackendSources(data);
}

export async function updateInputSource(input: TicketInputSourceConfig): Promise<TicketInputSourceConfig> {
  if (USE_MOCK) {
    const index = mockInputSources.findIndex((source) => source.id === input.id);
    if (index >= 0) mockInputSources[index] = input;
    return mockDelay(input);
  }
  const payload = toBackendPayload(input);
  const request = uuidPattern.test(input.id)
    ? api.patch<BackendInputSource>(`/input-sources/${input.id}`, payload)
    : api.post<BackendInputSource>("/input-sources", payload);
  const { data } = await request;
  return fromBackendSource(data, input);
}

function mergeBackendSources(sources: BackendInputSource[]): TicketInputSourceConfig[] {
  const byType = new Map(sources.map((source) => [source.sourceType, source]));
  const merged = mockInputSources.map((defaultSource) => {
    const backendSource = byType.get(defaultSource.type);
    return backendSource ? fromBackendSource(backendSource, defaultSource) : defaultSource;
  });
  const knownTypes = new Set(mockInputSources.map((source) => source.type));
  const extras = sources
    .filter((source) => !knownTypes.has(source.sourceType as TicketInputSourceType))
    .map((source) => fromBackendSource(source));
  return [...merged, ...extras];
}

function fromBackendSource(source: BackendInputSource, fallback?: TicketInputSourceConfig): TicketInputSourceConfig {
  const type = source.sourceType as TicketInputSourceType;
  return {
    id: source.id,
    type,
    name: source.name || fallback?.name || source.sourceType,
    description: fallback?.description ?? "External ticket input source.",
    status: toFrontendStatus(source.status, source.enabled, source.connectionValue),
    isEnabled: source.enabled,
    connectionLabel: fallback?.connectionLabel ?? "Connection",
    connectionValue: source.connectionValue,
    routingTeam: source.department,
    autoCreateTickets: Boolean(source.config?.autoCreateTickets ?? fallback?.autoCreateTickets ?? true),
    requiresAuth: Boolean(source.config?.requiresAuth ?? fallback?.requiresAuth ?? false),
    lastReceivedAt: fallback?.lastReceivedAt ?? null,
    createdToday: fallback?.createdToday ?? 0,
    settings: { ...(fallback?.settings ?? {}), ...(source.config ?? {}) },
  };
}

function toBackendPayload(source: TicketInputSourceConfig): SaveBackendInputSource {
  return {
    connectionValue: source.connectionValue,
    status: source.status,
    sourceType: source.type,
    name: source.name,
    enabled: source.isEnabled,
    teamId: null,
    config: {
      ...source.settings,
      autoCreateTickets: source.autoCreateTickets,
      requiresAuth: source.requiresAuth,
      routingTeam: source.routingTeam ?? "",
    },
  };
}

function toFrontendStatus(status: string, enabled: boolean, connectionValue: string): TicketInputSourceStatus {
  if (!enabled) return "disabled";
  if (status === "enabled" || status === "disabled" || status === "needs_setup") return status;
  return connectionValue ? "enabled" : "needs_setup";
}
