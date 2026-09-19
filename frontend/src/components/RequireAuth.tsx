"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";

// Next.js's App Router has no built-in route-guard primitive, so this
// wraps a page's content and redirects client-side if there's no token.
export default function RequireAuth({ children }: { children: React.ReactNode }) {
  const { token } = useAuth();  // get the value of token from AuthContext
  const router = useRouter();

  useEffect(() => {
    if (!token) router.replace("/login");
  }, [token, router]);

  if (!token) return null;
  return <>{children}</>;
}