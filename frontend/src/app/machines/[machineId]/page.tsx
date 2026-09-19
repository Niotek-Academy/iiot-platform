// Machine Detail with live WebSocket chart
"use client";

import { useEffect, useState, use as usePromise } from "react";
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";
import { Wifi, WifiOff, OctagonAlert } from "lucide-react";
import AppShell from "@/components/AppShell";
import HealthGauge from "@/components/HealthGauge";
import StatusBadge from "@/components/StatusBadge";
import { useAuth } from "@/context/AuthContext";
import { useMachineStream } from "@/hooks/useWebSocket";
import { getMachineOverview } from "@/api/machines";
import { listAlerts } from "@/api/alerts";
import { executeCommand, getCommandHistory } from "@/api/commands";
import type { Alert, CommandLog, CommandType, MachineOverview } from "@/types";

const MAX_LIVE_POINTS = 30;

export default function MachineDetailPage({
  params,
}: {
  params: Promise<{ machineId: string }>;
}) {
  const { machineId } = usePromise(params);
  const { token } = useAuth();
  const { latest, connected } = useMachineStream(machineId, token);

  const [overview, setOverview] = useState<MachineOverview | null>(null);
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [commandLog, setCommandLog] = useState<CommandLog[]>([]);
  const [liveSeries, setLiveSeries] = useState<Record<string, { t: number; v: number }[]>>({});
  const [actionError, setActionError] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState<CommandType | null>(null);

  useEffect(() => {
    getMachineOverview(machineId).then(setOverview).catch(() => {});
    listAlerts(machineId, false).then(setAlerts).catch(() => {});
    getCommandHistory(machineId, 10).then(setCommandLog).catch(() => {});
  }, [machineId]);

  useEffect(() => {
    if (!latest) return;

    setLiveSeries((prev) => {
      const next = { ...prev };
      for (const [metric, value] of Object.entries(latest.telemetry)) {
        const series = next[metric] ? [...next[metric]] : [];
        series.push({ t: Date.parse(latest.timestamp), v: value });
        next[metric] = series.slice(-MAX_LIVE_POINTS);
      }
      return next;
    });

    if (latest.ai_health) {
      setOverview((prev) =>
        prev
          ? {
              ...prev,
              current_health_score: latest.ai_health!.health_score,
              rul_hours: latest.ai_health!.rul_hours,
            }
          : prev
      );
    }
  }, [latest]);

  async function handleCommand(command: CommandType) {
    setActionError(null);
    setActionLoading(command);
    try {
      await executeCommand(machineId, command, `Triggered from dashboard: ${command}`);
      const [freshOverview, freshHistory] = await Promise.all([
        getMachineOverview(machineId),
        getCommandHistory(machineId, 10),
      ]);
      setOverview(freshOverview);
      setCommandLog(freshHistory);
    } catch {
      setActionError(`Failed to execute ${command}`);
    } finally {
      setActionLoading(null);
    }
  }

  if (!overview) {
    return (
      <AppShell>
        <p className="text-sm text-slate-500">Loading machine...</p>
      </AppShell>
    );
  }

  const chartMetrics = Object.keys(liveSeries);

  return (
    <AppShell>
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <h1 className="text-lg font-medium text-slate-900">{overview.name}</h1>
          <StatusBadge status={overview.status} />
        </div>
        <span className="flex items-center gap-1 text-xs text-slate-500">
          {connected ? <Wifi className="w-3.5 h-3.5 text-blue-600" /> : <WifiOff className="w-3.5 h-3.5 text-red-600" />}
          {connected ? "Live" : "Reconnecting..."}
        </span>
      </div>

      <div className="grid grid-cols-3 gap-3 mb-6">
        <StatBox label="Health score">
          <div className="flex items-center gap-3">
            <HealthGauge score={overview.current_health_score ?? 0} size={44} />
            <span className="text-xl font-medium text-slate-900">{overview.current_health_score?.toFixed(0) ?? "--"}</span>
          </div>
        </StatBox>
        <StatBox label="RUL hours">
          <p className="text-xl font-medium text-slate-900">{overview.rul_hours?.toFixed(1) ?? "--"}</p>
        </StatBox>
        <StatBox label="Active alerts">
          <p className="text-xl font-medium text-red-600">{overview.active_alerts_count}</p>
        </StatBox>
      </div>

      <div className="grid grid-cols-2 gap-3 mb-6">
        {chartMetrics.length === 0 && (
          <p className="text-sm text-slate-500 col-span-2">Waiting for live telemetry...</p>
        )}
        {chartMetrics.map((metric) => (
          <div key={metric} className="bg-white rounded-lg border border-slate-200 p-3 shadow-sm">
            <p className="text-xs text-slate-500 mb-2">{metric}</p>
            <ResponsiveContainer width="100%" height={100}>
              <LineChart data={liveSeries[metric]}>
                <XAxis dataKey="t" hide />
                <YAxis hide domain={["auto", "auto"]} />
                <Tooltip labelFormatter={(t) => new Date(Number(t)).toLocaleTimeString()} />
                <Line type="monotone" dataKey="v" stroke="#2563eb" dot={false} strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-2 gap-3 mb-6">
        <div className="bg-white rounded-lg border border-slate-200 p-4 shadow-sm">
          <p className="text-sm font-medium text-slate-900 mb-2">Active alerts</p>
          {alerts.length === 0 && <p className="text-sm text-slate-500">No active alerts.</p>}
          {alerts.map((a) => (
            <div key={a.alert_id} className="flex items-start gap-2 py-1.5 border-t border-slate-100 first:border-t-0">
              <OctagonAlert className="w-3.5 h-3.5 text-red-600 mt-0.5 shrink-0" />
              <div>
                <span className="text-xs text-red-600 font-medium">{a.severity}</span>
                <p className="text-sm text-slate-700">{a.message}</p>
              </div>
            </div>
          ))}
        </div>

        <div className="bg-white rounded-lg border border-slate-200 p-4 shadow-sm">
          <p className="text-sm font-medium text-slate-900 mb-2">Recent commands</p>
          {commandLog.length === 0 && <p className="text-sm text-slate-500">No commands yet.</p>}
          {commandLog.map((c) => (
            <div key={c.command_id} className="flex justify-between py-1.5 border-t border-slate-100 first:border-t-0 text-sm">
              <span className="text-slate-700">{c.command_type}</span>
              <span className="text-slate-400 text-xs">{c.issued_by}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="flex gap-2">
        <button
          onClick={() => handleCommand("EMERGENCY_STOP")}
          disabled={actionLoading !== null}
          className="flex-1 bg-red-600 text-white rounded-lg py-2 text-sm font-medium hover:bg-red-700 transition-colors disabled:opacity-50"
        >
          {actionLoading === "EMERGENCY_STOP" ? "Stopping..." : "Emergency stop"}
        </button>
        <button
          onClick={() => handleCommand("START")}
          disabled={actionLoading !== null}
          className="flex-1 bg-blue-600 text-white rounded-lg py-2 text-sm font-medium hover:bg-blue-700 transition-colors disabled:opacity-50"
        >
          {actionLoading === "START" ? "Starting..." : "Start"}
        </button>
        <button
          onClick={() => handleCommand("RESET_ALERTS")}
          disabled={actionLoading !== null}
          className="flex-1 bg-white border border-slate-300 text-slate-700 rounded-lg py-2 text-sm font-medium hover:bg-slate-50 transition-colors disabled:opacity-50"
        >
          {actionLoading === "RESET_ALERTS" ? "Resetting..." : "Reset alerts"}
        </button>
      </div>
      {actionError && <p className="text-sm text-red-600 mt-2">{actionError}</p>}
    </AppShell>
  );
}

function StatBox({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="bg-white rounded-lg border border-slate-200 p-3 shadow-sm">
      <p className="text-xs text-slate-500 mb-1">{label}</p>
      {children}
    </div>
  );
}