package com.autostrike.backend.messaging;

import com.autostrike.backend.campaign.Campaign;
import com.autostrike.backend.config.AutostrikeProperties;
import com.autostrike.backend.config.RabbitConfig;
import com.autostrike.backend.web.ScenarioResponse;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.core.MessageBuilder;
import org.springframework.amqp.core.MessageDeliveryMode;
import org.springframework.amqp.core.MessageProperties;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;
import tools.jackson.databind.json.JsonMapper;

/**
 * Publie du JSON brut et persistant, sans en-tête propre à Spring, pour rester lisible par les
 * services Rust. Le schéma du message est provisoire (ticket #20).
 */
@Component
public class RabbitCampaignEventPublisher implements CampaignEventPublisher {

  private final RabbitTemplate rabbitTemplate;
  private final JsonMapper jsonMapper;
  private final String exchange;

  public RabbitCampaignEventPublisher(
      RabbitTemplate rabbitTemplate, JsonMapper jsonMapper, AutostrikeProperties properties) {
    this.rabbitTemplate = rabbitTemplate;
    this.jsonMapper = jsonMapper;
    this.exchange = properties.rabbitmq().exchange();
  }

  @Override
  public void campaignCreated(Campaign campaign) {
    Message message =
        MessageBuilder.withBody(jsonMapper.writeValueAsBytes(ScenarioResponse.from(campaign)))
            .setContentType(MessageProperties.CONTENT_TYPE_JSON)
            .setDeliveryMode(MessageDeliveryMode.PERSISTENT)
            .build();
    rabbitTemplate.send(exchange, RabbitConfig.CAMPAIGN_CREATED, message);
  }
}
