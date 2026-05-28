import { ArrowLeft, Lock, Plus } from "lucide-react";
import { useMemo, useState } from "react";
import type { ReactNode } from "react";
import { Link, useParams } from "react-router-dom";
import { CategoryBadge } from "@/components/teams/CategoryBadge";
import { SlackIcon, TeamsIcon } from "@/components/teams/IntegrationBadge";
import { TeamMemberRow } from "@/components/teams/TeamMemberRow";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Modal } from "@/components/ui/Modal";
import { SkeletonRows } from "@/components/ui/Skeleton";
import { Table } from "@/components/ui/Table";
import { useToast } from "@/components/ui/Toast";
import { useSlackChannels, useTeam, useTeamMutations, useTeamsChannels } from "@/hooks/useTeams";
import { openOAuthPopup } from "@/utils/oauth";
import { maskWebhook } from "@/utils/format";
import type { TeamMemberRole } from "@/types";

export function TeamDetail() {
  const { id = "" } = useParams();
  const team = useTeam(id);
  const slackChannels = useSlackChannels(id);
  const teamsChannels = useTeamsChannels(id);
  const mutations = useTeamMutations();
  const toast = useToast();
  const [inviteOpen, setInviteOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [invite, setInvite] = useState({ name: "", email: "", role: "agent" as TeamMemberRole });
  const filteredSlack = useMemo(() => (slackChannels.data ?? []).filter((channel) => channel.name.includes(search)), [slackChannels.data, search]);

  if (team.isLoading) return <SkeletonRows rows={8} />;
  if (!team.data) return <div className="text-red-200">Team not found.</div>;
  const item = team.data;

  const sendInvite = async () => {
    await mutations.invite.mutateAsync({ teamId: item.id, input: invite });
    toast.push(`Invitation sent to ${invite.email}`, "success");
    setInviteOpen(false);
  };

  const connect = (provider: "slack" | "teams") => {
    openOAuthPopup(`/api/integrations/${provider}/authorize?team_id=${item.id}`, () => toast.push(`${provider} connected`, "success"), (err) => toast.push(err, "error"));
  };

  return (
    <div className="space-y-6">
      <Link to="/teams" className="inline-flex items-center gap-2 text-sm text-slate-300 hover:text-white"><ArrowLeft className="h-4 w-4" /> Back to Teams</Link>
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">{item.department}</h1>
          <div className="mt-2 flex flex-wrap gap-2">{item.categories.map((category) => <CategoryBadge key={category} category={category} />)}</div>
        </div>
        <Button variant="danger">Delete Team</Button>
      </div>
      <div className="grid gap-6 xl:grid-cols-[3fr_2fr]">
        <Card>
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold">Members</h2>
            <Button onClick={() => setInviteOpen(true)}><Plus className="h-4 w-4" /> Add Member</Button>
          </div>
          <Table>
            <thead className="bg-surface-2 text-left text-xs uppercase tracking-wide text-slate-400">
              <tr>{["Name", "Role", "Email", "Status", "Integrations", "Actions"].map((header) => <th key={header} className="px-4 py-3">{header}</th>)}</tr>
            </thead>
            <tbody className="divide-y divide-slate-800">
              {item.members.map((member) => (
                <TeamMemberRow
                  key={member.id}
                  member={member}
                  onResend={(inviteId) => {
                    void mutations.resend.mutateAsync({ teamId: item.id, inviteId }).then(() => toast.push("Invitation resent", "success"));
                  }}
                />
              ))}
            </tbody>
          </Table>
        </Card>
        <div className="space-y-6">
          <Card>
            <h2 className="mb-4 text-lg font-semibold">Slack Integration</h2>
            {item.integrations.slackWebhook ? (
              <IntegrationConfigured icon={<SlackIcon />} title="Connected to Slack" channel={item.integrations.slackChannelName ?? "No channel"} webhook={item.integrations.slackWebhook} onTest={() => void mutations.test.mutateAsync({ teamId: item.id, provider: "slack" }).then(() => toast.push("Slack test sent", "success"))} />
            ) : (
              <Button variant="brand-slack" onClick={() => connect("slack")}><SlackIcon /> Connect with Slack</Button>
            )}
            <div className="mt-5">
              <input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search Slack channels" className="mb-3 w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2 text-sm" />
              <div className="space-y-2">
                {filteredSlack.map((channel) => <div key={channel.id} className="flex items-center justify-between rounded-lg bg-surface-2 p-2 text-sm"><span>{channel.isPrivate ? <Lock className="mr-1 inline h-3 w-3" /> : null}#{channel.name}</span><span className="text-slate-400">{channel.memberCount} members</span></div>)}
              </div>
            </div>
          </Card>
          <Card>
            <h2 className="mb-4 text-lg font-semibold">Microsoft Teams Integration</h2>
            {item.integrations.teamsWebhook ? (
              <IntegrationConfigured icon={<TeamsIcon />} title="Connected to Microsoft Teams" channel={item.integrations.teamsChannelName ?? "No channel"} webhook={item.integrations.teamsWebhook} onTest={() => void mutations.test.mutateAsync({ teamId: item.id, provider: "teams" }).then(() => toast.push("Teams test sent", "success"))} />
            ) : (
              <Button variant="brand-teams" onClick={() => connect("teams")}><TeamsIcon /> Connect with Microsoft Teams</Button>
            )}
            <div className="mt-5 space-y-2">
              {(teamsChannels.data ?? []).map((channel) => <div key={channel.id} className="rounded-lg bg-surface-2 p-2 text-sm">{channel.teamName} › {channel.name}</div>)}
            </div>
          </Card>
        </div>
      </div>
      <Modal title="Invite Team Member" open={inviteOpen} onClose={() => setInviteOpen(false)}>
        <div className="space-y-3">
          <input value={invite.name} onChange={(event) => setInvite({ ...invite, name: event.target.value })} placeholder="Name" className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2" />
          <input value={invite.email} onChange={(event) => setInvite({ ...invite, email: event.target.value })} placeholder="Email address" className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2" />
          <select value={invite.role} onChange={(event) => setInvite({ ...invite, role: event.target.value as TeamMemberRole })} className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2">
            <option value="agent">Agent</option>
            <option value="senior">Senior</option>
          </select>
          <p className="text-sm text-slate-400">An invitation email will be sent to this address.</p>
          <div className="flex justify-end gap-2">
            <Button variant="secondary" onClick={() => setInviteOpen(false)}>Cancel</Button>
            <Button disabled={!invite.name || !invite.email} loading={mutations.invite.isPending} onClick={() => void sendInvite()}>Send Invite</Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

function IntegrationConfigured({ icon, title, channel, webhook, onTest }: { icon: ReactNode; title: string; channel: string; webhook: string; onTest: () => void }) {
  return (
    <div className="space-y-3 text-sm">
      <div className="flex items-center gap-2 text-emerald-300">{icon} {title}</div>
      <div className="text-slate-300">{channel}</div>
      <div className="font-mono text-xs text-slate-500">{maskWebhook(webhook)}</div>
      <div className="flex gap-2">
        <Button variant="secondary" onClick={onTest}>Test</Button>
        <Button variant="ghost">Change</Button>
        <Button variant="ghost">Disconnect</Button>
      </div>
    </div>
  );
}
