package com.autostrike.backend.config;

import org.springframework.amqp.core.Binding;
import org.springframework.amqp.core.BindingBuilder;
import org.springframework.amqp.core.Declarables;
import org.springframework.amqp.core.Queue;
import org.springframework.amqp.core.QueueBuilder;
import org.springframework.amqp.core.TopicExchange;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * Topologie RabbitMQ vue du Back-end (docs/ARCHITECTURE.md §4 et §5) : il publie {@code
 * campaign.created} et lit {@code backend_notifications}, abonnée à {@code campaign.updated} et
 * {@code task.updated}.
 *
 * <p>L'exchange doit être déclaré à l'identique par tous les services, sinon RabbitMQ refuse la
 * seconde déclaration : topic et durable, comme dans {@code Server/common/src/event} (Rust). La
 * queue est durable, non exclusive et sans auto-suppression, comme celles du serveur.
 */
@Configuration
public class RabbitConfig {

  public static final String CAMPAIGN_CREATED = "campaign.created";
  public static final String CAMPAIGN_UPDATED = "campaign.updated";
  public static final String TASK_UPDATED = "task.updated";

  @Bean
  Declarables backendTopology(AutostrikeProperties properties) {
    TopicExchange exchange = new TopicExchange(properties.rabbitmq().exchange(), true, false);
    Queue queue = QueueBuilder.durable(properties.rabbitmq().queue()).build();
    Binding campaignUpdated = BindingBuilder.bind(queue).to(exchange).with(CAMPAIGN_UPDATED);
    Binding taskUpdated = BindingBuilder.bind(queue).to(exchange).with(TASK_UPDATED);
    return new Declarables(exchange, queue, campaignUpdated, taskUpdated);
  }
}
