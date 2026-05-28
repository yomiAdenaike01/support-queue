export function openOAuthPopup(url: string, onSuccess: () => void, onError: (err: string) => void): void {
  const width = 600;
  const height = 700;
  const left = window.screenX + (window.outerWidth - width) / 2;
  const top = window.screenY + (window.outerHeight - height) / 2;
  const popup = window.open(url, "oauth", `width=${width},height=${height},left=${left},top=${top}`);

  const handler = (event: MessageEvent<{ type?: string; error?: string }>) => {
    if (event.data?.type === "OAUTH_SUCCESS") {
      window.removeEventListener("message", handler);
      popup?.close();
      onSuccess();
    }
    if (event.data?.type === "OAUTH_ERROR") {
      window.removeEventListener("message", handler);
      popup?.close();
      onError(event.data.error ?? "OAuth failed");
    }
  };

  window.addEventListener("message", handler);

  const interval = window.setInterval(() => {
    if (popup?.closed) {
      window.clearInterval(interval);
      window.removeEventListener("message", handler);
    }
  }, 500);
}
