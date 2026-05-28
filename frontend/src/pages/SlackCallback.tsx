import { useEffect } from "react";

export function SlackCallback() {
  useEffect(() => {
    window.opener?.postMessage({ type: "OAUTH_SUCCESS" }, "*");
    window.setTimeout(() => window.close(), 600);
  }, []);
  return <CallbackShell text="Connecting to Slack..." />;
}

export function CallbackShell({ text }: { text: string }) {
  return (
    <div className="grid min-h-screen place-items-center bg-[--color-bg] text-white">
      <div className="text-center">
        <div className="mx-auto mb-4 h-10 w-10 animate-spin rounded-full border-2 border-slate-700 border-t-blue-400" />
        <div>{text}</div>
      </div>
    </div>
  );
}
