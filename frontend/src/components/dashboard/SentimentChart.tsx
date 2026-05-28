import { Bar, BarChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Card } from "@/components/ui/Card";
import type { Ticket } from "@/types";

export function SentimentChart({ tickets }: { tickets: Ticket[] }) {
  const data = ["HIGH", "MEDIUM", "LOW"].map((sentiment) => ({
    sentiment,
    count: tickets.filter((ticket) => ticket.sentimentLabel === sentiment).length,
  }));
  return (
    <Card>
      <h2 className="mb-4 text-lg font-semibold">Sentiment</h2>
      <div className="h-64">
        <ResponsiveContainer>
          <BarChart data={data}>
            <XAxis dataKey="sentiment" stroke="#94a3b8" />
            <YAxis stroke="#94a3b8" />
            <Tooltip contentStyle={{ background: "#111827", border: "1px solid #1f2937" }} />
            <Bar dataKey="count" fill="#f59e0b" radius={[6, 6, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
}
