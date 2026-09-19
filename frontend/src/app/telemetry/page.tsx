"use client";

import { useEffect, useState } from "react";
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";
import AppShell from "@/components/AppShell";
import { listMachines } from "@/api/machines";
import { listSensors } from "@/api/sensors";
import { listAlerts } from "@/api/alerts";
import { getTelemetryHistory } from "@/api/telemetry";
import type { Machine, Sensor, Alert, TelemetryDataPoint } from "@/types";

const SEVERITY_STYLES: Record<string, string> = {
  INFO: "bg-blue-100 text-blue-700",
  WARNING: "bg-yellow-100 text-yellow-700",
  CRITICAL: "bg-red-100 text-red-700",
};

export default function TelemetryHistoryPage() {
  const [machines, setMachines] = useState<Machine[]>([]);
  const [selectedMachine, setSelectedMachine] = useState("");
  const [sensors, setSensors] = useState<Sensor[]>([]);
  const [selectedSensor, setSelectedSensor] = useState("");
  const [points, setPoints] = useState<TelemetryDataPoint[]>([]);
  const [unit, setUnit] = useState("");
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [alertFilter, setAlertFilter] = useState<"all" | "resolved" | "unresolved">("all");
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    listMachines().then((list) => {
      setMachines(list);
      if (list.length > 0) setSelectedMachine(list[0].machine_id);
    });
  }, []);

  useEffect(() => {
    if (!selectedMachine) return;
    listSensors(selectedMachine).then((list) => {
      setSensors(list);
      if (list.length > 0) setSelectedSensor(list[0].sensor_id);
    });
    refreshAlerts();
  }, [selectedMachine, alertFilter]);

  useEffect(() => {
    if (!selectedSensor) return;
    getTelemetryHistory(selectedSensor, 100).then((res) => {
      setPoints(res.data_points);
      setUnit(res.unit);
    });
  }, [selectedSensor]);

  function refreshAlerts() {
    const isResolved = alertFilter === "all" ? undefined : alertFilter === "resolved";
    listAlerts(selectedMachine, isResolved).then(setAlerts).catch(() => {});
  }

  const chartData = points.map((p) => ({ t: Date.parse(p.recorded_at), v: p.value }));

  const filteredAlerts = alerts.filter((a) =>
    a.message.toLowerCase().includes(searchQuery.toLowerCase()) ||
    a.severity.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <AppShell searchQuery={searchQuery} setSearchQuery={setSearchQuery}>
      <h1 className="text-lg font-medium text-slate-900 mb-4">Telemetry & alerts history</h1>

      <div className="flex gap-2 mb-4">
        <select
          value={selectedMachine}
          onChange={(e) => setSelectedMachine(e.target.value)}
          className="border border-slate-300 rounded-lg px-3 py-1.5 text-sm text-slate-900 bg-white focus:outline-none focus:border-blue-600 shadow-sm"
        >
          {machines.map((m) => <option key={m.machine_id} value={m.machine_id}>{m.name}</option>)}
        </select>
        <select
          value={selectedSensor}
          onChange={(e) => setSelectedSensor(e.target.value)}
          className="border border-slate-300 rounded-lg px-3 py-1.5 text-sm text-slate-900 bg-white focus:outline-none focus:border-blue-600 shadow-sm"
        >
          {sensors.map((s) => <option key={s.sensor_id} value={s.sensor_id}>{s.sensor_id}</option>)}
        </select>
      </div>

      <div className="bg-white rounded-lg border border-slate-200 p-4 mb-6 shadow-sm">
        <p className="text-sm font-medium text-slate-900 mb-2">{selectedSensor} history {unit && `(${unit})`}</p>
        {chartData.length === 0 ? (
          <p className="text-sm text-slate-500">No data yet for this sensor.</p>
        ) : (
          <ResponsiveContainer width="100%" height={220}>
            <LineChart data={chartData}>
              <XAxis dataKey="t" tickFormatter={(t) => new Date(Number(t)).toLocaleTimeString()} tick={{ fontSize: 11 }} />
              <YAxis tick={{ fontSize: 11 }} />
              <Tooltip labelFormatter={(t) => new Date(Number(t)).toLocaleString()} />
              <Line type="monotone" dataKey="v" stroke="#2563eb" dot={false} strokeWidth={2} />
            </LineChart>
          </ResponsiveContainer>
        )}
      </div>

      <div className="flex justify-between items-center mb-2">
        <p className="text-sm font-medium text-slate-900">Alerts</p>
        <div className="flex gap-1">
          {(["all", "unresolved", "resolved"] as const).map((f) => (
            <button
              key={f}
              onClick={() => setAlertFilter(f)}
              className={`text-xs px-3 py-1 rounded-lg font-medium transition-colors ${
                alertFilter === f
                  ? "bg-blue-600 text-white shadow-sm"
                  : "bg-white border border-slate-300 text-slate-700 hover:bg-slate-50"
              }`}
            >
              {f[0].toUpperCase() + f.slice(1)}
            </button>
          ))}
        </div>
      </div>

      <div className="bg-white border border-slate-200 rounded-lg overflow-hidden shadow-sm">
        <div className="flex px-4 py-2 bg-slate-50 text-xs font-medium text-slate-500">
          <span className="flex-1">Severity</span>
          <span className="flex-[3]">Message</span>
          <span className="flex-1">Time</span>
          <span className="flex-1">Status</span>
        </div>
        {filteredAlerts.map((a) => (
          <div key={a.alert_id} className="flex items-center px-4 py-2 border-t border-slate-100 text-sm text-slate-800">
            <span className="flex-1">
              <span className={`text-xs px-2 py-0.5 rounded-md font-medium ${SEVERITY_STYLES[a.severity]}`}>
                {a.severity}
              </span>
            </span>
            <span className="flex-[3] text-slate-900 font-medium">{a.message}</span>
            <span className="flex-1 text-xs text-slate-500">{new Date(a.created_at).toLocaleString()}</span>
            <span className={`flex-1 text-xs font-medium ${a.is_resolved ? "text-slate-500" : "text-blue-600"}`}>
              {a.is_resolved ? "Resolved" : "Open"}
            </span>
          </div>
        ))}
        {filteredAlerts.length === 0 && <p className="p-4 text-sm text-slate-500">No alerts match this filter.</p>}
      </div>
    </AppShell>
  );
}