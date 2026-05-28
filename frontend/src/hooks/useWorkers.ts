import { useQuery } from "@tanstack/react-query";
import { getWorkers } from "@/api/workers";

export function useWorkers() {
  return useQuery({ queryKey: ["workers"], queryFn: getWorkers, refetchInterval: 15_000, staleTime: 10_000 });
}
