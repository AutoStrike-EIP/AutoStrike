package com.autostrike.backend.messaging;

import com.autostrike.backend.campaign.CampaignService;
import com.autostrike.backend.config.RabbitConfig;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import java.util.UUID;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.AmqpRejectAndDontRequeueException;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;
import tools.jackson.core.JacksonException;
import tools.jackson.databind.json.JsonMapper;

/** Consomme {@code backend_notifications} ; le message est acquitté une fois traité. */
@Component
public class BackendNotificationsListener {

  private static final Logger log = LoggerFactory.getLogger(BackendNotificationsListener.class);

  /** Seul l'identifiant est lu : l'état fait foi en base, pas dans le message (ticket #20). */
  @JsonIgnoreProperties(ignoreUnknown = true)
  record CampaignUpdated(UUID id) {}

  private final CampaignService campaigns;
  private final JsonMapper jsonMapper;

  public BackendNotificationsListener(CampaignService campaigns, JsonMapper jsonMapper) {
    this.campaigns = campaigns;
    this.jsonMapper = jsonMapper;
  }

  @RabbitListener(queues = "${autostrike.rabbitmq.queue}")
  public void onMessage(Message message) {
    String routingKey = message.getMessageProperties().getReceivedRoutingKey();
    if (!RabbitConfig.CAMPAIGN_UPDATED.equals(routingKey)) {
      // task.updated : le Back-end ne gère pas encore les tâches.
      log.debug("Ignoring {} notification", routingKey);
      return;
    }
    CampaignUpdated event;
    try {
      event = jsonMapper.readValue(message.getBody(), CampaignUpdated.class);
    } catch (JacksonException e) {
      throw new AmqpRejectAndDontRequeueException("Malformed campaign.updated message", e);
    }
    if (event.id() == null) {
      throw new AmqpRejectAndDontRequeueException("campaign.updated message without id");
    }
    campaigns.onCampaignUpdated(event.id());
  }
}
