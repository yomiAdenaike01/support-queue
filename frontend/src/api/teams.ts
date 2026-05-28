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

export interface InviteMemberInput {
  name: string;
  email: string;
  role: TeamMemberRole;
}

export async function getTeams(): Promise<Team[]> {
  if (USE_MOCK) return mockDelay(mockTeams);
  const { data } = await api.get<Team[]>("/api/teams");
  return data;
}

export async function getTeam(id: string): Promise<Team> {
  if (USE_MOCK) {
    const team = mockTeams.find((item) => item.id === id);
    if (!team) throw new Error("Team not found");
    return mockDelay(team);
  }
  const { data } = await api.get<Team>(`/api/teams/${id}`);
  return data;
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
  const { data } = await api.post<Team>("/api/teams", input);
  return data;
}

export async function deleteTeam(id: string): Promise<{ deleted: boolean }> {
  if (USE_MOCK) return mockDelay({ deleted: true });
  const { data } = await api.delete<{ deleted: boolean }>(`/api/teams/${id}`);
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
  const { data } = await api.post<MemberInvite>(`/api/teams/${teamId}/invites`, input);
  return data;
}

export async function resendInvite(teamId: string, inviteId: string): Promise<MemberInvite> {
  if (USE_MOCK) {
    const invite = mockInvites.find((item) => item.id === inviteId) ?? mockInvites[0];
    return mockDelay({ ...invite, teamId, invitedAt: new Date().toISOString(), status: "pending" });
  }
  const { data } = await api.post<MemberInvite>(`/api/teams/${teamId}/invites/${inviteId}/resend`);
  return data;
}

export async function getSlackChannels(teamId: string): Promise<SlackChannel[]> {
  if (USE_MOCK) return mockDelay(mockSlackChannels);
  const { data } = await api.get<SlackChannel[]>("/api/integrations/slack/channels", { params: { team_id: teamId } });
  return data;
}

export async function getTeamsChannels(teamId: string): Promise<TeamsChannel[]> {
  if (USE_MOCK) return mockDelay(mockTeamsChannels);
  const { data } = await api.get<TeamsChannel[]>("/api/integrations/teams/teams-list", { params: { team_id: teamId } });
  return data;
}

export async function testIntegration(teamId: string, provider: "slack" | "teams"): Promise<{ success: boolean }> {
  if (USE_MOCK) return mockDelay({ success: true });
  const { data } = await api.post<{ success: boolean }>(`/api/integrations/${provider}/test`, undefined, {
    params: { team_id: teamId },
  });
  return data;
}
