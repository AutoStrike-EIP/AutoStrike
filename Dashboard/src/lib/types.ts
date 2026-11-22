/**
 * Types partagés pour le Dashboard.
 *
 * Hypothèse de contrat avec le Back-end (ticket #14) : aucun schéma de données
 * officiel n'est encore défini côté RabbitMQ/PostgreSQL (ticket #20, statut
 * BACKLOG). Ces types sont une estimation raisonnable à partir de la note de
 * caractérisation (statuts "Créé / En cours d'exécution / Succès / Echec",
 * ticket #15) — à ajuster dès que le contrat réel du Back-end est connu.
 */

export type ScenarioStatus = "CREATED" | "RUNNING" | "SUCCESS" | "FAILED";

export interface ScenarioResult {
  id: string;
  name: string;
  status: ScenarioStatus;
  createdAt: string;
  updatedAt: string;
}

/** Message reçu depuis le WebSocket du Back-end pour un scénario. */
export interface ScenarioStatusEvent {
  type: "scenario.updated";
  payload: ScenarioResult;
}

export type BackendSocketEvent = ScenarioStatusEvent;

export type ConnectionStatus = "connecting" | "connected" | "disconnected";
