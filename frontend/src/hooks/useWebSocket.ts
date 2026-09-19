"use client";

import { useEffect, useRef, useState } from "react";
import type { WsPayload } from "@/types";

export function useMachineStream(machineId: string, token: string | null) {
  const [latest, setLatest] = useState<WsPayload | null>(null);
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!token) return;

    let cancelled = false;
    let reconnectTimer: ReturnType<typeof setTimeout>;

    function connect() {
      const url = `${process.env.NEXT_PUBLIC_WS_URL}?token=${token}&machine_id=${machineId}`;
      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => !cancelled && setConnected(true);

      ws.onmessage = (event) => {
        try {
          const payload: WsPayload = JSON.parse(event.data);
          if (!cancelled) setLatest(payload);
        } catch {
          // ignore malformed frames
        }
      };

      ws.onclose = () => {
        if (cancelled) return;
        setConnected(false);
        reconnectTimer = setTimeout(connect, 2000);
      };

      ws.onerror = () => ws.close();
    }

    connect();

    return () => {
      cancelled = true;
      clearTimeout(reconnectTimer);
      wsRef.current?.close();
    };
  }, [machineId, token]);

  return { latest, connected };
}