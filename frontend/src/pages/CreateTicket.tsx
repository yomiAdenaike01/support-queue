import { useState } from "react";
import type { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { useToast } from "@/components/ui/Toast";
import { useTicketMutations } from "@/hooks/useTickets";

export function CreateTicket() {
  const [customerEmail, setCustomerEmail] = useState("");
  const [subject, setSubject] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();
  const toast = useToast();
  const { create } = useTicketMutations();
  const emailValid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(customerEmail);
  const valid = emailValid && subject.length >= 5 && subject.length <= 200 && message.length >= 10;

  const submit = async () => {
    try {
      const ticket = await create.mutateAsync({ customerEmail, subject, message });
      toast.push("Ticket created", "success");
      navigate(`/tickets/${ticket.id}`);
    } catch (error) {
      toast.push(error instanceof Error ? error.message : "Unable to create ticket", "error");
    }
  };

  return (
    <Card className="max-w-3xl">
      <h1 className="mb-6 text-2xl font-semibold">Create Ticket</h1>
      <div className="space-y-4">
        <Field label="Customer Email" error={customerEmail && !emailValid ? "Enter a valid email address" : ""}>
          <input value={customerEmail} onChange={(event) => setCustomerEmail(event.target.value)} className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2" />
        </Field>
        <Field label="Subject" error={subject && subject.length < 5 ? "Subject must be at least 5 characters" : ""}>
          <input value={subject} onChange={(event) => setSubject(event.target.value)} maxLength={200} className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2" />
        </Field>
        <Field label="Message" error={message && message.length < 10 ? "Message must be at least 10 characters" : ""}>
          <textarea value={message} onChange={(event) => setMessage(event.target.value)} rows={6} className="w-full rounded-lg border border-slate-700 bg-surface-2 px-3 py-2" />
        </Field>
        <Button disabled={!valid} loading={create.isPending} onClick={() => void submit()}>Submit Ticket</Button>
      </div>
    </Card>
  );
}

function Field({ label, error, children }: { label: string; error: string; children: ReactNode }) {
  return <label className="block text-sm font-medium">{label}<div className="mt-2">{children}</div>{error ? <div className="mt-1 text-xs text-red-300">{error}</div> : null}</label>;
}
