"use client";

import { useEffect, useState } from "react";
import { Plus, Trash2, Pencil } from "lucide-react";
import AppShell from "@/components/AppShell";
import RequireAdmin from "@/components/RequireAdmin";
import StatusBadge from "@/components/StatusBadge";
import { useAuth } from "@/context/AuthContext";
import { listMachines, createMachine, deleteMachine } from "@/api/machines";
import { listSensors, createSensor, deleteSensor } from "@/api/sensors";
import { listUsers, updateUserRole, registerUser } from "@/api/users";
import type { Machine, Sensor, User, Role } from "@/types";

type Tab = "machines" | "sensors" | "users";

export default function AdminPage() {
  const [tab, setTab] = useState<Tab>("machines");
  const [searchQuery, setSearchQuery] = useState("");

  return (
    <AppShell searchQuery={searchQuery} setSearchQuery={setSearchQuery}>
      <RequireAdmin>
        <h1 className="text-lg font-medium text-slate-900 mb-4">Admin management</h1>

        <div className="flex gap-2 mb-4">
          {(["machines", "sensors", "users"] as Tab[]).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={`text-sm px-4 py-1.5 rounded-lg font-medium ${
                tab === t ? "bg-blue-600 text-white shadow-sm" : "bg-white border border-slate-300 text-slate-700 hover:bg-slate-50"
              }`}
            >
              {t[0].toUpperCase() + t.slice(1)}
            </button>
          ))}
        </div>

        {tab === "machines" && <MachinesTab searchQuery={searchQuery} />}
        {tab === "sensors" && <SensorsTab />}
        {tab === "users" && <UsersTab />}
      </RequireAdmin>
    </AppShell>
  );
}

function MachinesTab({ searchQuery }: { searchQuery: string }) {
  const [machines, setMachines] = useState<Machine[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ machine_id: "", name: "", location: "", control_address: "" });
  const [error, setError] = useState<string | null>(null);

  function refresh() {
    listMachines().then(setMachines).catch(() => {});
  }

  useEffect(refresh, []);

  async function handleCreate() {
    setError(null);
    try {
      await createMachine({
        machine_id: form.machine_id,
        name: form.name,
        location: form.location,
        control_address: form.control_address || undefined,
      });
      setForm({ machine_id: "", name: "", location: "", control_address: "" });
      setShowForm(false);
      refresh();
    } catch {
      setError("Could not create machine — check the machine_id isn't already taken");
    }
  }

  async function handleDelete(machineId: string) {
    if (!confirm(`Delete machine ${machineId}?`)) return;
    await deleteMachine(machineId);
    refresh();
  }

  const filteredMachines = machines.filter((m) =>
    m.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    m.location.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div>
      <div className="flex justify-between items-center mb-2">
        <span className="text-sm font-medium text-slate-900">Machines</span>
        <button
          onClick={() => setShowForm((v) => !v)}
          className="flex items-center gap-1 bg-blue-600 text-white text-xs px-3 py-1.5 rounded-lg font-medium hover:bg-blue-700 transition-colors"
        >
          <Plus className="w-3.5 h-3.5" /> Add machine
        </button>
      </div>

      {showForm && (
        <div className="bg-white border border-slate-200 rounded-lg p-3 mb-3 grid grid-cols-4 gap-2 shadow-sm">
          <input placeholder="machine_id" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.machine_id} onChange={(e) => setForm({ ...form, machine_id: e.target.value })} />
          <input placeholder="name" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
          <input placeholder="location" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.location} onChange={(e) => setForm({ ...form, location: e.target.value })} />
          <input placeholder="control_address (optional)" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.control_address} onChange={(e) => setForm({ ...form, control_address: e.target.value })} />
          <button onClick={handleCreate} className="col-span-4 bg-blue-600 text-white text-sm rounded py-1.5 font-medium hover:bg-blue-700 transition-colors">Save</button>
          {error && <p className="col-span-4 text-xs text-red-600">{error}</p>}
        </div>
      )}

      <div className="border border-slate-200 rounded-lg overflow-hidden bg-white shadow-sm">
        <div className="flex px-4 py-2 bg-slate-50 text-xs font-medium text-slate-500">
          <span className="flex-[2]">Name</span>
          <span className="flex-1">Status</span>
          <span className="flex-1">Control address</span>
          <span className="flex-1">Actions</span>
        </div>
        {filteredMachines.map((m) => (
          <div key={m.machine_id} className="flex items-center px-4 py-2 border-t border-slate-100 text-sm text-slate-800">
            <span className="flex-[2] font-medium text-slate-900">{m.name}</span>
            <span className="flex-1"><StatusBadge status={m.status} /></span>
            <span className="flex-1 text-xs text-slate-500">{m.control_address ?? "—"}</span>
            <span className="flex-1">
              <button onClick={() => handleDelete(m.machine_id)} className="p-1 hover:bg-red-50 rounded transition-colors">
                <Trash2 className="w-4 h-4 text-red-600" />
              </button>
            </span>
          </div>
        ))}
        {filteredMachines.length === 0 && <p className="p-4 text-sm text-slate-500">No machines found.</p>}
      </div>
    </div>
  );
}

function SensorsTab() {
  const [machines, setMachines] = useState<Machine[]>([]);
  const [selectedMachine, setSelectedMachine] = useState("");
  const [sensors, setSensors] = useState<Sensor[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ sensor_id: "", metric_name: "", unit: "", source_address: "" });
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    listMachines().then((list) => {
      setMachines(list);
      if (list.length > 0) setSelectedMachine(list[0].machine_id);
    });
  }, []);

  function refresh() {
    if (selectedMachine) listSensors(selectedMachine).then(setSensors).catch(() => {});
  }

  useEffect(refresh, [selectedMachine]);

  async function handleCreate() {
    setError(null);
    try {
      await createSensor({
        sensor_id: form.sensor_id,
        machine_id: selectedMachine,
        metric_name: form.metric_name,
        unit: form.unit,
        source_address: form.source_address || undefined,
      });
      setForm({ sensor_id: "", metric_name: "", unit: "", source_address: "" });
      setShowForm(false);
      refresh();
    } catch {
      setError("Could not create sensor — check the sensor_id isn't already taken");
    }
  }

  async function handleDelete(sensorId: string) {
    if (!confirm(`Delete sensor ${sensorId}?`)) return;
    await deleteSensor(sensorId);
    refresh();
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-2">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-slate-900">Sensors —</span>
          <select value={selectedMachine} onChange={(e) => setSelectedMachine(e.target.value)} className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600 bg-white">
            {machines.map((m) => (
              <option key={m.machine_id} value={m.machine_id}>{m.name}</option>
            ))}
          </select>
        </div>
        <button onClick={() => setShowForm((v) => !v)} className="flex items-center gap-1 bg-blue-600 text-white text-xs px-3 py-1.5 rounded-lg font-medium hover:bg-blue-700 transition-colors">
          <Plus className="w-3.5 h-3.5" /> Add sensor
        </button>
      </div>

      {showForm && (
        <div className="bg-white border border-slate-200 rounded-lg p-3 mb-3 grid grid-cols-4 gap-2 shadow-sm">
          <input placeholder="sensor_id" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.sensor_id} onChange={(e) => setForm({ ...form, sensor_id: e.target.value })} />
          <input placeholder="metric_name" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.metric_name} onChange={(e) => setForm({ ...form, metric_name: e.target.value })} />
          <input placeholder="unit" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.unit} onChange={(e) => setForm({ ...form, unit: e.target.value })} />
          <input placeholder="source_address (optional)" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.source_address} onChange={(e) => setForm({ ...form, source_address: e.target.value })} />
          <button onClick={handleCreate} className="col-span-4 bg-blue-600 text-white text-sm rounded py-1.5 font-medium hover:bg-blue-700 transition-colors">Save</button>
          {error && <p className="col-span-4 text-xs text-red-600">{error}</p>}
        </div>
      )}

      <div className="border border-slate-200 rounded-lg overflow-hidden bg-white shadow-sm">
        <div className="flex px-4 py-2 bg-slate-50 text-xs font-medium text-slate-500">
          <span className="flex-1">Sensor ID</span>
          <span className="flex-[2]">Unit</span>
          <span className="flex-1">Actions</span>
        </div>
        {sensors.map((s) => (
          <div key={s.sensor_id} className="flex items-center px-4 py-2 border-t border-slate-100 text-sm text-slate-800">
            <span className="flex-1 font-medium text-slate-900">{s.sensor_id}</span>
            <span className="flex-[2] text-slate-600">{s.unit}</span>
            <span className="flex-1">
              <button onClick={() => handleDelete(s.sensor_id)} className="p-1 hover:bg-red-50 rounded transition-colors">
                <Trash2 className="w-4 h-4 text-red-600" />
              </button>
            </span>
          </div>
        ))}
        {sensors.length === 0 && <p className="p-4 text-sm text-slate-500">No sensors for this machine yet.</p>}
      </div>
    </div>
  );
}

function UsersTab() {
  const { user: currentUser } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ username: "", password: "", role: "OPERATOR" as Role });
  const [error, setError] = useState<string | null>(null);

  function refresh() {
    listUsers().then(setUsers).catch(() => {});
  }

  useEffect(refresh, []);

  async function handleRoleChange(userId: number, role: Role) {
    try {
      await updateUserRole(userId, role);
      refresh();
    } catch {
      setError("Could not update role — you cannot change your own role");
    }
  }

  async function handleCreate() {
    setError(null);
    
    if (form.password.length < 6) {
      setError("Password must be at least 6 characters long");
      return;
    }

    try {
      await registerUser(form.username, form.password, form.role);
      setForm({ username: "", password: "", role: "OPERATOR" });
      setShowForm(false);
      refresh();
    } catch (err: any) {
      setError("Could not create user — check if username is taken or password is too short");
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-2">
        <span className="text-sm font-medium text-slate-900">User management</span>
        <button onClick={() => setShowForm((v) => !v)} className="flex items-center gap-1 bg-blue-600 text-white text-xs px-3 py-1.5 rounded-lg font-medium hover:bg-blue-700 transition-colors">
          <Plus className="w-3.5 h-3.5" /> Add user
        </button>
      </div>

      {showForm && (
        <div className="bg-white border border-slate-200 rounded-lg p-3 mb-3 grid grid-cols-4 gap-2 shadow-sm">
          <input placeholder="username" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} />
          <input placeholder="password (min 6 chars)" type="password" className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} />
          <select className="border border-slate-300 rounded px-2 py-1 text-sm text-slate-900 focus:outline-none focus:border-blue-600 bg-white" value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value as Role })}>
            <option value="OPERATOR">Operator</option>
            <option value="ADMIN">Admin</option>
          </select>
          <button onClick={handleCreate} className="bg-blue-600 text-white text-sm rounded py-1.5 font-medium hover:bg-blue-700 transition-colors">Save</button>
          {error && <p className="col-span-4 text-xs text-red-600 mt-1">{error}</p>}
        </div>
      )}

      <div className="border border-slate-200 rounded-lg overflow-hidden bg-white shadow-sm">
        <div className="flex px-4 py-2 bg-slate-50 text-xs font-medium text-slate-500">
          <span className="flex-[2]">Username</span>
          <span className="flex-1">Role</span>
          <span className="flex-1">Actions</span>
        </div>
        {users.map((u) => (
          <div key={u.user_id} className="flex items-center px-4 py-2 border-t border-slate-100 text-sm text-slate-800">
            <span className="flex-[2] font-medium text-slate-900">{u.username}</span>
            <span className="flex-1">
              <span className={`text-xs px-2 py-0.5 rounded-md font-medium ${u.role === "ADMIN" ? "bg-blue-100 text-blue-700" : "bg-slate-100 text-slate-700"}`}>
                {u.role}
              </span>
            </span>
            <span className="flex-1">
              {u.user_id !== currentUser?.user_id ? (
                <button onClick={() => handleRoleChange(u.user_id, u.role === "ADMIN" ? "OPERATOR" : "ADMIN")} className="flex items-center gap-1 text-xs text-blue-600 hover:text-blue-800 font-medium">
                  <Pencil className="w-3 h-3" /> Toggle role
                </button>
              ) : (
                <span className="text-xs text-slate-400">(you)</span>
              )}
            </span>
          </div>
        ))}
        {users.length === 0 && <p className="p-4 text-sm text-slate-500">No users found.</p>}
      </div>
    </div>
  );
}