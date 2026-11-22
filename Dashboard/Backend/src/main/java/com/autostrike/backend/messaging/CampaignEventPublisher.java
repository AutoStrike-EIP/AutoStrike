package com.autostrike.backend.messaging;

import com.autostrike.backend.campaign.Campaign;
import com.autostrike.backend.scenario.ScenarioTemplate;

public interface CampaignEventPublisher {

  /**
   * Annonce une nouvelle campagne au Task Manager, avec les tâches à créer.
   *
   * @throws org.springframework.amqp.AmqpException si le broker est injoignable
   */
  void campaignCreated(Campaign campaign, ScenarioTemplate template);
}
