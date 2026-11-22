/**
 * Types partagés pour le Dashboard, alignés sur le contrat du Back-end
 * (Dashboard/Backend/README.md, section « Contrat exposé au Dashboard »).
 * Le schéma des données reste provisoire tant que le ticket #20 n'est pas tranché.
 */

export type ScenarioStatus = "CREATED" | "RUNNING" | "SUCCESS" | "FAILED";

/** Criticité métier d'un hôte : départage les résultats dans la priorisation. */
export type HostCriticality = "PRODUCTION" | "TEST";

export interface ScenarioTarget {
  host: string;
  criticality: HostCriticality;
}

export interface ScenarioResult {
  id: string;
  name: string;
  templateId: string;
  targets: ScenarioTarget[];
  status: ScenarioStatus;
  createdAt: string;
  updatedAt: string;
}

/** Une étape d'un modèle de scénario : une technique MITRE ATT&CK. */
export interface ScenarioStep {
  order: number;
  tactic: string;
  technique: string;
  techniqueName: string;
  action: string;
  dependsOn: number | null;
  cleanup: boolean;
}

export interface ScenarioTemplate {
  id: string;
  name: string;
  description: string;
  mode: string;
  steps: ScenarioStep[];
}

/** Corps de `POST /scenarios`. */
export interface LaunchScenarioRequest {
  templateId: string;
  name?: string;
  targets: ScenarioTarget[];
}

/** Message reçu depuis le WebSocket du Back-end pour un scénario. */
export interface ScenarioStatusEvent {
  type: "scenario.updated";
  payload: ScenarioResult;
}

export type BackendSocketEvent = ScenarioStatusEvent;

export type ConnectionStatus = "connecting" | "connected" | "disconnected";
