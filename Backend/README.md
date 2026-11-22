# AutoStrike — Back-end

API REST + WebSocket entre le Dashboard, PostgreSQL et RabbitMQ (ticket #14).
Spring Boot 4 + Java 21. Voir `docs/ARCHITECTURE.md` pour la vue d'ensemble.

## Prérequis

- JDK 21
- PostgreSQL et RabbitMQ (`docker compose up -d` démarre les deux en local)

Maven n'a pas besoin d'être installé : le wrapper `./mvnw` s'en charge.

## Commandes

```bash
./mvnw spring-boot:run   # lance le serveur sur http://localhost:3001
./mvnw verify            # compile + tests (ni PostgreSQL ni RabbitMQ requis)
./mvnw spotless:check    # linter — à lancer avant toute PR
./mvnw spotless:apply    # corrige le formatage
```

### Réseau lent vers Maven Central (Chine continentale)

Aucun miroir n'est imposé par défaut. Si Maven Central est lent ou coupe les gros
téléchargements (`Connection reset`), passer par le miroir fourni :

```bash
./mvnw -s .mvn/settings-cn.xml verify
```

Pour ne pas répéter `-s` à chaque commande, copier ce fichier en `~/.m2/settings.xml`
(réglage propre à la machine, non versionné).

## Contrat exposé au Dashboard

| Route | Rôle |
|---|---|
| `GET /scenarios` | Liste des scénarios, du plus récent au plus ancien |
| `GET /scenarios/{id}` | Un scénario (404 s'il n'existe pas) |
| `POST /scenarios` | Lance un scénario — corps `{"name": "..."}`, réponse 201 |
| `WS /scenarios` | Pousse `{"type": "scenario.updated", "payload": <scénario>}` à chaque création ou changement de statut |

Un scénario : `{id, name, status, createdAt, updatedAt}`, `status` parmi
`CREATED`, `RUNNING`, `SUCCESS`, `FAILED`, dates en ISO-8601 UTC.

`POST /scenarios` répond 400 si `name` est vide et 503 si RabbitMQ est injoignable
(le scénario est alors enregistré en `FAILED`).

## RabbitMQ

- Exchange direct durable `autostrike.events` (variable `RABBITMQ_EXCHANGE`).
- Publie `campaign.created` (JSON persistant, même forme que le scénario ci-dessus).
- Lit la queue durable `backend_notifications`, liée à `campaign.updated` et `task.updated`.
  Sur `campaign.updated`, seul le champ `id` du message est lu : le Back-end relit la campagne
  en base et la diffuse sur le WebSocket. `task.updated` est ignoré pour l'instant.

Le nom de l'exchange et le format des messages sont provisoires (ticket #20).

## Configuration

| Variable | Défaut |
|---|---|
| `SERVER_PORT` | `3001` |
| `DATABASE_URL` | `jdbc:postgresql://localhost:5432/autostrike` |
| `DATABASE_USERNAME` / `DATABASE_PASSWORD` | `autostrike` / `autostrike` |
| `RABBITMQ_HOST` / `RABBITMQ_PORT` | `localhost` / `5672` |
| `RABBITMQ_USERNAME` / `RABBITMQ_PASSWORD` | `guest` / `guest` |
| `RABBITMQ_EXCHANGE` | `autostrike.events` |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:3000` |

Les valeurs par défaut sont des identifiants de développement, à surcharger ailleurs.

Le schéma de la base est géré par Flyway (`src/main/resources/db/migration`).
