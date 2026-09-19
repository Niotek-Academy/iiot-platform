import client from "./client";
import type { ApiEnvelope, Sensor } from "@/types";

export async function listSensors(machineId?: string) {
  const res = await client.get<ApiEnvelope<Sensor[]>>("/sensors", {
    params: machineId ? { machine_id: machineId } : undefined,
  });
  return res.data.data!;
}

export async function createSensor(input: {
  sensor_id: string;
  machine_id: string;
  metric_name: string;
  unit: string;
  source_address?: string;
}) {
  const res = await client.post<ApiEnvelope<Sensor>>("/sensors", input);
  return res.data.data!;
}

export async function deleteSensor(sensorId: string) {
  await client.delete(`/sensors/${sensorId}`);
}