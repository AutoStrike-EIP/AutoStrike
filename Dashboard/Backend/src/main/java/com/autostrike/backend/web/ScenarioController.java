package com.autostrike.backend.web;

import com.autostrike.backend.campaign.CampaignService;
import jakarta.validation.Valid;
import java.util.List;
import java.util.UUID;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

@RestController
@RequestMapping("/scenarios")
public class ScenarioController {

  private final CampaignService campaigns;

  public ScenarioController(CampaignService campaigns) {
    this.campaigns = campaigns;
  }

  /**
   * {@code /scenarios} sert aussi de point d'entrée WebSocket : une requête {@code Upgrade:
   * websocket} n'est pas pour ce contrôleur et retombe sur le handler WebSocket.
   */
  @GetMapping(headers = "Upgrade!=websocket")
  public List<ScenarioResponse> list() {
    return campaigns.findAll().stream().map(ScenarioResponse::from).toList();
  }

  @GetMapping("/{id}")
  public ScenarioResponse get(@PathVariable UUID id) {
    return campaigns
        .findById(id)
        .map(ScenarioResponse::from)
        .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND));
  }

  @PostMapping
  @ResponseStatus(HttpStatus.CREATED)
  public ScenarioResponse launch(@Valid @RequestBody CreateScenarioRequest request) {
    return ScenarioResponse.from(campaigns.launch(request.name()));
  }
}
