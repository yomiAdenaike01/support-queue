import type { Worker } from "@/types";

export const mockWorkers: Worker[] = [
  {
    name: "python-worker-1",
    language: "Python",
    status: "ONLINE",
    processedCount: 892,
    failedCount: 12,
    lastHeartbeat: new Date().toISOString(),
  },
  {
    name: "python-worker-2",
    language: "Python",
    status: "ONLINE",
    processedCount: 445,
    failedCount: 7,
    lastHeartbeat: new Date().toISOString(),
  },
  {
    name: "rust-worker-1",
    language: "Rust",
    status: "DEGRADED",
    processedCount: 215,
    failedCount: 16,
    lastHeartbeat: new Date(Date.now() - 45_000).toISOString(),
  },
];
