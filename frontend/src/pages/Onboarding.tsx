import { useState } from "react";
import type { ReactNode } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { SlackIcon, TeamsIcon } from "@/components/teams/IntegrationBadge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { openOAuthPopup } from "@/utils/oauth";

export function Onboarding() {
  const { inviteToken = "" } = useParams();
  const navigate = useNavigate();
  const [step, setStep] = useState(0);
  const [slack, setSlack] = useState(false);
  const [teams, setTeams] = useState(false);

  const connect = (provider: "slack" | "teams") => {
    openOAuthPopup(`/api/onboarding/${inviteToken}/${provider}`, () => {
      if (provider === "slack") setSlack(true);
      if (provider === "teams") setTeams(true);
    }, () => undefined);
  };

  return (
    <div className="grid min-h-screen place-items-center bg-[--color-bg] p-4">
      <Card className="w-full max-w-xl text-center">
        <div className="mx-auto mb-6 grid h-12 w-12 place-items-center rounded-lg bg-accent font-mono font-bold">SO</div>
        {step === 0 ? (
          <>
            <h1 className="text-2xl font-semibold">You've been invited to join Billing Team at Acme Corp</h1>
            <p className="mt-3 text-slate-400">as Senior Support Agent by Sarah Mitchell</p>
            <Button className="mt-6" onClick={() => setStep(1)}>Get Started</Button>
          </>
        ) : null}
        {step === 1 ? (
          <Step title="Connect your Slack account" step="Step 1 of 2" connected={slack} connectedText="Connected as @sarah.mitchell" onConnect={() => connect("slack")} onNext={() => setStep(2)} button={<><SlackIcon /> Connect with Slack</>} />
        ) : null}
        {step === 2 ? (
          <Step title="Connect your Microsoft Teams account" step="Step 2 of 2" connected={teams} connectedText="Connected as sarah.mitchell@acme.com" onConnect={() => connect("teams")} onNext={() => setStep(3)} button={<><TeamsIcon /> Connect with Microsoft Teams</>} />
        ) : null}
        {step === 3 ? (
          <>
            <h1 className="text-2xl font-semibold">You're all set!</h1>
            <p className="mt-3 text-slate-400">You've joined Billing Team as Senior Support Agent.</p>
            {!slack || !teams ? <div className="mt-4 rounded-lg border border-amber-500/30 bg-amber-500/10 p-3 text-sm text-amber-100">You skipped some integrations. You can connect them later from profile settings.</div> : null}
            <Button className="mt-6" onClick={() => navigate("/")}>Open Dashboard</Button>
          </>
        ) : null}
      </Card>
    </div>
  );
}

function Step({ title, step, connected, connectedText, onConnect, onNext, button }: { title: string; step: string; connected: boolean; connectedText: string; onConnect: () => void; onNext: () => void; button: ReactNode }) {
  return (
    <>
      <div className="text-sm text-slate-400">{step}</div>
      <h1 className="mt-2 text-2xl font-semibold">{title}</h1>
      {connected ? <div className="mt-5 rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-3 text-emerald-200">{connectedText}</div> : <Button className="mt-6" variant={title.includes("Slack") ? "brand-slack" : "brand-teams"} onClick={onConnect}>{button}</Button>}
      <div className="mt-6 flex justify-center gap-2">
        <Button variant="ghost" onClick={onNext}>Skip for now</Button>
        <Button onClick={onNext}>Continue</Button>
      </div>
    </>
  );
}
