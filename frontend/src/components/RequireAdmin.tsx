"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";

export default function RequireAdmin({ children }: { children: React.ReactNode }) {
  const { isAdmin } = useAuth();  // get the value of isAdmin from AuthContext
  const router = useRouter();

  useEffect(() => {
    if (!isAdmin) router.replace("/");
  }, [isAdmin, router]);

  if (!isAdmin) return null;
  return <>{children}</>;
}