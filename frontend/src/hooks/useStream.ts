import { useQuery } from "@tanstack/react-query";
import { getStreamStats } from "@/api/stream";

export function useStream() {
  return useQuery({ queryKey: ["stream"], queryFn: getStreamStats, refetchInterval: 30_000, staleTime: 10_000 });
}
