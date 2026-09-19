import client from "./client";
import type { ApiEnvelope, CommandLog, CommandType } from "@/types";

interface CommandResponse {
  command_id: number;
  machine_id: string;
  command: string;
  status: string;
  executed_at: string;
}

export async function executeCommand(machineId: string, command: CommandType, reason: string) {
  const res = await client.post<ApiEnvelope<CommandResponse>>(`/machines/${machineId}/command`, {
    command,
    reason,
  });
  return res.data.data!;
}

interface CommandLogsResponse {
  machine_id: string;
  commands: CommandLog[];
}

export async function getCommandHistory(machineId: string, limit = 20) {
  const res = await client.get<ApiEnvelope<CommandLogsResponse>>(`/machines/${machineId}/commands`, {
    params: { limit },
  });
  return res.data.data!.commands;
}