package com.autostrike.backend.config;

import java.util.List;
import org.springframework.boot.context.properties.ConfigurationProperties;

/** Configuration propre à AutoStrike (préfixe {@code autostrike} dans application.yaml). */
@ConfigurationProperties(prefix = "autostrike")
public record AutostrikeProperties(Cors cors, Rabbitmq rabbitmq) {

  /** Origines autorisées à appeler l'API REST et le WebSocket (le Dashboard). */
  public record Cors(List<String> allowedOrigins) {}

  /** Exchange partagé entre les services et queue propre au Back-end. */
  public record Rabbitmq(String exchange, String queue) {}
}
