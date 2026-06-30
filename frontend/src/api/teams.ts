import { api, mockDelay, USE_MOCK } from "@/api/client";
import { mockInvites, mockSlackChannels, mockTeams, mockTeamsChannels } from "@/mocks/teams";
import type {
  CreateTeamPayload,
  MemberInvite,
  SlackChannel,
  Team,
  TeamMember,
  TeamMemberRole,
  TeamsChannel,
} from "@/types";

interface BackendTeam {
  id: string;
  department: string;
  integrations?: Record<string, unknown> | null;
  members?: Array<{
    id: string;
    name: string;
    role: TeamMemberRole;
    email_address?: string | null;
    phone_number?: string | null;
  }>;
}

export interface InviteMemberInput {
  name: string;
  email: string;
  role: TeamMemberRole;
}

export async function getTeams(): Promise<Team[]> {
  if (USE_MOCK) return mockDelay(mockTeams);
  const { data } = await api.get<BackendTeam[]>("/team");
  return data.map(toTeam);
}

export async function getTeam(id: string): Promise<Team> {
  if (USE_MOCK) {
    const team = mockTeams.find((item) => item.id === id);
    if (!team) throw new Error("Team not found");
    return mockDelay(team);
  }
  const { data } = await api.get<Team>(`/teams/${id}`);
  return data;
}

function toTeam(team: BackendTeam): Team {
  return {
    id: team.id,
    department: team.department,
    categories: [],
    integrations: {
      slackWebhook: null,
      slackChannelId: null,
      slackChannelName: null,
      teamsWebhook: null,
      teamsChannelId: null,
      teamsChannelName: null,
    },
    members: (team.members ?? []).map((member) => ({
      id: member.id,
      name: member.name,
      role: member.role,
      emailAddress: member.email_address ?? null,
      phoneNumber: member.phone_number ?? null,
      integrations: { slackId: null, teamsId: null },
    })),
    createdAt: new Date().toISOString(),
  };
}

export async function createTeam(input: CreateTeamPayload): Promise<Team> {
  if (USE_MOCK) {
    const team: Team = {
      id: `team-${String(mockTeams.length + 1).padStart(3, "0")}`,
      ...input,
      members: [],
      activeTicketCount: 0,
      createdAt: new Date().toISOString(),
    };
    mockTeams.push(team);
    return mockDelay(team);
  }
  const { data } = await api.post<Team>("/teams", input);
  return data;
}

export async function deleteTeam(id: string): Promise<{ deleted: boolean }> {
  if (USE_MOCK) return mockDelay({ deleted: true });
  const { data } = await api.delete<{ deleted: boolean }>(`/teams/${id}`);
  return data;
}

export async function inviteMember(teamId: string, input: InviteMemberInput): Promise<MemberInvite> {
  if (USE_MOCK) {
    const invite: MemberInvite = {
      id: `inv-${Date.now()}`,
      teamId,
      email: input.email,
      name: input.name,
      role: input.role,
      status: "pending",
      invitedAt: new Date().toISOString(),
      expiresAt: new Date(Date.now() + 48 * 60 * 60 * 1000).toISOString(),
      acceptedAt: null,
    };
    mockInvites.push(invite);
    const team = mockTeams.find((item) => item.id === teamId);
    const member: TeamMember = {
      id: `mem-${Date.now()}`,
      name: input.name,
      emailAddress: input.email,
      phoneNumber: null,
      role: input.role,
      integrations: { slackId: null, teamsId: null },
      inviteStatus: "pending",
      inviteId: invite.id,
    };
    team?.members.push(member);
    return mockDelay(invite);
  }
  const { data } = await api.post<MemberInvite>(`/teams/${teamId}/invites`, input);
  return data;
}

export async function resendInvite(teamId: string, inviteId: string): Promise<MemberInvite> {
  if (USE_MOCK) {
    const invite = mockInvites.find((item) => item.id === inviteId) ?? mockInvites[0];
    return mockDelay({ ...invite, teamId, invitedAt: new Date().toISOString(), status: "pending" });
  }
  const { data } = await api.post<MemberInvite>(`/teams/${teamId}/invites/${inviteId}/resend`);
  return data;
}

export async function getSlackChannels(teamId: string): Promise<SlackChannel[]> {
  if (USE_MOCK) return mockDelay(mockSlackChannels);
  const { data } = await api.get<SlackChannel[]>("/integrations/slack/channels", { params: { team_id: teamId } });
  return data;
}

export async function getTeamsChannels(teamId: string): Promise<TeamsChannel[]> {
  if (USE_MOCK) return mockDelay(mockTeamsChannels);
  const { data } = await api.get<TeamsChannel[]>("/integrations/teams/teams-list", { params: { team_id: teamId } });
  return data;
}

export async function testIntegration(teamId: string, provider: "slack" | "teams"): Promise<{ success: boolean }> {
  if (USE_MOCK) return mockDelay({ success: true });
  const { data } = await api.post<{ success: boolean }>(`/integrations/${provider}/test`, undefined, {
    params: { team_id: teamId },
  });
  return data;
}
