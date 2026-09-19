import client from "./client";
import type { ApiEnvelope, AuthResponse } from "@/types";

export async function login(username: string, password: string) {
  const res = await client.post<ApiEnvelope<AuthResponse>>("/auth/login", {
    username,
    password,
  });
  return res.data.data!;
}