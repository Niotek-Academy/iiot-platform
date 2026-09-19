import client from "./client";
import type { ApiEnvelope, Role, User } from "@/types";

export async function listUsers() {
  const res = await client.get<ApiEnvelope<User[]>>("/users");
  return res.data.data!;
}

export async function updateUserRole(userId: number, role: Role) {
  const res = await client.patch<ApiEnvelope<User>>(`/users/${userId}/role`, { role });
  return res.data.data!;
}

export async function registerUser(username: string, password: string, role: Role) {
  const res = await client.post<ApiEnvelope<User>>("/auth/register", { username, password, role });
  return res.data.data!;
}