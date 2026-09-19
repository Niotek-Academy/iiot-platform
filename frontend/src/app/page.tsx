"use client";

import RequireAuth from "@/components/RequireAuth";

export default function DashboardPage() {
  return (
    <RequireAuth>
      <div>Dashboard placeholder</div>
    </RequireAuth>
  );
}