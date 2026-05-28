import { create } from "zustand";

interface ToastItem {
  id: number;
  message: string;
  tone: "success" | "error" | "info";
}

interface ToastState {
  items: ToastItem[];
  push: (message: string, tone?: ToastItem["tone"]) => void;
  remove: (id: number) => void;
}

export const useToast = create<ToastState>((set) => ({
  items: [],
  push: (message, tone = "info") => {
    const id = Date.now();
    set((state) => ({ items: [...state.items, { id, message, tone }] }));
    window.setTimeout(() => set((state) => ({ items: state.items.filter((item) => item.id !== id) })), 3500);
  },
  remove: (id) => set((state) => ({ items: state.items.filter((item) => item.id !== id) })),
}));

export function ToastViewport() {
  const { items, remove } = useToast();
  return (
    <div className="fixed right-4 top-4 z-[60] space-y-2">
      {items.map((item) => (
        <button
          key={item.id}
          onClick={() => remove(item.id)}
          className={`block w-80 rounded-lg border p-3 text-left text-sm shadow-xl ${
            item.tone === "error"
              ? "border-red-500/40 bg-red-950 text-red-100"
              : item.tone === "success"
                ? "border-emerald-500/40 bg-emerald-950 text-emerald-100"
                : "border-slate-700 bg-surface text-slate-100"
          }`}
        >
          {item.message}
        </button>
      ))}
    </div>
  );
}
