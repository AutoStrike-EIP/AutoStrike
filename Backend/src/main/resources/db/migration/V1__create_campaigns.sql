-- Schéma provisoire : le schéma officiel sera défini par les tickets #15 et #20.
CREATE TABLE campaigns (
    id         UUID PRIMARY KEY,
    name       VARCHAR(255)             NOT NULL,
    status     VARCHAR(16)              NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);
