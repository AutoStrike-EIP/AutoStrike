# AutoStrike — Base de données

PostgreSQL : stockage permanent des utilisateurs, campagnes, tâches et rapports
(voir `docs/ARCHITECTURE.md`).

## Lancer la base en local

```bash
docker compose up -d
```

Base `autostrike` sur `localhost:5432`, utilisateur et mot de passe `autostrike`
(identifiants de développement uniquement).

## Schéma

Le schéma est versionné par des migrations Flyway, appliquées au démarrage du Back-end :
`../Backend/src/main/resources/db/migration`. Toute évolution du schéma = une nouvelle
migration `V<n>__<description>.sql`, jamais la modification d'une migration existante.

| Table | Écrite par | Contenu |
|---|---|---|
| `campaigns` | Back-end (création), Task Manager (statut) | Une campagne lancée depuis le Dashboard |
| `campaign_targets` | Back-end | Hôtes ciblés par une campagne et leur criticité métier |

Les tables des tâches et des rapports appartiennent au Task Manager et restent à définir
(tickets #15 et #20).

RabbitMQ n'est pas ici : il se lance depuis `Broker/` (branche serveur, PR #23).
