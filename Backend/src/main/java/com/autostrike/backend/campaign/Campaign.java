package com.autostrike.backend.campaign;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.PrePersist;
import jakarta.persistence.PreUpdate;
import jakarta.persistence.Table;
import java.time.Instant;
import java.util.UUID;

/**
 * Une campagne de test. Le Dashboard parle de "scénario" ; la base et RabbitMQ parlent de
 * "campagne" (docs/ARCHITECTURE.md). Le Back-end la crée, le Task Manager fait évoluer son statut.
 */
@Entity
@Table(name = "campaigns")
public class Campaign {

  @Id
  @GeneratedValue(strategy = GenerationType.UUID)
  private UUID id;

  @Column(nullable = false)
  private String name;

  @Enumerated(EnumType.STRING)
  @Column(nullable = false, length = 16)
  private CampaignStatus status;

  @Column(name = "created_at", nullable = false, updatable = false)
  private Instant createdAt;

  @Column(name = "updated_at", nullable = false)
  private Instant updatedAt;

  protected Campaign() {}

  public Campaign(String name) {
    this.name = name;
    this.status = CampaignStatus.CREATED;
  }

  @PrePersist
  void onCreate() {
    createdAt = Instant.now();
    updatedAt = createdAt;
  }

  @PreUpdate
  void onUpdate() {
    updatedAt = Instant.now();
  }

  public void markFailed() {
    status = CampaignStatus.FAILED;
  }

  public UUID getId() {
    return id;
  }

  public String getName() {
    return name;
  }

  public CampaignStatus getStatus() {
    return status;
  }

  public Instant getCreatedAt() {
    return createdAt;
  }

  public Instant getUpdatedAt() {
    return updatedAt;
  }
}
