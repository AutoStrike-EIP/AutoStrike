package com.autostrike.backend.scenario;

import java.util.List;

/**
 * Un modèle de scénario : une chaîne fixe de techniques MITRE ATT&amp;CK que l'agent joue dans
 * l'ordre. Une campagne = un modèle appliqué à une liste d'hôtes cibles.
 *
 * @param mode toujours {@code non-destructive} pour l'instant (cadre juridique du projet)
 */
public record ScenarioTemplate(
    String id, String name, String description, String mode, List<Step> steps) {

  /**
   * Une étape de la chaîne.
   *
   * @param dependsOn ordre de l'étape dont celle-ci dépend, ou {@code null}
   * @param cleanup l'agent doit annuler l'action en fin de tâche
   */
  public record Step(
      int order,
      String tactic,
      String technique,
      String techniqueName,
      String action,
      Integer dependsOn,
      boolean cleanup) {}
}
