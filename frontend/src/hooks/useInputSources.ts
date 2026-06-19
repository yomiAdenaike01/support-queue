import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getInputSources, updateInputSource } from "@/api/inputSources";

export function useInputSources() {
  return useQuery({ queryKey: ["input-sources"], queryFn: getInputSources, staleTime: 20_000 });
}

export function useInputSourceMutations() {
  const queryClient = useQueryClient();
  return {
    update: useMutation({
      mutationFn: updateInputSource,
      onSuccess: () => {
        void queryClient.invalidateQueries({ queryKey: ["input-sources"] });
      },
    }),
  };
}
