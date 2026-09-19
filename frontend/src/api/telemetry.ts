import client from "./client";
import type { ApiEnvelope, TelemetryResponse } from "@/types";

export async function getTelemetryHistory(sensorId: string, limit = 50) {
  const res = await client.get<ApiEnvelope<TelemetryResponse>>("/telemetry", {
    params: { sensor_id: sensorId, limit },
  });
  return res.data.data!;
}