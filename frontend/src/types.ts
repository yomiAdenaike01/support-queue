export type TicketStatus = "PENDING" | "PROCESSING" | "PROCESSED" | "FAILED" | "RESOLVED";
export type TicketPriority = "LOW" | "MEDIUM" | "HIGH" | "URGENT";
export type TicketCategory =
  | "BILLING"
  | "ACCOUNT_ACCESS"
  | "TECHNICAL"
  | "DELIVERY"
  | "CANCELLATION"
  | "SUBSCRIPTIONS"
  | "GENERAL";
export type WorkerStatus = "ONLINE" | "OFFLINE" | "DEGRADED";
export type SentimentLabel = "HIGH" | "MEDIUM" | "LOW";

export interface Message {
  id: string;
  content: string;
  role: "customer" | "assistant";
  createdAt: string;
}

export interface Ticket {
  id: string;
  customerEmail: string;
  subject: string;
  status: TicketStatus;
  priority: TicketPriority | null;
  category: TicketCategory | null;
  assignedTeam: string | null;
  suggestedResponse: string | null;
  retryCount: number;
  errorMessage: string | null;
  urgencyFlag: boolean;
  urgencyReason: string | null;
  sentimentScore: number | null;
  sentimentLabel: SentimentLabel | null;
  messages: Message[];
  createdAt: string;
  updatedAt: string | null;
  processedAt: string | null;
}

export interface TicketEvent {
  id: string;
  ticketId: string;
  eventType:
    | "TICKET_CREATED"
    | "TICKET_PROCESSING"
    | "TICKET_CLASSIFIED"
    | "TICKET_PROCESSED"
    | "TICKET_FAILED"
    | "TICKET_REPROCESSED"
    | "TICKET_RESOLVED";
  payload: Record<string, unknown> | null;
  createdAt: string;
}

export interface Worker {
  name: string;
  language: string;
  status: WorkerStatus;
  processedCount: number;
  failedCount: number;
  lastHeartbeat: string;
}

export interface Metrics {
  totalTickets: number;
  pendingTickets: number;
  processingTickets: number;
  processedTickets: number;
  failedTickets: number;
  resolvedTickets: number;
  averageProcessingTimeMs: number;
  streamPendingMessages: number;
  deadLetterCount: number;
}

export interface StreamStats {
  pendingMessages: number;
  deadLetterCount: number;
  totalProcessed: number;
  consumerGroups: number;
}

export type TeamMemberRole = "senior" | "agent";
export type IntegrationProvider = "slack" | "microsoft_teams";
export type IntegrationStatus = "connected" | "disconnected" | "pending" | "error";
export type InviteStatus = "pending" | "accepted" | "expired";

export interface TeamIntegrations {
  slackWebhook: string | null;
  slackChannelId: string | null;
  slackChannelName: string | null;
  teamsWebhook: string | null;
  teamsChannelId: string | null;
  teamsChannelName: string | null;
}

export interface TeamMemberIntegrations {
  slackId: string | null;
  teamsId: string | null;
}

export interface TeamMember {
  id: string;
  name: string;
  emailAddress: string | null;
  phoneNumber: string | null;
  role: TeamMemberRole;
  integrations: TeamMemberIntegrations;
  inviteStatus?: InviteStatus;
  inviteId?: string;
}

export interface SlackChannel {
  id: string;
  name: string;
  isPrivate: boolean;
  memberCount: number;
}

export interface TeamsChannel {
  id: string;
  name: string;
  teamName: string;
  isPrivate: boolean;
}

export interface OAuthIntegration {
  provider: IntegrationProvider;
  status: IntegrationStatus;
  connectedAt: string | null;
  connectedBy: string | null;
  workspaceName: string | null;
  workspaceIcon: string | null;
  selectedChannel: SlackChannel | TeamsChannel | null;
  error: string | null;
}

export interface Team {
  id: string;
  department: string;
  categories: TicketCategory[];
  integrations: TeamIntegrations;
  members: TeamMember[];
  createdAt: string;
  oauthIntegrations?: OAuthIntegration[];
  activeTicketCount?: number;
}

export interface CreateTeamPayload {
  department: string;
  categories: TicketCategory[];
  integrations: TeamIntegrations;
}

export interface CreateTeamMemberPayload {
  name: string;
  emailAddress: string | null;
  phoneNumber: string | null;
  role: TeamMemberRole;
  integrations: TeamMemberIntegrations;
}

export interface MemberInvite {
  id: string;
  teamId: string;
  email: string;
  name: string;
  role: TeamMemberRole;
  status: InviteStatus;
  invitedAt: string;
  expiresAt: string;
  acceptedAt: string | null;
}

export interface TicketFilters {
  status?: TicketStatus | "ALL";
  priority?: TicketPriority | "ALL";
  category?: TicketCategory | "ALL";
  search?: string;
  page?: number;
  limit?: number;
}
