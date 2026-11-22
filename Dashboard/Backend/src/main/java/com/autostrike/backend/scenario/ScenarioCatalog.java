package com.autostrike.backend.scenario;

import com.autostrike.backend.scenario.ScenarioTemplate.Step;
import java.util.List;
import java.util.Optional;
import org.springframework.stereotype.Component;

/**
 * Modèles de scénarios disponibles. Codés en dur tant qu'il n'en existe qu'un : le scénario v1 du
 * livrable « Scénario de campagne v1 (Linux, basique) », repris dans la note de caractérisation du
 * point dur (section 7).
 */
@Component
public class ScenarioCatalog {

  public static final String LINUX_BASIC_V1 = "linux-basic-v1";

  private static final List<ScenarioTemplate> TEMPLATES =
      List.of(
          new ScenarioTemplate(
              LINUX_BASIC_V1,
              "Linux basique (v1)",
              "Reconnaissance locale, identifiants en clair puis persistance par cron, depuis un"
                  + " compte utilisateur non-root. Lecture seule ou réversible.",
              "non-destructive",
              List.of(
                  new Step(
                      1,
                      "Discovery",
                      "T1087.001",
                      "Account Discovery: Local Account",
                      "list_local_accounts",
                      null,
                      false),
                  new Step(
                      2,
                      "Credential Access",
                      "T1552.001",
                      "Unsecured Credentials: Credentials In Files",
                      "grep_credentials_in_files",
                      1,
                      false),
                  new Step(
                      3,
                      "Persistence",
                      "T1053.003",
                      "Scheduled Task/Job: Cron",
                      "cron_persistence_test",
                      2,
                      true))));

  public List<ScenarioTemplate> findAll() {
    return TEMPLATES;
  }

  public Optional<ScenarioTemplate> findById(String id) {
    return TEMPLATES.stream().filter(template -> template.id().equals(id)).findFirst();
  }
}
