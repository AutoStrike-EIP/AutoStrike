package com.autostrike.backend.messaging;

import static org.assertj.core.api.Assertions.assertThat;

import com.autostrike.backend.campaign.Campaign;
import com.autostrike.backend.campaign.CampaignTarget;
import com.autostrike.backend.campaign.HostCriticality;
import com.autostrike.backend.scenario.ScenarioCatalog;
import com.autostrike.backend.scenario.ScenarioTemplate;
import java.util.List;
import org.junit.jupiter.api.Test;
import tools.jackson.databind.JsonNode;
import tools.jackson.databind.json.JsonMapper;

/** Le Task Manager (Rust) lit ce JSON : ses noms de champs font partie du contrat. */
class CampaignCreatedMessageTest {

  private final JsonMapper jsonMapper = JsonMapper.builder().build();

  @Test
  void followsTheCampaignSchemaOfTheScenarioDocument() {
    ScenarioTemplate template =
        new ScenarioCatalog().findById(ScenarioCatalog.LINUX_BASIC_V1).orElseThrow();
    Campaign campaign =
        new Campaign(
            "Audit",
            template.id(),
            List.of(new CampaignTarget("10.0.0.12", HostCriticality.PRODUCTION)));

    JsonNode json =
        jsonMapper.readTree(
            jsonMapper.writeValueAsString(CampaignCreatedMessage.of(campaign, template)));

    JsonNode body = json.get("campaign");
    assertThat(body.get("name").asString()).isEqualTo("Audit");
    assertThat(body.get("template").asString()).isEqualTo("linux-basic-v1");
    assertThat(body.get("mode").asString()).isEqualTo("non-destructive");
    assertThat(body.get("target_hosts").get(0).get("host").asString()).isEqualTo("10.0.0.12");
    assertThat(body.get("target_hosts").get(0).get("criticality").asString())
        .isEqualTo("PRODUCTION");

    JsonNode tasks = body.get("tasks");
    assertThat(tasks.size()).isEqualTo(3);
    assertThat(tasks.get(0).get("technique").asString()).isEqualTo("T1087.001");
    assertThat(tasks.get(0).get("action").asString()).isEqualTo("list_local_accounts");
    assertThat(tasks.get(0).has("depends_on")).isFalse();
    assertThat(tasks.get(0).has("cleanup")).isFalse();
    assertThat(tasks.get(1).get("depends_on").asInt()).isEqualTo(1);
    assertThat(tasks.get(2).get("depends_on").asInt()).isEqualTo(2);
    assertThat(tasks.get(2).get("cleanup").asBoolean()).isTrue();
  }
}
