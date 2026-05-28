import type { ButtonHTMLAttributes, PropsWithChildren } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement>, PropsWithChildren {
  variant?: "primary" | "secondary" | "danger" | "ghost" | "brand-slack" | "brand-teams";
  loading?: boolean;
}

const variants = {
  primary: "bg-accent text-white hover:bg-blue-600",
  secondary: "bg-surface-2 text-slate-100 hover:bg-slate-700",
  danger: "bg-red-600 text-white hover:bg-red-500",
  ghost: "bg-transparent text-slate-300 hover:bg-surface-2",
  "brand-slack": "bg-[#4A154B] text-white hover:bg-[#5b1a5c]",
  "brand-teams": "bg-[#6264A7] text-white hover:bg-[#7173bd]",
};

export function Button({ children, className = "", variant = "primary", loading = false, disabled, ...props }: ButtonProps) {
  return (
    <button
      className={`inline-flex min-h-10 items-center justify-center gap-2 rounded-lg px-4 py-2 text-sm font-semibold transition disabled:cursor-not-allowed disabled:opacity-50 ${variants[variant]} ${className}`}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? <span className="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white" /> : null}
      {children}
    </button>
  );
}
