import { ArrowRight, Trash2 } from "lucide-react";
import { Link } from "react-router-dom";
import { CategoryBadge } from "@/components/teams/CategoryBadge";
import { IntegrationBadge } from "@/components/teams/IntegrationBadge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import type { Team } from "@/types";
import { initials } from "@/utils/format";

const avatarColors = ["bg-blue-500", "bg-emerald-500", "bg-amber-500", "bg-purple-500", "bg-rose-500"];

export function TeamCard({ team, onDelete }: { team: Team; onDelete: (team: Team) => void }) {
  const seniorCount = team.members.filter((member) => member.role === "senior").length;
  return (
    <Card>
      <div className="mb-4 flex items-start justify-between gap-4">
        <div>
          <h2 className="text-lg font-semibold">{team.department}</h2>
          <div className="mt-2 flex flex-wrap gap-2">{team.categories.map((category) => <CategoryBadge key={category} category={category} />)}</div>
        </div>
        <Button variant="ghost" className="h-9 w-9 p-0" onClick={() => onDelete(team)} aria-label="Delete team"><Trash2 className="h-4 w-4" /></Button>
      </div>
      <div className="mb-4 flex flex-wrap gap-2">
        <IntegrationBadge provider="slack" connected={Boolean(team.integrations.slackWebhook)} />
        <IntegrationBadge provider="teams" connected={Boolean(team.integrations.teamsWebhook)} />
      </div>
      <div className="mb-5 text-sm text-slate-400">{team.members.length - seniorCount} agents, {seniorCount} senior</div>
      <div className="mb-5 flex items-center">
        {team.members.slice(0, 3).map((member) => (
          <div
            key={member.id}
            title={`${member.name} - ${member.role}`}
            className={`-ml-1 grid h-9 w-9 place-items-center rounded-full border-2 border-surface text-xs font-bold text-white ${avatarColors[member.name.length % avatarColors.length]}`}
          >
            {initials(member.name)}
          </div>
        ))}
        {team.members.length > 3 ? <span className="ml-2 text-sm text-slate-400">+{team.members.length - 3} more</span> : null}
      </div>
      <Link to={`/teams/${team.id}`} className="inline-flex items-center gap-2 text-sm font-semibold text-blue-300 hover:text-blue-200">
        View Team <ArrowRight className="h-4 w-4" />
      </Link>
    </Card>
  );
}
