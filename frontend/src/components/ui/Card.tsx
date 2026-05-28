import type { PropsWithChildren } from "react";

interface CardProps extends PropsWithChildren {
  className?: string;
}

export function Card({ children, className = "" }: CardProps) {
  return <section className={`rounded-xl border border-[--color-border] bg-[--color-surface] p-6 ${className}`}>{children}</section>;
}
