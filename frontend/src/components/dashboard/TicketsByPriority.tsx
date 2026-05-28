import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Card } from "@/components/ui/Card";
import type { Ticket } from "@/types";

export function TicketsByPriority({ tickets }: { tickets: Ticket[] }) {
  const data = ["URGENT", "HIGH", "MEDIUM", "LOW"].map((priority) => ({
    priority,
    count: tickets.filter((ticket) => ticket.priority === priority).length,
  }));
  return (
    <Card>
      <h2 className="mb-4 text-lg font-semibold">By Priority</h2>
      <div className="h-64">
        <ResponsiveContainer>
          <BarChart data={data}>
            <CartesianGrid stroke="#1f2937" />
            <XAxis dataKey="priority" stroke="#94a3b8" />
            <YAxis stroke="#94a3b8" />
            <Tooltip contentStyle={{ background: "#111827", border: "1px solid #1f2937" }} />
            <Bar dataKey="count" fill="#3b82f6" radius={[6, 6, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
}
