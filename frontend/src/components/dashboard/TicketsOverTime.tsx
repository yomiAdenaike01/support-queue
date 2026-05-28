import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis, CartesianGrid } from "recharts";
import { Card } from "@/components/ui/Card";

const data = Array.from({ length: 24 }, (_, hour) => ({ hour: `${hour}:00`, tickets: 8 + ((hour * 7) % 23) }));

export function TicketsOverTime() {
  return (
    <Card>
      <h2 className="mb-4 text-lg font-semibold">Tickets Over Time</h2>
      <div className="h-72">
        <ResponsiveContainer>
          <LineChart data={data}>
            <CartesianGrid stroke="#1f2937" />
            <XAxis dataKey="hour" stroke="#94a3b8" />
            <YAxis stroke="#94a3b8" />
            <Tooltip contentStyle={{ background: "#111827", border: "1px solid #1f2937" }} />
            <Line type="monotone" dataKey="tickets" stroke="#10b981" strokeWidth={3} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
}
