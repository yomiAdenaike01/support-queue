import {
  AtSign,
  Code2,
  Facebook,
  FileText,
  Globe2,
  MessageCircle,
  MessagesSquare,
  Phone,
  Smartphone,
  Users,
} from "lucide-react";
import { useMemo, useState } from "react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Modal } from "@/components/ui/Modal";
import { SkeletonRows } from "@/components/ui/Skeleton";
import { useToast } from "@/components/ui/Toast";
import { useInputSourceMutations, useInputSources } from "@/hooks/useInputSources";
import type { TicketInputSourceConfig, TicketInputSourceStatus, TicketInputSourceType } from "@/types";
import { relativeTime } from "@/utils/format";

const sourceIcons: Record<TicketInputSourceType, typeof AtSign> = {
  EMAIL: AtSign,
  WEB_FORM: FileText,
  LIVE_CHAT: MessagesSquare,
  PHONE: Phone,
  SOCIAL_MEDIA: Facebook,
  WHATSAPP_BUSINESS: MessageCircle,
  SMS: Smartphone,
  IN_APP: Globe2,
  COMMUNITY_FORUM: Users,
  API: Code2,
};

const statusClasses: Record<TicketInputSourceStatus, string> = {
  enabled: "border-emerald-500/30 bg-emerald-500/15 text-emerald-300",
  disabled: "border-slate-600 bg-slate-800 text-slate-300",
  needs_setup: "border-amber-500/30 bg-amber-500/15 text-amber-300",
};

export function InputSources() {
  const sources = useInputSources();
  const mutations = useInputSourceMutations();
  const toast = useToast();
  const [editing, setEditing] = useState<TicketInputSourceConfig | null>(null);
  const [draft, setDraft] = useState<TicketInputSourceConfig | null>(null);
  const enabledCount = useMemo(() => (sources.data ?? []).filter((source) => source.isEnabled).length, [sources.data]);
  const todaysTotal = useMemo(() => (sources.data ?? []).reduce((sum, source) => sum + source.createdToday, 0), [sources.data]);

  const openEditor = (source: TicketInputSourceConfig) => {
    setEditing(source);
    setDraft({ ...source, settings: { ...source.settings } });
  };

  const saveDraft = async () => {
    if (!draft) return;
    const status: TicketInputSourceStatus = draft.isEnabled ? (draft.connectionValue ? "enabled" : "needs_setup") : "disabled";
    await mutations.update.mutateAsync({ ...draft, status });
    toast.push(`${draft.name} source updated`, "success");
    setEditing(null);
    setDraft(null);
  };

  const toggleSource = async (source: TicketInputSourceConfig) => {
    const isEnabled = !source.isEnabled;
    await mutations.update.mutateAsync({
      ...source,
      isEnabled,
      status: isEnabled ? (source.connectionValue ? "enabled" : "needs_setup") : "disabled",
    });
    toast.push(`${source.name} ${isEnabled ? "enabled" : "disabled"}`, "success");
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Input Sources</h1>
          <p className="mt-1 text-sm text-slate-400">Configure where tickets can be created from across customer channels.</p>
        </div>
        <div className="grid grid-cols-2 gap-3 text-sm">
          <Summary label="Enabled" value={`${enabledCount}/${sources.data?.length ?? 0}`} />
          <Summary label="Created today" value={String(todaysTotal)} />
        </div>
      </div>

      {sources.isLoading ? <SkeletonRows rows={8} /> : null}
      {sources.isError ? <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-3 text-red-200">Unable to load input sources.</div> : null}

      <div className="grid gap-4 lg:grid-cols-2 2xl:grid-cols-3">
        {(sources.data ?? []).map((source) => {
          const Icon = sourceIcons[source.type];
          return (
            <Card key={source.id} className="flex flex-col gap-5">
              <div className="flex items-start justify-between gap-4">
                <div className="flex gap-3">
                  <div className="grid h-11 w-11 shrink-0 place-items-center rounded-lg bg-surface-2 text-blue-300">
                    <Icon className="h-5 w-5" />
                  </div>
                  <div>
                    <h2 className="font-semibold">{source.name}</h2>
                    <p className="mt-1 text-sm leading-6 text-slate-400">{source.description}</p>
                  </div>
                </div>
                <span className={`shrink-0 rounded-full border px-2.5 py-1 text-xs font-semibold ${statusClasses[source.status]}`}>
                  {source.status.replace("_", " ")}
                </span>
              </div>

              <div className="grid gap-3 text-sm sm:grid-cols-2">
                <SourceStat label={source.connectionLabel} value={source.connectionValue || "Not configured"} />
                <SourceStat label="Routing team" value={source.routingTeam ?? "Default routing"} />
                <SourceStat label="Created today" value={String(source.createdToday)} />
                <SourceStat label="Last received" value={source.lastReceivedAt ? relativeTime(source.lastReceivedAt) : "Never"} />
              </div>

              <div className="mt-auto flex flex-wrap items-center justify-between gap-3 border-t border-slate-800 pt-4">
                <label className="inline-flex items-center gap-2 text-sm text-slate-300">
                  <input
                    type="checkbox"
                    checked={source.isEnabled}
                    onChange={() => void toggleSource(source)}
                    className="h-4 w-4 rounded border-slate-600 bg-surface-2"
                  />
                  Accept tickets
                </label>
                <Button variant="secondary" onClick={() => openEditor(source)}>Configure</Button>
              </div>
            </Card>
          );
        })}
      </div>

      <Modal title={editing ? `Configure ${editing.name}` : "Configure input source"} open={Boolean(editing && draft)} onClose={() => setEditing(null)}>
        {draft ? (
          <div className="space-y-4">
            <label className="flex items-center justify-between rounded-lg bg-surface-2 p-3 text-sm">
              <span>Accept tickets from this source</span>
              <input type="checkbox" checked={draft.isEnabled} onChange={(event) => setDraft({ ...draft, isEnabled: event.target.checked })} />
            </label>
            <Field label={draft.connectionLabel}>
              <input
                value={draft.connectionValue}
                onChange={(event) => setDraft({ ...draft, connectionValue: event.target.value })}
                placeholder={connectionPlaceholder(draft.type)}
                className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2 text-sm outline-none focus:border-accent"
              />
            </Field>
            <Field label="Routing team">
              <input
                value={draft.routingTeam ?? ""}
                onChange={(event) => setDraft({ ...draft, routingTeam: event.target.value || null })}
                placeholder="Default routing"
                className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2 text-sm outline-none focus:border-accent"
              />
            </Field>
            <div className="grid gap-3 sm:grid-cols-2">
              <label className="flex items-center gap-2 rounded-lg bg-surface-2 p-3 text-sm">
                <input type="checkbox" checked={draft.autoCreateTickets} onChange={(event) => setDraft({ ...draft, autoCreateTickets: event.target.checked })} />
                Auto-create tickets
              </label>
              <label className="flex items-center gap-2 rounded-lg bg-surface-2 p-3 text-sm">
                <input type="checkbox" checked={draft.requiresAuth} onChange={(event) => setDraft({ ...draft, requiresAuth: event.target.checked })} />
                Require authentication
              </label>
            </div>
            <SourceSpecificFields draft={draft} onChange={setDraft} />
            <div className="flex justify-end gap-2 pt-2">
              <Button variant="secondary" onClick={() => setEditing(null)}>Cancel</Button>
              <Button loading={mutations.update.isPending} onClick={() => void saveDraft()}>Save Source</Button>
            </div>
          </div>
        ) : null}
      </Modal>
    </div>
  );
}

function Summary({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg border border-slate-800 bg-surface px-4 py-3">
      <div className="text-xs text-slate-400">{label}</div>
      <div className="mt-1 font-mono text-xl font-semibold">{value}</div>
    </div>
  );
}

function SourceStat({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg bg-surface-2 p-3">
      <div className="text-xs text-slate-500">{label}</div>
      <div className="mt-1 truncate text-slate-200" title={value}>{value}</div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="block text-sm font-medium">
      {label}
      <div className="mt-2">{children}</div>
    </label>
  );
}

function SourceSpecificFields({ draft, onChange }: { draft: TicketInputSourceConfig; onChange: (draft: TicketInputSourceConfig) => void }) {
  const entries = sourceFields(draft.type);
  if (entries.length === 0) return null;
  return (
    <div className="grid gap-3 sm:grid-cols-2">
      {entries.map((field) => (
        <Field key={field.key} label={field.label}>
          <input
            value={String(draft.settings[field.key] ?? "")}
            onChange={(event) => onChange({ ...draft, settings: { ...draft.settings, [field.key]: event.target.value } })}
            placeholder={field.placeholder}
            className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2 text-sm outline-none focus:border-accent"
          />
        </Field>
      ))}
    </div>
  );
}

function sourceFields(type: TicketInputSourceType): Array<{ key: string; label: string; placeholder: string }> {
  switch (type) {
    case "EMAIL":
      return [{ key: "inboundAddress", label: "Inbound address", placeholder: "support@example.com" }];
    case "WEB_FORM":
      return [{ key: "allowedDomain", label: "Allowed domain", placeholder: "www.example.com" }];
    case "LIVE_CHAT":
      return [{ key: "widgetOrigin", label: "Widget origin", placeholder: "https://www.example.com" }];
    case "PHONE":
      return [{ key: "defaultQueue", label: "Default queue", placeholder: "frontline-phone" }];
    case "SOCIAL_MEDIA":
      return [
        { key: "xHandle", label: "Twitter/X handle", placeholder: "@example" },
        { key: "facebookPageId", label: "Facebook page ID", placeholder: "123456789" },
      ];
    case "WHATSAPP_BUSINESS":
      return [
        { key: "phoneNumberId", label: "Phone number ID", placeholder: "100000000000000" },
        { key: "webhookVerifyToken", label: "Webhook verify token", placeholder: "••••••••" },
      ];
    case "SMS":
      return [
        { key: "provider", label: "Provider", placeholder: "Twilio" },
        { key: "phoneNumber", label: "Phone number", placeholder: "+15551234567" },
      ];
    case "IN_APP":
      return [{ key: "appKey", label: "App key", placeholder: "app_prod_01" }];
    case "COMMUNITY_FORUM":
      return [{ key: "forumUrl", label: "Forum URL", placeholder: "https://community.example.com" }];
    case "API":
      return [
        { key: "keyMode", label: "API key mode", placeholder: "Scoped keys" },
        { key: "rateLimit", label: "Rate limit", placeholder: "600/min" },
      ];
  }
}

function connectionPlaceholder(type: TicketInputSourceType): string {
  const placeholders: Record<TicketInputSourceType, string> = {
    EMAIL: "support@example.com",
    WEB_FORM: "www.example.com",
    LIVE_CHAT: "https://www.example.com",
    PHONE: "frontline-phone",
    SOCIAL_MEDIA: "@example or page ID",
    WHATSAPP_BUSINESS: "Phone number ID",
    SMS: "+15551234567",
    IN_APP: "app_prod_01",
    COMMUNITY_FORUM: "https://community.example.com",
    API: "Scoped keys",
  };
  return placeholders[type];
}
