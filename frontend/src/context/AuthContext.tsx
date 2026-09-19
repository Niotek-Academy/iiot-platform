"use client";

import { createContext, useContext, useState, useCallback, type ReactNode } from "react";
import { jwtDecode } from "jwt-decode";
import { setAuthToken } from "@/api/client";
import type { Role, User } from "@/types";

// Define the structure of the JWT token data
interface JwtClaims {
  user_id: number;
  username: string;
  role: Role;
  exp: number;
}

// Define what values and functions the Auth Context will provide
interface AuthContextValue {
  user: User | null;
  token: string | null;
  login: (token: string, user: User) => void;
  logout: () => void;
  isAdmin: boolean;
}

// Create the context
const AuthContext = createContext<AuthContextValue | null>(null);

// Provider component to wrap around the app and manage auth state
export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<User | null>(null);

  // Login function to save token and user data
  const login = useCallback((newToken: string, newUser: User) => {
    setAuthToken(newToken);
    setToken(newToken);
    setUser(newUser);
  }, []);

  // Logout function to clear token and user data
  const logout = useCallback(() => {
    setAuthToken(null);
    setToken(null);
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isAdmin: user?.role === "ADMIN" }}>
      {children}
    </AuthContext.Provider>
  );
}

// Custom hook to easily use the auth context in components
export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider");
  return ctx;
}

// Helper function to decode the JWT token
export function decodeToken(token: string): JwtClaims {
  return jwtDecode<JwtClaims>(token);
}