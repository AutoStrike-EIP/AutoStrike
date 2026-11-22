# AutoStrike — Stack technique & Architecture

> Document vivant : à mettre à jour à chaque évolution de l'architecture (nouveau composant, changement de techno, nouvelle décision structurante). Dernière mise à jour : **12 novembre 2026**.

## 1. Stack technique

| Composant | Techno | Rôle |
|---|---|---|
| Dashboard | Next.js + TypeScript | Interface de contrôle : lancement des campagnes, visualisation des rapports et des agents connectés |
| Back-end | Spring Boot + Java | Intermédiaire entre le Dashboard et la base de données ; expose l'API REST et le WebSocket |
| Task Manager | Rust | Détecte les nouvelles campagnes et crée les tâches correspondantes ; met à jour le statut des campagnes/tâches en base |
| Task Transmitter | Rust | Transmet les tâches une à une aux agents et relaie leurs résultats |
| Agent | Rust | Exécute les tâches reçues et renvoie un rapport |
| Message broker | RabbitMQ | Messagerie inter-services (événements `task.*` et `campaign.*`) |
| Base de données | PostgreSQL | Stockage permanent (utilisateurs, campagnes, tâches, rapports) |

**Task Manager** et **Task Transmitter** forment ensemble le composant "Serveur" côté orchestration.

*Redis a été retiré de la stack (cf. décision ADR-001 ci-dessous).*

*Le Back-end est passé de NestJS/TypeScript à Spring Boot/Java (cf. décision ADR-002 ci-dessous).*

## 2. Composants et rôles

- **Dashboard** — interface web pour lancer des campagnes de test et consulter les résultats.
- **Back-end** — seul composant à parler directement à PostgreSQL pour les utilisateurs et les campagnes (pas les tâches) ; expose REST (auth, création/lancement de campagnes) et WebSocket (statuts en direct vers le Dashboard).
- **Task Manager** — à la création d'une campagne, génère les tâches associées et les écrit en base ; tient PostgreSQL à jour sur le statut des campagnes/tâches.
- **Task Transmitter** — reçoit les tâches créées et les distribue aux agents connectés ; renvoie leurs résultats au reste du système.
- **Agent** — exécute une simulation et produit un rapport.
- **RabbitMQ** — fait transiter les événements entre Back-end, Task Manager et Task Transmitter (fan-out fiable, avec persistance).
- **PostgreSQL** — stockage permanent de toutes les données.

## 3. Flux de communication

```
Dashboard <--WebSocket (statuts live)--> Back-end <--REST (auth, create/run)--> Dashboard
Back-end <--> PostgreSQL        (users, campagnes — PAS les tâches)
Back-end <--> RabbitMQ          (écrit campaign.created / lit backend_notifications)
RabbitMQ <--> Task Manager      (lit manager_notifications / écrit task.created)
Task Manager <--> PostgreSQL    (crée les tâches, met à jour campagnes/tâches)
RabbitMQ <--> Task Transmitter  (lit transmitter_notifications / écrit task.updated)
Task Transmitter <--> Agent     (envoie les tâches, reçoit les résultats)
```

## 4. Événements RabbitMQ

| Événement | Publié par | Consommé par | Déclencheur |
|---|---|---|---|
| `campaign.created` | Back-end | Task Manager | Une nouvelle campagne est créée côté Dashboard/API |
| `campaign.updated` | Task Manager | Back-end | Le statut d'une campagne change en base |
| `task.created` | Task Manager | Task Transmitter | Une tâche vient d'être générée pour une campagne |
| `task.updated` | Task Transmitter | Back-end, Task Manager | Un agent a renvoyé un résultat/statut pour une tâche |
| `task.deleted` / `campaign.deleted` | — | — | Prévus, pas encore implémentés |

## 5. Queues RabbitMQ

| Queue | Propriétaire | Abonnée à |
|---|---|---|
| `backend_notifications` | Back-end | `campaign.updated`, `task.updated` |
| `manager_notifications` | Task Manager | `campaign.created`, `task.updated` |
| `transmitter_notifications` | Task Transmitter | `task.created` |

Chaque queue est durable et propre à son service : si un service redémarre, ses messages en attente restent disponibles au lieu d'être perdus.

## 6. Points d'attention sur l'architecture RabbitMQ (à traiter)

- **Producteur manquant pour `campaign.updated`** — `backend_notifications` y est abonné, mais aucun composant ne le publie actuellement. Task Manager doit publier `campaign.updated` quand le statut global d'une campagne change.
- **Type d'exchange non précisé** — un exchange direct suffit vu qu'il n'y a pas de besoin de wildcard ; un topic exchange (`task.#`) peut être envisagé si des sous-types d'événements apparaissent plus tard.
- **Durabilité = une question de configuration, pas de choix de techno** — sans queues `durable:true`, messages `persistent` (delivery_mode=2) et ack manuel après écriture en base, RabbitMQ perd des messages exactement comme Redis Pub/Sub.
- **Pas de dead-letter queue définie** — sans DLX + limite de retry, risque de boucle infinie ou de perte silencieuse selon la config de requeue.
- **Idempotence et ordre non traités** — RabbitMQ est at-least-once (doublons possibles) et ne garantit pas l'ordre entre queues différentes.
- **Scalabilité du Back-end à anticiper** — si le Back-end est scalé horizontalement, une `backend_notifications` unique ne livrera chaque événement qu'à une seule instance ; il faudra une queue par instance (exclusive/auto-delete).

## 7. Journal des décisions techniques (ADR)

### ADR-001 — Message broker inter-services : RabbitMQ (remplace Redis Pub/Sub)

- **Date** : 11 septembre 2026
- **Contexte** : l'architecture initiale (juillet 2026) prévoyait Redis Pub/Sub comme bus d'événements. Plusieurs événements (ex. `task.updated`) doivent atteindre plusieurs consommateurs à la fois, et leur perte entraîne une désynchronisation silencieuse entre l'état réel d'une tâche et son statut en base.
- **Options considérées** : Redis Pub/Sub (pas de persistance) ; Redis Streams + consumer groups (jamais utilisé par l'équipe) ; RabbitMQ (queue durable par service, ack manuel).
- **Choix** : RabbitMQ.
- **Raison** : seule option garantissant qu'un événement métier n'est jamais perdu même si un service tombe ou redémarre.
- **Conséquences acceptées** : Redis retiré de la stack ; RabbitMQ devient un point de défaillance unique à documenter ; un service supplémentaire à maintenir.

### ADR-002 — Back-end : Spring Boot + Java (remplace NestJS + TypeScript)

- **Date** : 12 novembre 2026
- **Contexte** : la stack initiale prévoyait NestJS + TypeScript pour rester mono-langage avec le Dashboard. Aucun développement significatif du Back-end n'avait encore démarré au moment du changement.
- **Choix** : Spring Boot + Java.
- **Raison** : préférence et expérience de l'équipe — plus à l'aise en Java/Spring qu'en TypeScript/NestJS.
- **Conséquences acceptées** : le projet n'est plus mono-langage TypeScript/Rust ; Java s'ajoute comme quatrième langage (à anticiper pour le dimensionnement Docker/CI) ; le contrat REST/WebSocket côté Dashboard reste HTTP/WebSocket standard, non affecté dans son principe.

---

*Prochaines mises à jour attendues : détail du format des messages (schéma JSON par événement), authentification/permissions RabbitMQ, stratégie de dead-letter queue, mise à jour du dimensionnement Proxmox, et spécification précise des routes REST/WebSocket du Back-end Spring Boot.*
