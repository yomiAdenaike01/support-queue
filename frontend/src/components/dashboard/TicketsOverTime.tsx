import {
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  CartesianGrid,
} from "recharts";
import { Card } from "@/components/ui/Card";
import { useQuery } from "@tanstack/react-query";
import { getTicketsOvertime } from "@/api/metrics";

export function TicketsOverTime() {
  const { data: response } = useQuery({
    queryKey: ["metrics-tickets-over-time"],
    queryFn: () => getTicketsOvertime(),
  });

  console.log({ response });

  if (!response) return null;

  return (
    <Card>
      <h2 className="mb-4 text-lg font-semibold">Tickets Over Time</h2>
      <div className="h-72">
        <ResponsiveContainer>
          <LineChart data={response}>
            <CartesianGrid stroke="#1f2937" />
            <XAxis dataKey="hour" stroke="#94a3b8" />
            <YAxis stroke="#94a3b8" />
            <Tooltip
              contentStyle={{
                background: "#111827",
                border: "1px solid #1f2937",
              }}
            />
            <Line
              type="monotone"
              dataKey="tickets"
              stroke="#10b981"
              strokeWidth={3}
              dot={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
}
