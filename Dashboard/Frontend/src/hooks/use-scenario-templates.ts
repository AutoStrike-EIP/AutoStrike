"use client";

import { useEffect, useState } from "react";
import { API_URL } from "@/lib/config";
import type { ScenarioTemplate } from "@/lib/types";

/**
 * Charge une fois les modèles de scénarios (`GET /scenario-templates`).
 * `error` reste à `true` si le Back-end est injoignable au chargement.
 */
export function useScenarioTemplates() {
  const [templates, setTemplates] = useState<ScenarioTemplate[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    let cancelled = false;
    fetch(`${API_URL}/scenario-templates`)
      .then((res) => (res.ok ? res.json() : Promise.reject(res.status)))
      .then((data: ScenarioTemplate[]) => {
        if (!cancelled) setTemplates(data);
      })
      .catch(() => {
        if (!cancelled) setError(true);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  return { templates, loading, error };
}
