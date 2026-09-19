"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { OctagonAlert } from "lucide-react";
import AppShell from "@/components/AppShell";
import HealthGauge from "@/components/HealthGauge";
import StatusBadge from "@/components/StatusBadge";
import { listMachines, getMachineOverview } from "@/api/machines";
import type { MachineOverview } from "@/types";

export default function DashboardPage() {
  const [overviews, setOverviews] = useState<MachineOverview[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState(""); // 1. أضفنا الـ state الخاصة بالبحث

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const machines = await listMachines();
        const results = await Promise.all(machines.map((m) => getMachineOverview(m.machine_id)));
        if (!cancelled) setOverviews(results);
      } catch {
        if (!cancelled) setError("Could not load machines");
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    load();
    const interval = setInterval(load, 10_000);
    return () => {
      cancelled = true;
      clearInterval(interval);
    };
  }, []);

  // 2. تصفية الماكينات بناءً على الحروف المكتوبة في السيرش
  const filteredMachines = overviews.filter((m) =>
    m.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const totalMachines = overviews.length;
  const running = overviews.filter((m) => m.status === "OPERATIONAL").length;
  const totalAlerts = overviews.reduce((sum, m) => sum + m.active_alerts_count, 0);
  const stopped = overviews.filter((m) => m.status === "STOPPED").length;

  return (
    <AppShell searchQuery={searchQuery} setSearchQuery={setSearchQuery}>
      <h1 className="text-lg font-medium text-slate-900 mb-4">Live dashboard</h1>

      {loading && <p className="text-sm text-slate-500">Loading machines...</p>}
      {error && <p className="text-sm text-red-600">{error}</p>}

      {!loading && !error && (
        <>
          <div className="grid grid-cols-4 gap-3 mb-6">
            <StatCard label="Total machines" value={totalMachines} />
            <StatCard label="Active running" value={running} color="text-green-600" />
            <StatCard label="Active alerts" value={totalAlerts} color="text-yellow-600" />
            <StatCard label="Stopped" value={stopped} color="text-red-600" />
          </div>

          <div className="grid grid-cols-3 gap-3">
            {/* 3. بنعرض الماكينات المصفاة (filteredMachines) بدل القائمة كلها */}
            {filteredMachines.map((m) => (
              <Link
                key={m.machine_id}
                href={`/machines/${m.machine_id}`}
                className="bg-white rounded-lg border border-slate-200 p-4 flex items-center gap-3 hover:border-blue-300 transition-colors"
              >
                <HealthGauge score={m.current_health_score ?? 0} />
                <div className="min-w-0">
                  <p className="text-sm font-medium text-slate-900 truncate">{m.name}</p>
                  <StatusBadge status={m.status} />
                  {m.active_alerts_count > 0 && (
                    <span className="flex items-center gap-1 text-xs text-red-600 mt-1">
                      <OctagonAlert className="w-3 h-3" />
                      {m.active_alerts_count} active alert{m.active_alerts_count > 1 ? "s" : ""}
                    </span>
                  )}
                </div>
              </Link>
            ))}
          </div>

          {filteredMachines.length === 0 && (
            <p className="text-sm text-slate-500 mt-4">No matching machines found.</p>
          )}
        </>
      )}
    </AppShell>
  );
}

function StatCard({ label, value, color }: { label: string; value: number; color?: string }) {
  return (
    <div className="bg-white rounded-lg border border-slate-200 p-3">
      <p className="text-xs text-slate-500">{label}</p>
      <p className={`text-xl font-medium mt-0.5 ${color ?? "text-slate-900"}`}>{value}</p>
    </div>
  );
}