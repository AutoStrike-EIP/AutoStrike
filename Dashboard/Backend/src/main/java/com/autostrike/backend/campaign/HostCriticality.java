package com.autostrike.backend.campaign;

/**
 * Criticité métier d'un hôte, déclarée au lancement de la campagne. Sert à départager les résultats
 * qualifiés à profondeur égale (étage 2 du filtre de priorisation).
 */
public enum HostCriticality {
  PRODUCTION,
  TEST
}
