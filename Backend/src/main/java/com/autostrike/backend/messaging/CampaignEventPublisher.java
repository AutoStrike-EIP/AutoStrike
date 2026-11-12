package com.autostrike.backend.messaging;

import com.autostrike.backend.campaign.Campaign;

public interface CampaignEventPublisher {

  /**
   * Annonce une nouvelle campagne au Task Manager.
   *
   * @throws org.springframework.amqp.AmqpException si le broker est injoignable
   */
  void campaignCreated(Campaign campaign);
}
