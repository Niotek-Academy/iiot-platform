import client from "./client";
import type { ApiEnvelope, Machine, MachineOverview } from "@/types";

export async function listMachines() {
  const res = await client.get<ApiEnvelope<Machine[]>>("/machines");
  return res.data.data!;
}

export async function getMachineOverview(machineId: string) {
  const res = await client.get<ApiEnvelope<MachineOverview>>(`/machines/${machineId}`);
  return res.data.data!;
}

export async function createMachine(input: {
  machine_id: string;
  name: string;
  location: string;
  status?: string;
  control_address?: string;
}) {
  const res = await client.post<ApiEnvelope<Machine>>("/machines", input);
  return res.data.data!;
}

export async function updateMachine(
  machineId: string,
  input: Partial<{ name: string; location: string; status: string; control_address: string }>
) {
  const res = await client.patch<ApiEnvelope<Machine>>(`/machines/${machineId}`, input);
  return res.data.data!;
}

export async function deleteMachine(machineId: string) {
  await client.delete(`/machines/${machineId}`);
}