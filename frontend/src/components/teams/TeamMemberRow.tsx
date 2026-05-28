import { RotateCw, Trash2 } from "lucide-react";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { IntegrationBadge } from "@/components/teams/IntegrationBadge";
import type { TeamMember } from "@/types";

export function TeamMemberRow({ member, onResend }: { member: TeamMember; onResend: (inviteId: string) => void }) {
  const inviteStatus = member.inviteStatus ?? "accepted";
  return (
    <tr className="odd:bg-surface even:bg-surface-2/40">
      <td className="px-4 py-3 font-medium">{member.name}</td>
      <td className="px-4 py-3">
        <Badge className={member.role === "senior" ? "border-blue-500/30 bg-blue-500/15 text-blue-300" : "border-slate-600 bg-slate-800 text-slate-300"}>
          {member.role}
        </Badge>
      </td>
      <td className="px-4 py-3 text-slate-300">{member.emailAddress ?? "No email"}</td>
      <td className="px-4 py-3">
        {inviteStatus === "pending" ? (
          <span className="inline-flex items-center gap-2 text-amber-300"><span className="h-2 w-2 animate-pulse rounded-full bg-amber-300" /> Pending</span>
        ) : inviteStatus === "expired" ? (
          <span title="Invite expired after 48 hours" className="text-red-300">Expired</span>
        ) : (
          <span className="text-emerald-300">Active</span>
        )}
      </td>
      <td className="px-4 py-3">
        <div className="flex flex-wrap gap-2">
          <IntegrationBadge provider="slack" connected={Boolean(member.integrations.slackId)} />
          <IntegrationBadge provider="teams" connected={Boolean(member.integrations.teamsId)} />
        </div>
      </td>
      <td className="px-4 py-3">
        <div className="flex gap-2">
          {member.inviteId && inviteStatus !== "accepted" ? (
            <Button variant="secondary" className="h-8 px-3" onClick={() => onResend(member.inviteId ?? "")}>
              <RotateCw className="h-4 w-4" /> Resend
            </Button>
          ) : null}
          <Button variant="ghost" className="h-8 w-8 p-0" aria-label="Delete member"><Trash2 className="h-4 w-4" /></Button>
        </div>
      </td>
    </tr>
  );
}
