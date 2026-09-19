import client from "./client";
import type { ApiEnvelope, Alert } from "@/types";

interface AlertsResponse {
  machine_id: string;
  alerts: Alert[];
}

export async function listAlerts(machineId: string, isResolved?: boolean) {
  const res = await client.get<ApiEnvelope<AlertsResponse>>("/alerts", {
    params: { machine_id: machineId, is_resolved: isResolved },
  });
  return res.data.data!.alerts;
}