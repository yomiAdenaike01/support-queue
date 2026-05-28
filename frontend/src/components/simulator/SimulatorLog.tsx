import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import type { SimulatorLogItem } from "@/components/simulator/SimulatorControls";
import { formatDate } from "@/utils/format";

export function SimulatorLog({ items, onClear }: { items: SimulatorLogItem[]; onClear: () => void }) {
  return (
    <Card>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold">Simulator Log</h2>
        <Button variant="secondary" onClick={onClear}>Clear Log</Button>
      </div>
      <div className="max-h-80 space-y-2 overflow-y-auto font-mono text-sm">
        {items.length === 0 ? <div className="text-slate-500">No simulator events yet.</div> : null}
        {items.map((item) => (
          <div key={item.id} className="rounded-lg bg-surface-2 px-3 py-2 text-slate-300">
            <span className="text-slate-500">{formatDate(item.timestamp)}</span> · {item.type}
            {item.ticketId ? <span className="text-blue-300"> · {item.ticketId}</span> : null}
          </div>
        ))}
      </div>
    </Card>
  );
}
