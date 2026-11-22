package com.autostrike.backend.config;

import com.autostrike.backend.ws.ScenarioSocketHandler;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.socket.config.annotation.EnableWebSocket;
import org.springframework.web.socket.config.annotation.WebSocketConfigurer;
import org.springframework.web.socket.config.annotation.WebSocketHandlerRegistry;

@Configuration
@EnableWebSocket
public class WebSocketConfig implements WebSocketConfigurer {

  private final ScenarioSocketHandler handler;
  private final AutostrikeProperties properties;

  public WebSocketConfig(ScenarioSocketHandler handler, AutostrikeProperties properties) {
    this.handler = handler;
    this.properties = properties;
  }

  @Override
  public void registerWebSocketHandlers(WebSocketHandlerRegistry registry) {
    registry
        .addHandler(handler, "/scenarios")
        .setAllowedOrigins(properties.cors().allowedOrigins().toArray(String[]::new));
  }
}
