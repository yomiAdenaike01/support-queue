import { Plus } from "lucide-react";
import { Link } from "react-router-dom";
import { TeamCard } from "@/components/teams/TeamCard";
import { Button } from "@/components/ui/Button";
import { EmptyState } from "@/components/ui/EmptyState";
import { Modal } from "@/components/ui/Modal";
import { SkeletonRows } from "@/components/ui/Skeleton";
import { useToast } from "@/components/ui/Toast";
import { useTeamMutations, useTeams } from "@/hooks/useTeams";
import type { Team, TicketCategory } from "@/types";
import { useState } from "react";

const allCategories: TicketCategory[] = ["BILLING", "ACCOUNT_ACCESS", "TECHNICAL", "DELIVERY", "CANCELLATION", "SUBSCRIPTIONS", "GENERAL"];

export function Teams() {
  const teams = useTeams();
  const mutations = useTeamMutations();
  const toast = useToast();
  const [deleting, setDeleting] = useState<Team | null>(null);
  const unassigned = allCategories.filter((category) => !(teams.data ?? []).some((team) => team.categories.includes(category)));

  const confirmDelete = async () => {
    if (!deleting) return;
    await mutations.delete.mutateAsync(deleting.id);
    toast.push("Team deleted", "success");
    setDeleting(null);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Teams</h1>
          <p className="text-sm text-slate-400">Manage support teams, routing categories, and integrations.</p>
        </div>
        <Link to="/teams/new"><Button><Plus className="h-4 w-4" /> Create Team</Button></Link>
      </div>
      {unassigned.length > 0 ? <div className="rounded-lg border border-amber-500/30 bg-amber-500/10 p-3 text-sm text-amber-100">Warning: {unassigned.join(", ")} tickets have no assigned team and will not be routed.</div> : null}
      {teams.isLoading ? <SkeletonRows rows={5} /> : teams.data?.length ? (
        <div className="grid gap-6 xl:grid-cols-2">{teams.data.map((team) => <TeamCard key={team.id} team={team} onDelete={setDeleting} />)}</div>
      ) : <EmptyState title="No teams yet" message="Create your first team to start routing tickets." />}
      <Modal title="Delete team" open={Boolean(deleting)} onClose={() => setDeleting(null)}>
        <p className="mb-5 text-sm text-slate-300">
          This team has {deleting?.activeTicketCount ?? 0} active tickets. Deleting it will not affect existing tickets but new tickets will not be routed here.
        </p>
        <div className="flex justify-end gap-2">
          <Button variant="secondary" onClick={() => setDeleting(null)}>Cancel</Button>
          <Button variant="danger" loading={mutations.delete.isPending} onClick={() => void confirmDelete()}>Delete</Button>
        </div>
      </Modal>
    </div>
  );
}
