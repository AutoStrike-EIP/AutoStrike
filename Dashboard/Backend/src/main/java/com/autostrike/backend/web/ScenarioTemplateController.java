package com.autostrike.backend.web;

import com.autostrike.backend.scenario.ScenarioCatalog;
import com.autostrike.backend.scenario.ScenarioTemplate;
import java.util.List;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ScenarioTemplateController {

  private final ScenarioCatalog catalog;

  public ScenarioTemplateController(ScenarioCatalog catalog) {
    this.catalog = catalog;
  }

  @GetMapping("/scenario-templates")
  public List<ScenarioTemplate> list() {
    return catalog.findAll();
  }
}
