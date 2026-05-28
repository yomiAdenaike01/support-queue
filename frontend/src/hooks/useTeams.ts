import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  type InviteMemberInput,
  createTeam,
  deleteTeam,
  getSlackChannels,
  getTeam,
  getTeams,
  getTeamsChannels,
  inviteMember,
  resendInvite,
  testIntegration,
} from "@/api/teams";

export function useTeams() {
  return useQuery({ queryKey: ["teams"], queryFn: getTeams, staleTime: 15_000 });
}

export function useTeam(id: string) {
  return useQuery({ queryKey: ["team", id], queryFn: () => getTeam(id), enabled: Boolean(id), refetchInterval: 3000 });
}

export function useSlackChannels(teamId: string) {
  return useQuery({ queryKey: ["slack-channels", teamId], queryFn: () => getSlackChannels(teamId), enabled: Boolean(teamId) });
}

export function useTeamsChannels(teamId: string) {
  return useQuery({ queryKey: ["teams-channels", teamId], queryFn: () => getTeamsChannels(teamId), enabled: Boolean(teamId) });
}

export function useTeamMutations() {
  const queryClient = useQueryClient();
  return {
    create: useMutation({ mutationFn: createTeam }),
    delete: useMutation({ mutationFn: deleteTeam, onSuccess: () => queryClient.invalidateQueries({ queryKey: ["teams"] }) }),
    invite: useMutation({
      mutationFn: ({ teamId, input }: { teamId: string; input: InviteMemberInput }) => inviteMember(teamId, input),
      onSuccess: (_, variables) => queryClient.invalidateQueries({ queryKey: ["team", variables.teamId] }),
    }),
    resend: useMutation({
      mutationFn: ({ teamId, inviteId }: { teamId: string; inviteId: string }) => resendInvite(teamId, inviteId),
    }),
    test: useMutation({
      mutationFn: ({ teamId, provider }: { teamId: string; provider: "slack" | "teams" }) => testIntegration(teamId, provider),
    }),
  };
}
