package com.autostrike.backend.web;

import com.autostrike.backend.campaign.Campaign;
import com.autostrike.backend.campaign.CampaignStatus;
import java.time.Instant;
import java.util.UUID;

/** Forme attendue par le Dashboard ({@code ScenarioResult} dans Dashboard/src/lib/types.ts). */
public record ScenarioResponse(
    UUID id, String name, CampaignStatus status, Instant createdAt, Instant updatedAt) {

  public static ScenarioResponse from(Campaign campaign) {
    return new ScenarioResponse(
        campaign.getId(),
        campaign.getName(),
        campaign.getStatus(),
        campaign.getCreatedAt(),
        campaign.getUpdatedAt());
  }
}
