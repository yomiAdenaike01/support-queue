import axios from "axios";

export const BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
export const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== "false";

export const api = axios.create({
  baseURL: BASE_URL,
  headers: { "Content-Type": "application/json" },
});

export async function mockDelay<T>(value: T, ms = 180): Promise<T> {
  await new Promise((resolve) => window.setTimeout(resolve, ms));
  return value;
}
