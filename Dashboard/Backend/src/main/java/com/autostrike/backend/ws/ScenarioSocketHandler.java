package com.autostrike.backend.ws;

import com.autostrike.backend.campaign.Campaign;
import com.autostrike.backend.web.ScenarioResponse;
import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import org.springframework.web.socket.CloseStatus;
import org.springframework.web.socket.TextMessage;
import org.springframework.web.socket.WebSocketSession;
import org.springframework.web.socket.handler.ConcurrentWebSocketSessionDecorator;
import org.springframework.web.socket.handler.TextWebSocketHandler;
import tools.jackson.databind.json.JsonMapper;

/** Diffuse l'état des scénarios à tous les clients connectés sur {@code WS /scenarios}. */
@Component
public class ScenarioSocketHandler extends TextWebSocketHandler {

  private static final Logger log = LoggerFactory.getLogger(ScenarioSocketHandler.class);
  private static final int SEND_TIME_LIMIT_MS = 5_000;
  private static final int SEND_BUFFER_LIMIT_BYTES = 64 * 1024;

  /** Message envoyé aux clients ({@code ScenarioStatusEvent} côté Dashboard). */
  record ScenarioUpdated(String type, ScenarioResponse payload) {
    ScenarioUpdated(ScenarioResponse payload) {
      this("scenario.updated", payload);
    }
  }

  private final Map<String, WebSocketSession> sessions = new ConcurrentHashMap<>();
  private final JsonMapper jsonMapper;

  public ScenarioSocketHandler(JsonMapper jsonMapper) {
    this.jsonMapper = jsonMapper;
  }

  @Override
  public void afterConnectionEstablished(WebSocketSession session) {
    // Le décorateur sérialise les envois : broadcast() est appelé depuis plusieurs threads.
    sessions.put(
        session.getId(),
        new ConcurrentWebSocketSessionDecorator(
            session, SEND_TIME_LIMIT_MS, SEND_BUFFER_LIMIT_BYTES));
  }

  @Override
  public void afterConnectionClosed(WebSocketSession session, CloseStatus status) {
    sessions.remove(session.getId());
  }

  public void broadcast(Campaign campaign) {
    TextMessage message =
        new TextMessage(
            jsonMapper.writeValueAsString(new ScenarioUpdated(ScenarioResponse.from(campaign))));
    for (WebSocketSession session : sessions.values()) {
      try {
        session.sendMessage(message);
      } catch (IOException | RuntimeException e) {
        log.warn("Dropping WebSocket session {}: {}", session.getId(), e.getMessage());
        sessions.remove(session.getId());
      }
    }
  }
}
