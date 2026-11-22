package com.autostrike.backend.campaign;

import com.autostrike.backend.messaging.CampaignEventPublisher;
import com.autostrike.backend.scenario.ScenarioCatalog;
import com.autostrike.backend.scenario.ScenarioTemplate;
import com.autostrike.backend.ws.ScenarioSocketHandler;
import java.util.List;
import java.util.Optional;
import java.util.UUID;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.AmqpException;
import org.springframework.stereotype.Service;

@Service
public class CampaignService {

  private static final Logger log = LoggerFactory.getLogger(CampaignService.class);

  private final CampaignRepository repository;
  private final ScenarioCatalog catalog;
  private final CampaignEventPublisher publisher;
  private final ScenarioSocketHandler socket;

  public CampaignService(
      CampaignRepository repository,
      ScenarioCatalog catalog,
      CampaignEventPublisher publisher,
      ScenarioSocketHandler socket) {
    this.repository = repository;
    this.catalog = catalog;
    this.publisher = publisher;
    this.socket = socket;
  }

  public List<Campaign> findAll() {
    return repository.findAllByOrderByCreatedAtDesc();
  }

  public Optional<Campaign> findById(UUID id) {
    return repository.findById(id);
  }

  /**
   * Crée la campagne puis annonce {@code campaign.created} au Task Manager.
   *
   * <p>La méthode n'est volontairement pas transactionnelle : la ligne est validée en base avant la
   * publication, pour que le Task Manager la trouve. Si le broker est injoignable, la campagne ne
   * sera jamais exécutée : elle passe en FAILED.
   *
   * @param name libellé de la campagne ; le nom du modèle s'il est vide
   * @throws UnknownTemplateException si le modèle n'existe pas
   * @throws CampaignLaunchException si la publication a échoué
   */
  public Campaign launch(String templateId, String name, List<CampaignTarget> targets) {
    ScenarioTemplate template =
        catalog.findById(templateId).orElseThrow(() -> new UnknownTemplateException(templateId));
    String label = name == null || name.isBlank() ? template.name() : name.strip();

    Campaign campaign = repository.save(new Campaign(label, template.id(), targets));
    try {
      publisher.campaignCreated(campaign, template);
    } catch (AmqpException e) {
      log.error("Could not publish campaign.created for campaign {}", campaign.getId(), e);
      campaign.markFailed();
      socket.broadcast(repository.save(campaign));
      throw new CampaignLaunchException(e);
    }
    socket.broadcast(campaign);
    return campaign;
  }

  /** Le statut a changé en base (écrit par le Task Manager) : on le relit et on le diffuse. */
  public void onCampaignUpdated(UUID id) {
    repository
        .findById(id)
        .ifPresentOrElse(
            socket::broadcast, () -> log.warn("campaign.updated for unknown campaign {}", id));
  }
}
