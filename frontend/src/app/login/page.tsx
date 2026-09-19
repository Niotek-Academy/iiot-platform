"use client";

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { ChartBar, Eye, EyeOff } from "lucide-react";
import { login as loginApi } from "@/api/auth";
import { useAuth } from "@/context/AuthContext";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    e.stopPropagation();
    setError(null);
    setLoading(true);
    
    try {
      const { token, user } = await loginApi(username, password);
      login(token, user);
      router.push("/");
    } catch (err: any) {
      console.log("Login error details:", err); // اطبع الخطأ لمعرفته
      setError("Invalid username or password");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50">
      {/* أضفنا onSubmit هنا وتأكدنا من منع السلوك الافتراضي */}
      <form onSubmit={handleSubmit} className="bg-white border border-slate-200 rounded-xl p-8 w-80 shadow-sm">
        <div className="flex items-center gap-2 mb-6">
          <div className="w-8 h-8 rounded-lg bg-blue-100 flex items-center justify-center">
            <ChartBar className="w-4 h-4 text-blue-600" />
          </div>
          <span className="font-medium text-slate-900">Niotek IIoT</span>
        </div>

        <label className="text-sm text-slate-500 block mb-1">Username</label>
        <input
          type="text"
          className="w-full mb-3 border border-slate-300 rounded-lg px-3 py-2 text-sm text-slate-900 bg-white"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />

        <label className="text-sm text-slate-500 block mb-1">Password</label>
        
        <div className="relative mb-5">
          <input
            type={showPassword ? "text" : "password"}
            className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm text-slate-900 bg-white pr-10"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
          >
            {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
          </button>
        </div>

        {error && <p className="text-sm text-red-600 mb-3">{error}</p>}

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white rounded-lg py-2 text-sm font-medium disabled:opacity-50"
        >
          {loading ? "Logging in..." : "Log in"}
        </button>
      </form>
    </div>
  );
}