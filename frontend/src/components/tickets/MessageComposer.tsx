import { Send } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/Button";
import type { Message } from "@/types";

interface MessageComposerProps {
  disabled?: boolean;
  loading?: boolean;
  onSubmit: (message: Pick<Message, "content" | "role">) => Promise<void>;
}

export function MessageComposer({ disabled = false, loading = false, onSubmit }: MessageComposerProps) {
  const [content, setContent] = useState("");
  const [role, setRole] = useState<Message["role"]>("assistant");
  const isValid = content.trim().length >= 2;

  const submit = async () => {
    if (!isValid) return;
    await onSubmit({ content: content.trim(), role });
    setContent("");
    setRole("assistant");
  };

  return (
    <div className="mt-6 rounded-xl border border-slate-700 bg-surface-2 p-4">
      <div className="mb-3 flex flex-wrap items-center justify-between gap-3">
        <h3 className="text-sm font-semibold">Add message</h3>
        <select
          value={role}
          onChange={(event) => setRole(event.target.value as Message["role"])}
          className="rounded-lg border border-slate-700 bg-surface px-3 py-2 text-sm outline-none focus:border-accent"
          disabled={disabled || loading}
        >
          <option value="assistant">Assistant</option>
          <option value="customer">Customer</option>
        </select>
      </div>
      <textarea
        value={content}
        onChange={(event) => setContent(event.target.value)}
        rows={4}
        placeholder="Write the next message in this ticket thread"
        className="w-full resize-y rounded-lg border border-slate-700 bg-surface px-3 py-2 text-sm leading-6 outline-none focus:border-accent"
        disabled={disabled || loading}
      />
      <div className="mt-3 flex justify-end">
        <Button onClick={() => void submit()} disabled={!isValid || disabled} loading={loading}>
          <Send className="h-4 w-4" /> Add Message
        </Button>
      </div>
    </div>
  );
}
