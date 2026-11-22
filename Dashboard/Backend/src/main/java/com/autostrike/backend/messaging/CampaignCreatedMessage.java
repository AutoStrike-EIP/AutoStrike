package com.autostrike.backend.messaging;

import com.autostrike.backend.campaign.Campaign;
import com.autostrike.backend.campaign.CampaignTarget;
import com.autostrike.backend.scenario.ScenarioTemplate;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;
import java.util.UUID;

/**
 * Corps de l'événement {@code campaign.created}, lu par le Task Manager (Rust). Reprend le schéma
 * de campagne du livrable « Scénario de campagne v1 », en snake_case, avec deux ajouts : l'{@code
 * id} de la campagne en base, et la criticité de chaque hôte (nécessaire à la priorisation). Schéma
 * provisoire jusqu'au ticket #20.
 */
public record CampaignCreatedMessage(Body campaign) {

  public record Body(
      UUID id,
      String name,
      String template,
      String mode,
      @JsonProperty("target_hosts") List<CampaignTarget> targetHosts,
      List<Task> tasks) {}

  @JsonInclude(JsonInclude.Include.NON_NULL)
  public record Task(
      int order,
      String technique,
      String action,
      @JsonProperty("depends_on") Integer dependsOn,
      Boolean cleanup) {}

  static CampaignCreatedMessage of(Campaign campaign, ScenarioTemplate template) {
    List<Task> tasks =
        template.steps().stream()
            .map(step ->
                new Task(
                    step.order(),
                    step.technique(),
                    step.action(),
                    step.dependsOn(),
                    step.cleanup() ? Boolean.TRUE : null))
            .toList();
    return new CampaignCreatedMessage(
        new Body(
            campaign.getId(),
            campaign.getName(),
            template.id(),
            template.mode(),
            campaign.getTargets(),
            tasks));
  }
}
