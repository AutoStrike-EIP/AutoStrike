package com.autostrike.backend.campaign;

import java.util.List;
import java.util.UUID;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;

public interface CampaignRepository extends JpaRepository<Campaign, UUID> {

  /** Charge les cibles dans la même requête, plutôt qu'une requête par campagne. */
  @EntityGraph(attributePaths = "targets")
  List<Campaign> findAllByOrderByCreatedAtDesc();
}
