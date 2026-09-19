import type { MachineStatus } from "@/types";

const STYLES: Record<MachineStatus, string> = {
  OPERATIONAL: "bg-green-100 text-green-700",
  WARNING: "bg-yellow-100 text-yellow-700",
  CRITICAL: "bg-red-100 text-red-700",
  STOPPED: "bg-red-100 text-red-700",
};

const LABELS: Record<MachineStatus, string> = {
  OPERATIONAL: "Operational",
  WARNING: "Warning",
  CRITICAL: "Critical",
  STOPPED: "Stopped",
};

export default function StatusBadge({ status }: { status: MachineStatus }) {
  return (
    <span className={`text-xs px-2.5 py-0.5 rounded-md font-medium ${STYLES[status]}`}>
      {LABELS[status]}
    </span>
  );
}