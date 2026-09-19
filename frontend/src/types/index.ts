export interface ApiEnvelope<T> {
  success: boolean;
  data?: T;
  error?: { code: string; message: string };
}

export type Role = "ADMIN" | "OPERATOR";

export interface User {
  user_id: number;
  username: string;
  role: Role;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export type MachineStatus = "OPERATIONAL" | "WARNING" | "CRITICAL" | "STOPPED";

export interface Machine {
  machine_id: string;
  name: string;
  location: string;
  status: MachineStatus;
  control_address?: string;
  created_at: string;
}

export interface MachineOverview {
  machine_id: string;
  name: string;
  location: string;
  status: MachineStatus;
  current_health_score: number | null;
  rul_hours: number | null;
  active_alerts_count: number;
  updated_at: string;
}

export interface Sensor {
  sensor_id: string;
  machine_id: string;
  metric_name: string;
  unit: string;
  source_address?: string;
  created_at: string;
}

export interface TelemetryDataPoint {
  value: number;
  recorded_at: string;
}

export interface TelemetryResponse {
  sensor_id: string;
  unit: string;
  count: number;
  data_points: TelemetryDataPoint[];
}

export type AlertSeverity = "INFO" | "WARNING" | "CRITICAL";

export interface Alert {
  alert_id: number;
  severity: AlertSeverity;
  message: string;
  is_resolved: boolean;
  created_at: string;
}

export type CommandType = "EMERGENCY_STOP" | "START" | "RESET_ALERTS";

export interface CommandLog {
  command_id: number;
  command_type: CommandType;
  issued_by: string;
  reason: string;
  executed_at: string;
}

export interface WsAiHealth {
  health_score: number;
  is_anomaly: boolean;
  anomaly_features: string[];
  rul_hours: number;
}

export interface WsPayload {
  event_type: "TELEMETRY_AND_AI_UPDATE";
  timestamp: string;
  machine_id: string;
  telemetry: Record<string, number>;
  ai_health?: WsAiHealth;
}