"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type {
  BackendSocketEvent,
  ConnectionStatus,
  ScenarioResult,
} from "@/lib/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3001";
const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:3001";

const RECONNECT_DELAYS_MS = [1000, 2000, 5000, 10000, 15000];

/**
 * Récupère la liste des scénarios (GET #14) puis se connecte au WebSocket du
 * Back-end pour recevoir les mises à jour en direct. Se reconnecte
 * automatiquement avec un backoff si la connexion tombe (le Back-end n'est
 * pas encore déployé au moment de l'écriture de ce hook — ticket #14 encore
 * TODO — donc l'échec de connexion doit rester silencieux et visible dans
 * `status`, jamais bloquant pour le reste de la page).
 */
export function useBackendSocket() {
  const [scenarios, setScenarios] = useState<ScenarioResult[]>([]);
  const [status, setStatus] = useState<ConnectionStatus>("connecting");
  const [launchError, setLaunchError] = useState<string | null>(null);

  const socketRef = useRef<WebSocket | null>(null);
  const reconnectAttemptRef = useRef(0);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(
    null,
  );
  const unmountedRef = useRef(false);
  // Évite l'auto-référence de `connect` dans sa propre définition : la
  // reconnexion planifiée appelle toujours la dernière version via la ref.
  const connectRef = useRef<() => void>(() => {});

  const upsertScenario = useCallback((scenario: ScenarioResult) => {
    setScenarios((prev) => {
      const index = prev.findIndex((s) => s.id === scenario.id);
      if (index === -1) return [scenario, ...prev];
      const next = [...prev];
      next[index] = scenario;
      return next;
    });
  }, []);

  const connect = useCallback(() => {
    if (unmountedRef.current) return;

    setStatus((prev) => (prev === "connected" ? prev : "connecting"));

    const socket = new WebSocket(`${WS_URL}/scenarios`);
    socketRef.current = socket;

    socket.onopen = () => {
      reconnectAttemptRef.current = 0;
      setStatus("connected");
    };

    socket.onmessage = (event) => {
      try {
        const data: BackendSocketEvent = JSON.parse(event.data);
        if (data.type === "scenario.updated") {
          upsertScenario(data.payload);
        }
      } catch {
        // Message non conforme au contrat attendu — ignoré plutôt que de
        // planter la page, tant que le schéma n'est pas figé (ticket #20).
      }
    };

    const scheduleReconnect = () => {
      if (unmountedRef.current) return;
      setStatus("disconnected");
      const delay =
        RECONNECT_DELAYS_MS[
          Math.min(reconnectAttemptRef.current, RECONNECT_DELAYS_MS.length - 1)
        ];
      reconnectAttemptRef.current += 1;
      reconnectTimeoutRef.current = setTimeout(() => connectRef.current(), delay);
    };

    socket.onclose = scheduleReconnect;
    socket.onerror = () => {
      socket.close();
    };
  }, [upsertScenario]);

  useEffect(() => {
    connectRef.current = connect;
  }, [connect]);

  useEffect(() => {
    unmountedRef.current = false;

    // Chargement initial via REST (GET /scenarios, ticket #14) avant de
    // basculer en live sur le WebSocket.
    fetch(`${API_URL}/scenarios`)
      .then((res) => (res.ok ? res.json() : Promise.reject(res.status)))
      .then((data: ScenarioResult[]) => setScenarios(data))
      .catch(() => {
        // Back-end indisponible au chargement : on part d'une liste vide,
        // le WebSocket complètera dès qu'une connexion sera établie.
      });

    connect();

    return () => {
      unmountedRef.current = true;
      if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);
      socketRef.current?.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps -- ne doit tourner qu'au montage
  }, []);

  const launchScenario = useCallback(async (name: string) => {
    setLaunchError(null);
    try {
      const res = await fetch(`${API_URL}/scenarios`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name }),
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const created: ScenarioResult = await res.json();
      upsertScenario(created);
    } catch (err) {
      setLaunchError(
        err instanceof Error
          ? err.message
          : "Impossible de contacter le back-end",
      );
    }
  }, [upsertScenario]);

  return { scenarios, status, launchScenario, launchError };
}
