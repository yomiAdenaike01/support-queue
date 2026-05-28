import { useQuery } from "@tanstack/react-query";
import { getMetrics } from "@/api/metrics";

const refreshInterval = Number(import.meta.env.VITE_REFRESH_INTERVAL ?? 30000);

export function useMetrics() {
  return useQuery({ queryKey: ["metrics"], queryFn: getMetrics, refetchInterval: refreshInterval, staleTime: 10_000 });
}
