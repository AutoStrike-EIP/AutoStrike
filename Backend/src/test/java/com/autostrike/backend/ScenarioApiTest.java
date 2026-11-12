package com.autostrike.backend;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.verify;

import com.autostrike.backend.campaign.Campaign;
import com.autostrike.backend.campaign.CampaignRepository;
import com.autostrike.backend.config.RabbitConfig;
import com.autostrike.backend.messaging.BackendNotificationsListener;
import com.autostrike.backend.messaging.CampaignEventPublisher;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.http.WebSocket;
import java.nio.charset.StandardCharsets;
import java.util.UUID;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.CompletionStage;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.TimeUnit;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.amqp.AmqpConnectException;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.core.MessageBuilder;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import tools.jackson.databind.JsonNode;
import tools.jackson.databind.json.JsonMapper;

/** Vérifie le contrat REST + WebSocket attendu par le Dashboard, sur un vrai serveur HTTP. */
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@ActiveProfiles("test")
class ScenarioApiTest {

  @Value("${local.server.port}")
  int port;

  @Autowired CampaignRepository repository;
  @Autowired BackendNotificationsListener listener;
  @Autowired JsonMapper jsonMapper;
  @MockitoBean CampaignEventPublisher publisher;

  final HttpClient http = HttpClient.newHttpClient();
  WebSocket socket;

  @BeforeEach
  void cleanDatabase() {
    repository.deleteAll();
  }

  @AfterEach
  void closeSocket() {
    if (socket != null) {
      socket.abort();
    }
  }

  @Test
  void postCreatesScenarioAndAnnouncesItToTheTaskManager() throws Exception {
    HttpResponse<String> response = post("{\"name\":\"Port scan\"}");

    assertThat(response.statusCode()).isEqualTo(201);
    JsonNode body = jsonMapper.readTree(response.body());
    assertThat(body.get("id").asString()).isNotBlank();
    assertThat(body.get("name").asString()).isEqualTo("Port scan");
    assertThat(body.get("status").asString()).isEqualTo("CREATED");
    assertThat(body.get("createdAt").asString()).matches("\\d{4}-\\d{2}-\\d{2}T.*Z");
    assertThat(body.get("updatedAt").asString()).matches("\\d{4}-\\d{2}-\\d{2}T.*Z");
    verify(publisher).campaignCreated(any(Campaign.class));
  }

  @Test
  void getReturnsScenariosMostRecentFirst() throws Exception {
    post("{\"name\":\"first\"}");
    Thread.sleep(5);
    post("{\"name\":\"second\"}");

    HttpResponse<String> response = get("/scenarios");

    assertThat(response.statusCode()).isEqualTo(200);
    JsonNode body = jsonMapper.readTree(response.body());
    assertThat(body.size()).isEqualTo(2);
    assertThat(body.get(0).get("name").asString()).isEqualTo("second");
    assertThat(body.get(1).get("name").asString()).isEqualTo("first");
  }

  @Test
  void getByIdReturnsTheScenarioOr404() throws Exception {
    String id = jsonMapper.readTree(post("{\"name\":\"one\"}").body()).get("id").asString();

    assertThat(jsonMapper.readTree(get("/scenarios/" + id).body()).get("name").asString())
        .isEqualTo("one");
    assertThat(get("/scenarios/00000000-0000-0000-0000-000000000000").statusCode()).isEqualTo(404);
  }

  @Test
  void postRejectsBlankName() throws Exception {
    assertThat(post("{\"name\":\"  \"}").statusCode()).isEqualTo(400);
    assertThat(repository.count()).isZero();
  }

  @Test
  void postFailsScenarioWhenBrokerIsDown() throws Exception {
    doThrow(new AmqpConnectException(new RuntimeException("broker down")))
        .when(publisher)
        .campaignCreated(any(Campaign.class));

    assertThat(post("{\"name\":\"lost\"}").statusCode()).isEqualTo(503);

    JsonNode body = jsonMapper.readTree(get("/scenarios").body());
    assertThat(body.get(0).get("status").asString()).isEqualTo("FAILED");
  }

  @Test
  void websocketBroadcastsCreationThenStatusChanges() throws Exception {
    BlockingQueue<String> received = connectWebSocket();

    String id = jsonMapper.readTree(post("{\"name\":\"live\"}").body()).get("id").asString();
    JsonNode created = jsonMapper.readTree(received.poll(5, TimeUnit.SECONDS));
    assertThat(created.get("type").asString()).isEqualTo("scenario.updated");
    assertThat(created.get("payload").get("id").asString()).isEqualTo(id);
    assertThat(created.get("payload").get("status").asString()).isEqualTo("CREATED");

    // Le Task Manager met le statut à jour en base, puis publie campaign.updated.
    Campaign campaign = repository.findById(UUID.fromString(id)).orElseThrow();
    campaign.markFailed();
    repository.save(campaign);
    listener.onMessage(notification(RabbitConfig.CAMPAIGN_UPDATED, "{\"id\":\"" + id + "\"}"));

    JsonNode updated = jsonMapper.readTree(received.poll(5, TimeUnit.SECONDS));
    assertThat(updated.get("type").asString()).isEqualTo("scenario.updated");
    assertThat(updated.get("payload").get("id").asString()).isEqualTo(id);
    assertThat(updated.get("payload").get("status").asString()).isEqualTo("FAILED");
  }

  @Test
  void taskUpdatedNotificationsAreIgnored() throws Exception {
    BlockingQueue<String> received = connectWebSocket();

    listener.onMessage(notification(RabbitConfig.TASK_UPDATED, "{\"id\":\"whatever\"}"));

    assertThat(received.poll(300, TimeUnit.MILLISECONDS)).isNull();
  }

  private HttpResponse<String> post(String json) throws Exception {
    HttpRequest request =
        HttpRequest.newBuilder(URI.create("http://localhost:" + port + "/scenarios"))
            .header("Content-Type", "application/json")
            .POST(HttpRequest.BodyPublishers.ofString(json))
            .build();
    return http.send(request, HttpResponse.BodyHandlers.ofString());
  }

  private HttpResponse<String> get(String path) throws Exception {
    HttpRequest request =
        HttpRequest.newBuilder(URI.create("http://localhost:" + port + path)).GET().build();
    return http.send(request, HttpResponse.BodyHandlers.ofString());
  }

  private static Message notification(String routingKey, String json) {
    return MessageBuilder.withBody(json.getBytes(StandardCharsets.UTF_8))
        .setReceivedRoutingKey(routingKey)
        .build();
  }

  private BlockingQueue<String> connectWebSocket() {
    BlockingQueue<String> received = new LinkedBlockingQueue<>();
    WebSocket.Listener collector =
        new WebSocket.Listener() {
          private final StringBuilder buffer = new StringBuilder();

          @Override
          public CompletionStage<?> onText(WebSocket webSocket, CharSequence data, boolean last) {
            buffer.append(data);
            if (last) {
              received.add(buffer.toString());
              buffer.setLength(0);
            }
            webSocket.request(1);
            return null;
          }
        };
    socket =
        http.newWebSocketBuilder()
            .buildAsync(URI.create("ws://localhost:" + port + "/scenarios"), collector)
            .join();
    return received;
  }
}
