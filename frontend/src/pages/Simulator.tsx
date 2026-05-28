import { useState } from "react";
import { SimulatorControls, type SimulatorLogItem } from "@/components/simulator/SimulatorControls";
import { SimulatorLog } from "@/components/simulator/SimulatorLog";
import { Card } from "@/components/ui/Card";

export function Simulator() {
  const [items, setItems] = useState<SimulatorLogItem[]>([]);
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold">Simulator</h1>
        <p className="text-sm text-slate-400">Generate ticket load and failure cases for the worker pipeline.</p>
      </div>
      <SimulatorControls onLog={(item) => setItems((current) => [...current, item])} />
      <div className="grid gap-6 xl:grid-cols-[2fr_1fr]">
        <SimulatorLog items={items} onClear={() => setItems([])} />
        <Card>
          <h2 className="mb-3 text-lg font-semibold">Instructions</h2>
          <p className="text-sm leading-6 text-slate-300">
            Seed creates a batch of tickets, auto generation starts continuous traffic, and quick actions create one-off tickets.
            Workers process tickets automatically through the configured stream.
          </p>
        </Card>
      </div>
    </div>
  );
}
