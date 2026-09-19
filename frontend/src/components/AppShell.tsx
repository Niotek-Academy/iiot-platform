"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChartBar, LayoutDashboard, FileText, Users, Bell } from "lucide-react";
import RequireAuth from "./RequireAuth";
import { useAuth } from "@/context/AuthContext";

const NAV_ITEMS = [
  { href: "/", label: "Live dashboard", icon: LayoutDashboard },
  { href: "/telemetry", label: "Telemetry logs", icon: FileText },
  { href: "/admin", label: "Admin management", icon: Users, adminOnly: true },
];

function Sidebar() {
  const pathname = usePathname();
  const { isAdmin } = useAuth();

  return (
    <div className="w-44 shrink-0 bg-slate-50 border-r border-slate-200 py-4">
      <div className="flex items-center gap-2 px-4 pb-4">
        <div className="w-7 h-7 rounded-lg bg-blue-100 flex items-center justify-center">
          <ChartBar className="w-4 h-4 text-blue-600" />
        </div>
        <span className="font-medium text-sm text-slate-900">Niotek IIoT</span>
      </div>

      {NAV_ITEMS.filter((item) => !item.adminOnly || isAdmin).map((item) => {
        const Icon = item.icon;
        const active = pathname === item.href;
        return (
          <Link
            key={item.href}
            href={item.href}
            className={`flex items-center gap-2 px-4 py-2 text-sm ${
              active ? "bg-blue-50 border-l-2 border-blue-600 text-blue-600" : "text-slate-600 hover:bg-slate-100"
            }`}
          >
            <Icon className="w-4 h-4" />
            {item.label}
          </Link>
        );
      })}
    </div>
  );
}

function TopBar({
  searchQuery,
  setSearchQuery,
  placeholder,
}: {
  searchQuery: string;
  setSearchQuery: (q: string) => void;
  placeholder: string;
}) {
  const { user, logout } = useAuth();

  return (
    <div className="flex items-center justify-between px-6 py-3 border-b border-slate-200 bg-white">
      <input
        type="text"
        placeholder={placeholder}
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        className="border border-slate-300 rounded-lg px-3 py-1.5 text-sm w-56 text-slate-900 focus:outline-none focus:border-blue-600"
      />
      <div className="flex items-center gap-4">
        <Bell className="w-4 h-4 text-slate-500" />
        <span className="text-sm text-slate-700">
          {user?.username} <span className="text-slate-400">({user?.role.toLowerCase()})</span>
        </span>
        <button onClick={logout} className="text-sm text-slate-500 hover:text-slate-800">
          Log out
        </button>
      </div>
    </div>
  );
}

interface AppShellProps {
  children: React.ReactNode;
  searchQuery?: string;
  setSearchQuery?: (q: string) => void;
  searchPlaceholder?: string;
}

export default function AppShell({
  children,
  searchQuery: controlledQuery,
  setSearchQuery: controlledSetQuery,
  searchPlaceholder = "Search machines",
}: AppShellProps) {
  const [internalQuery, setInternalQuery] = useState("");
  const searchQuery = controlledQuery ?? internalQuery;
  const setSearchQuery = controlledSetQuery ?? setInternalQuery;

  return (
    <RequireAuth>
      <div className="flex min-h-screen bg-slate-100">
        <Sidebar />
        <div className="flex-1 min-w-0 flex flex-col">
          <TopBar searchQuery={searchQuery} setSearchQuery={setSearchQuery} placeholder={searchPlaceholder} />
          <div className="p-6 flex-1 overflow-auto">{children}</div>
        </div>
      </div>
    </RequireAuth>
  );
}