import type { Message } from "@/types";
import { formatDate } from "@/utils/format";

export function MessageThread({ messages }: { messages: Message[] }) {
  return (
    <div className="space-y-4">
      {messages
        .slice()
        .sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime())
        .map((message) => (
          <div key={message.id} className={`flex ${message.role === "assistant" ? "justify-end" : "justify-start"}`}>
            <div className={`max-w-[80%] rounded-xl border p-4 ${message.role === "assistant" ? "border-blue-500/30 bg-blue-500/10" : "border-slate-700 bg-surface-2"}`}>
              <div className="mb-2 flex items-center justify-between gap-4 text-xs text-slate-400">
                <span className="font-semibold uppercase">{message.role}</span>
                <span>{formatDate(message.createdAt)}</span>
              </div>
              <p className="whitespace-pre-wrap text-sm leading-6 text-slate-100">{message.content}</p>
            </div>
          </div>
        ))}
    </div>
  );
}
