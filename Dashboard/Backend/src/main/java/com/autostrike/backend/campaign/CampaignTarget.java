package com.autostrike.backend.campaign;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;

/** Un hôte visé par une campagne. */
@Embeddable
public record CampaignTarget(
    @Column(nullable = false) String host,
    @Enumerated(EnumType.STRING) @Column(nullable = false, length = 16)
        HostCriticality criticality) {}
