import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { CategoryBadge } from "@/components/teams/CategoryBadge";
import { useTeamMutations } from "@/hooks/useTeams";
import type { TicketCategory } from "@/types";

const categories: TicketCategory[] = ["BILLING", "TECHNICAL", "DELIVERY", "SUBSCRIPTIONS", "GENERAL", "ACCOUNT_ACCESS", "CANCELLATION"];

export function CreateTeam() {
  const [department, setDepartment] = useState("");
  const [selected, setSelected] = useState<TicketCategory[]>(["BILLING"]);
  const [slackWebhook, setSlackWebhook] = useState("");
  const [slackChannelName, setSlackChannelName] = useState("");
  const [teamsWebhook, setTeamsWebhook] = useState("");
  const [teamsChannelName, setTeamsChannelName] = useState("");
  const mutations = useTeamMutations();
  const navigate = useNavigate();
  const valid = department.length >= 3 && department.length <= 100 && selected.length > 0;

  const toggle = (category: TicketCategory) => setSelected((current) => current.includes(category) ? current.filter((item) => item !== category) : [...current, category]);
  const submit = async () => {
    const team = await mutations.create.mutateAsync({
      department,
      categories: selected,
      integrations: {
        slackWebhook: slackWebhook || null,
        slackChannelId: null,
        slackChannelName: slackChannelName || null,
        teamsWebhook: teamsWebhook || null,
        teamsChannelId: null,
        teamsChannelName: teamsChannelName || null,
      },
    });
    navigate(`/teams/${team.id}`);
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Create Team</h1>
      <div className="grid gap-6 xl:grid-cols-2">
        <Card>
          <h2 className="mb-4 text-lg font-semibold">Team Details</h2>
          <label className="block text-sm font-medium">Department Name
            <input value={department} onChange={(event) => setDepartment(event.target.value)} className="mt-2 w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2" />
          </label>
          <div className="mt-5">
            <div className="mb-3 text-sm font-medium">Handles Categories</div>
            <div className="grid gap-2 sm:grid-cols-2">
              {categories.map((category) => (
                <label key={category} className="flex items-center gap-3 rounded-lg bg-surface-2 p-3">
                  <input type="checkbox" checked={selected.includes(category)} onChange={() => toggle(category)} />
                  <CategoryBadge category={category} />
                </label>
              ))}
            </div>
          </div>
        </Card>
        <Card>
          <h2 className="mb-4 text-lg font-semibold">Integrations</h2>
          <IntegrationFields title="Slack Integration" webhook={slackWebhook} setWebhook={setSlackWebhook} channel={slackChannelName} setChannel={setSlackChannelName} />
          <div className="my-6 border-t border-slate-800" />
          <IntegrationFields title="Microsoft Teams Integration" webhook={teamsWebhook} setWebhook={setTeamsWebhook} channel={teamsChannelName} setChannel={setTeamsChannelName} />
        </Card>
      </div>
      <div className="flex justify-end gap-2">
        <Button variant="secondary" onClick={() => navigate("/teams")}>Cancel</Button>
        <Button disabled={!valid} loading={mutations.create.isPending} onClick={() => void submit()}>Create Team</Button>
      </div>
    </div>
  );
}

function IntegrationFields({ title, webhook, setWebhook, channel, setChannel }: { title: string; webhook: string; setWebhook: (value: string) => void; channel: string; setChannel: (value: string) => void }) {
  return (
    <div className="space-y-3">
      <h3 className="font-semibold">{title}</h3>
      <input value={webhook} onChange={(event) => setWebhook(event.target.value)} placeholder="Webhook URL" className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2" />
      <input value={channel} onChange={(event) => setChannel(event.target.value)} placeholder="Channel name" className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2" />
      <Button variant="secondary" type="button">Test Connection</Button>
    </div>
  );
}
