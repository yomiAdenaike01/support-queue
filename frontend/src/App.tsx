import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Layout } from "@/components/layout/Layout";
import { ToastViewport } from "@/components/ui/Toast";
import { CreateTeam } from "@/pages/CreateTeam";
import { CreateTicket } from "@/pages/CreateTicket";
import { Dashboard } from "@/pages/Dashboard";
import { Onboarding } from "@/pages/Onboarding";
import { Settings } from "@/pages/Settings";
import { Simulator } from "@/pages/Simulator";
import { SlackCallback } from "@/pages/SlackCallback";
import { TeamDetail } from "@/pages/TeamDetail";
import { Teams } from "@/pages/Teams";
import { TeamsCallback } from "@/pages/TeamsCallback";
import { TicketDetail } from "@/pages/TicketDetail";
import { Tickets } from "@/pages/Tickets";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/integrations/slack/callback" element={<SlackCallback />} />
          <Route path="/integrations/teams/callback" element={<TeamsCallback />} />
          <Route path="/onboarding/:inviteToken" element={<Onboarding />} />
          <Route element={<Layout />}>
            <Route path="/" element={<Dashboard />} />
            <Route path="/tickets" element={<Tickets />} />
            <Route path="/tickets/new" element={<CreateTicket />} />
            <Route path="/tickets/:id" element={<TicketDetail />} />
            <Route path="/teams" element={<Teams />} />
            <Route path="/teams/new" element={<CreateTeam />} />
            <Route path="/teams/:id" element={<TeamDetail />} />
            <Route path="/simulator" element={<Simulator />} />
            <Route path="/settings" element={<Settings />} />
          </Route>
        </Routes>
      </BrowserRouter>
      <ToastViewport />
    </QueryClientProvider>
  );
}
