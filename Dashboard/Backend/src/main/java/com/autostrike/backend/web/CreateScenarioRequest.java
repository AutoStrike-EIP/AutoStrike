package com.autostrike.backend.web;

import com.autostrike.backend.campaign.CampaignTarget;
import com.autostrike.backend.campaign.HostCriticality;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import java.util.List;

/**
 * Lancement d'un scénario.
 *
 * @param name libellé libre ; le nom du modèle s'il est absent
 */
public record CreateScenarioRequest(
    @NotBlank String templateId,
    @Size(max = 255) String name,
    @NotEmpty @Size(max = 50) List<@Valid @NotNull Target> targets) {

  /**
   * Un hôte cible. Le nom d'hôte est restreint aux caractères d'un nom DNS ou d'une adresse IP : il
   * finira dans les tâches envoyées aux agents.
   */
  public record Target(
      @NotBlank @Size(max = 255) @Pattern(regexp = "[A-Za-z0-9.:-]+") String host,
      @NotNull HostCriticality criticality) {

    CampaignTarget toCampaignTarget() {
      return new CampaignTarget(host, criticality);
    }
  }

  List<CampaignTarget> campaignTargets() {
    return targets.stream().map(Target::toCampaignTarget).toList();
  }
}
