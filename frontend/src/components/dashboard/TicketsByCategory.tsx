import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import { Card } from "@/components/ui/Card";
import type { Ticket } from "@/types";

const colors = ["#06b6d4", "#3b82f6", "#10b981", "#ef4444", "#8b5cf6", "#64748b", "#f59e0b"];

export function TicketsByCategory({ tickets }: { tickets: Ticket[] }) {
  const data = Object.entries(
    tickets.reduce<Record<string, number>>((acc, ticket) => {
      const key = ticket.category ?? "UNCLASSIFIED";
      acc[key] = (acc[key] ?? 0) + 1;
      return acc;
    }, {}),
  ).map(([name, value]) => ({ name, value }));
  return (
    <Card>
      <h2 className="mb-4 text-lg font-semibold">By Category</h2>
      <div className="h-64">
        <ResponsiveContainer>
          <PieChart>
            <Pie data={data} dataKey="value" nameKey="name" innerRadius={48} outerRadius={84}>
              {data.map((entry, index) => <Cell key={entry.name} fill={colors[index % colors.length]} />)}
            </Pie>
            <Tooltip contentStyle={{ background: "#111827", border: "1px solid #1f2937" }} />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
}
