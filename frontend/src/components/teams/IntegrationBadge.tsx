import { Badge } from "@/components/ui/Badge";

export function SlackIcon({ className = "h-4 w-4" }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" className={className} aria-hidden="true">
      <path fill="#36C5F0" d="M8.4 2a2 2 0 0 0 0 4h2V4a2 2 0 0 0-2-2Z" />
      <path fill="#2EB67D" d="M13.6 2a2 2 0 0 0-2 2v5.2a2 2 0 1 0 4 0V4a2 2 0 0 0-2-2Z" />
      <path fill="#ECB22E" d="M22 8.4a2 2 0 0 0-4 0v2h2a2 2 0 0 0 2-2Z" />
      <path fill="#E01E5A" d="M14.8 11.6a2 2 0 1 0 0 4H20a2 2 0 1 0 0-4h-5.2Z" />
      <path fill="#36C5F0" d="M15.6 20a2 2 0 0 0-4 0 2 2 0 1 0 4 0Z" />
      <path fill="#2EB67D" d="M10.4 14.8a2 2 0 1 0-4 0V20a2 2 0 1 0 4 0v-5.2Z" />
      <path fill="#ECB22E" d="M2 15.6a2 2 0 0 0 4 0v-2H4a2 2 0 0 0-2 2Z" />
      <path fill="#E01E5A" d="M4 10.4h5.2a2 2 0 0 0 0-4H4a2 2 0 0 0 0 4Z" />
    </svg>
  );
}

export function TeamsIcon({ className = "h-4 w-4" }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" className={className} aria-hidden="true">
      <rect x="8" y="5" width="11" height="13" rx="3" fill="#6264A7" />
      <rect x="3" y="8" width="10" height="10" rx="2" fill="#464775" />
      <circle cx="17.5" cy="4.5" r="2.5" fill="#7B83EB" />
      <circle cx="20.5" cy="8.5" r="2" fill="#7B83EB" />
      <path fill="#fff" d="M5.5 10h6v1.4H9.2V16H7.8v-4.6H5.5V10Z" />
    </svg>
  );
}

export function IntegrationBadge({ provider, connected }: { provider: "slack" | "teams"; connected: boolean }) {
  const Icon = provider === "slack" ? SlackIcon : TeamsIcon;
  return (
    <Badge className={connected ? "border-emerald-500/30 bg-emerald-500/15 text-emerald-300" : "border-slate-600 bg-slate-800 text-slate-400"}>
      <Icon /> {provider === "slack" ? "Slack" : "Teams"} {connected ? "Connected" : "Not configured"}
    </Badge>
  );
}
