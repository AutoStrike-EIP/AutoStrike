package com.autostrike.backend.config;

import org.springframework.amqp.core.Binding;
import org.springframework.amqp.core.BindingBuilder;
import org.springframework.amqp.core.Declarables;
import org.springframework.amqp.core.DirectExchange;
import org.springframework.amqp.core.Queue;
import org.springframework.amqp.core.QueueBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * Topologie RabbitMQ vue du Back-end (docs/ARCHITECTURE.md §4 et §5) : il publie {@code
 * campaign.created} et lit {@code backend_notifications}, abonnée à {@code campaign.updated} et
 * {@code task.updated}. Exchange et queue sont durables.
 */
@Configuration
public class RabbitConfig {

  public static final String CAMPAIGN_CREATED = "campaign.created";
  public static final String CAMPAIGN_UPDATED = "campaign.updated";
  public static final String TASK_UPDATED = "task.updated";

  @Bean
  Declarables backendTopology(AutostrikeProperties properties) {
    DirectExchange exchange = new DirectExchange(properties.rabbitmq().exchange(), true, false);
    Queue queue = QueueBuilder.durable(properties.rabbitmq().queue()).build();
    Binding campaignUpdated = BindingBuilder.bind(queue).to(exchange).with(CAMPAIGN_UPDATED);
    Binding taskUpdated = BindingBuilder.bind(queue).to(exchange).with(TASK_UPDATED);
    return new Declarables(exchange, queue, campaignUpdated, taskUpdated);
  }
}
