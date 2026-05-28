import { create } from "zustand";
import type { TicketCategory, TicketPriority, TicketStatus } from "@/types";

interface UiState {
  sidebarOpen: boolean;
  search: string;
  status: TicketStatus | "ALL";
  priority: TicketPriority | "ALL";
  category: TicketCategory | "ALL";
  setSidebarOpen: (open: boolean) => void;
  setSearch: (search: string) => void;
  setStatus: (status: TicketStatus | "ALL") => void;
  setPriority: (priority: TicketPriority | "ALL") => void;
  setCategory: (category: TicketCategory | "ALL") => void;
  clearFilters: () => void;
}

export const useUiStore = create<UiState>((set) => ({
  sidebarOpen: true,
  search: "",
  status: "ALL",
  priority: "ALL",
  category: "ALL",
  setSidebarOpen: (sidebarOpen) => set({ sidebarOpen }),
  setSearch: (search) => set({ search }),
  setStatus: (status) => set({ status }),
  setPriority: (priority) => set({ priority }),
  setCategory: (category) => set({ category }),
  clearFilters: () => set({ search: "", status: "ALL", priority: "ALL", category: "ALL" }),
}));
