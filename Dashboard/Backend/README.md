# AutoStrike — Back-end

API REST + WebSocket entre le Dashboard, PostgreSQL et RabbitMQ (ticket #14).
Spring Boot 4 + Java 21. Voir `docs/ARCHITECTURE.md` pour la vue d'ensemble.

## Prérequis

- JDK 21
- PostgreSQL (`../Database`, `docker compose up -d`) et RabbitMQ (`Broker/` à la racine du repo)

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
| `GET /scenario-templates` | Modèles de scénarios disponibles et leur chaîne de techniques MITRE |
| `GET /scenarios` | Liste des scénarios lancés, du plus récent au plus ancien |
| `GET /scenarios/{id}` | Un scénario (404 s'il n'existe pas) |
| `POST /scenarios` | Lance un scénario, réponse 201 (corps ci-dessous) |
| `WS /scenarios` | Pousse `{"type": "scenario.updated", "payload": <scénario>}` à chaque création ou changement de statut |

Corps de `POST /scenarios` :

```json
{
  "templateId": "linux-basic-v1",
  "name": "Audit poste RH",
  "targets": [
    { "host": "10.0.0.12", "criticality": "PRODUCTION" },
    { "host": "lab-web.local", "criticality": "TEST" }
  ]
}
```

- `name` est facultatif : le nom du modèle est utilisé s'il est absent.
- 1 à 50 cibles ; `host` limité aux caractères d'un nom DNS ou d'une IP (`A-Z a-z 0-9 . : -`).
- `criticality` : `PRODUCTION` ou `TEST`. Elle départage les résultats à profondeur égale dans
  la priorisation (filtre à étages, étage 2).

Un scénario : `{id, name, templateId, targets, status, createdAt, updatedAt}`, `status` parmi
`CREATED`, `RUNNING`, `SUCCESS`, `FAILED`, dates en ISO-8601 UTC.

`POST /scenarios` répond 400 si la requête est invalide (modèle inconnu, aucune cible, hôte ou
criticité invalide) et 503 si RabbitMQ est injoignable (le scénario est alors enregistré en
`FAILED`).

Un seul modèle existe pour l'instant, `linux-basic-v1` : le scénario v1 du livrable « Scénario
de campagne v1 (Linux, basique) » (T1087.001 → T1552.001 → T1053.003).

## RabbitMQ

- Exchange **topic** durable `autostrike.events` (variable `RABBITMQ_EXCHANGE`), déclaré
  à l'identique par le serveur Rust (`Server/common/src/event`). Un type différent d'un côté
  ferait échouer la déclaration de l'autre (`PRECONDITION_FAILED`).
- Publie `campaign.created` (JSON persistant) avec le schéma de campagne du livrable, en
  snake_case :

  ```json
  {
    "campaign": {
      "id": "…", "name": "Audit poste RH", "template": "linux-basic-v1", "mode": "non-destructive",
      "target_hosts": [{ "host": "10.0.0.12", "criticality": "PRODUCTION" }],
      "tasks": [
        { "order": 1, "technique": "T1087.001", "action": "list_local_accounts" },
        { "order": 2, "technique": "T1552.001", "action": "grep_credentials_in_files", "depends_on": 1 },
        { "order": 3, "technique": "T1053.003", "action": "cron_persistence_test", "depends_on": 2, "cleanup": true }
      ]
    }
  }
  ```

  Écart au livrable : `target_hosts` contient des objets `{host, criticality}` et non de simples
  chaînes, pour que la criticité arrive jusqu'à la priorisation.
- Lit la queue durable `backend_notifications`, liée à `campaign.updated` et `task.updated`.
  Sur `campaign.updated`, seul le champ `id` du message est lu : le Back-end relit la campagne
  en base et la diffuse sur le WebSocket. `task.updated` est ignoré pour l'instant.

Le format des messages reste provisoire (ticket #20).

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
