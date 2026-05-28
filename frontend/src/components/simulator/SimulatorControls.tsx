import { useState } from "react";
import { Play, Square, Wand2 } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { forceFailure, randomTicket, seedTickets, startSimulator, stopSimulator } from "@/api/simulator";

export interface SimulatorLogItem {
  id: number;
  type: string;
  ticketId?: string;
  timestamp: string;
}

export function SimulatorControls({ onLog }: { onLog: (item: SimulatorLogItem) => void }) {
  const [count, setCount] = useState(10);
  const [running, setRunning] = useState(false);
  const [loading, setLoading] = useState<string | null>(null);

  const run = async (type: string, action: () => Promise<{ count?: number; started?: boolean; stopped?: boolean; id?: string }>) => {
    setLoading(type);
    const result = await action();
    onLog({ id: Date.now(), type, ticketId: "id" in result ? result.id : undefined, timestamp: new Date().toISOString() });
    setLoading(null);
  };

  return (
    <Card>
      <div className="grid gap-6 lg:grid-cols-3">
        <div>
          <h2 className="mb-3 font-semibold">Seed Tickets</h2>
          <input
            type="number"
            min={1}
            max={100}
            value={count}
            onChange={(event) => setCount(Number(event.target.value))}
            className="mb-3 w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2"
          />
          <Button loading={loading === "seed"} onClick={() => run("Seeded tickets", () => seedTickets(count))}>
            <Wand2 className="h-4 w-4" /> Generate
          </Button>
        </div>
        <div>
          <h2 className="mb-3 font-semibold">Auto Generator</h2>
          <div className={`mb-3 inline-flex rounded-full px-3 py-1 text-sm ${running ? "bg-emerald-500/15 text-emerald-300" : "bg-slate-700 text-slate-300"}`}>
            {running ? "RUNNING" : "STOPPED"}
          </div>
          <div className="flex gap-2">
            <Button onClick={() => { setRunning(true); void run("Started simulator", startSimulator); }}><Play className="h-4 w-4" /> Start</Button>
            <Button variant="secondary" onClick={() => { setRunning(false); void run("Stopped simulator", stopSimulator); }}><Square className="h-4 w-4" /> Stop</Button>
          </div>
        </div>
        <div>
          <h2 className="mb-3 font-semibold">Quick Actions</h2>
          <div className="flex flex-wrap gap-2">
            <Button variant="secondary" onClick={() => run("Generated random ticket", randomTicket)}>Random Ticket</Button>
            <Button variant="danger" onClick={() => run("Forced failure", forceFailure)}>Force Failure</Button>
          </div>
        </div>
      </div>
    </Card>
  );
}
