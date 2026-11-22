-- Une campagne applique un modèle de scénario à une liste d'hôtes cibles.
-- Les campagnes créées avant cette migration n'avaient qu'un modèle possible : linux-basic-v1.
ALTER TABLE campaigns ADD COLUMN template_id VARCHAR(64) DEFAULT 'linux-basic-v1' NOT NULL;
ALTER TABLE campaigns ALTER COLUMN template_id DROP DEFAULT;

CREATE TABLE campaign_targets (
    campaign_id  UUID         NOT NULL REFERENCES campaigns (id) ON DELETE CASCADE,
    target_order INTEGER      NOT NULL,
    host         VARCHAR(255) NOT NULL,
    criticality  VARCHAR(16)  NOT NULL,
    PRIMARY KEY (campaign_id, target_order)
);
