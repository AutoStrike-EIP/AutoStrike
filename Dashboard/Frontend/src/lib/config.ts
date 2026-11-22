/** URL de l'API REST du Back-end (voir .env.local.example). */
export const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3001";

/** URL du serveur WebSocket du Back-end. */
export const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:3001";
