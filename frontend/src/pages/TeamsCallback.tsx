import { useEffect } from "react";
import { CallbackShell } from "@/pages/SlackCallback";

export function TeamsCallback() {
  useEffect(() => {
    window.opener?.postMessage({ type: "OAUTH_SUCCESS" }, "*");
    window.setTimeout(() => window.close(), 600);
  }, []);
  return <CallbackShell text="Connecting to Microsoft Teams..." />;
}
