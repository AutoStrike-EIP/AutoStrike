# AutoStrike — Dashboard

Tout ce qui sert l'interface web d'AutoStrike, regroupé au même endroit :

| Dossier | Rôle | Stack |
|---|---|---|
| [`Frontend/`](Frontend/) | Interface de contrôle : lancement des campagnes, suivi en direct | Next.js + TypeScript |
| [`Backend/`](Backend/) | API REST + WebSocket, pont vers PostgreSQL et RabbitMQ | Spring Boot + Java |
| [`Database/`](Database/) | PostgreSQL local et documentation du schéma | PostgreSQL |

## Lancer l'ensemble en local

Depuis la racine du repo, tout se lance avec Docker (base, broker, back, front) :

```bash
docker compose up --build
```

Pour développer un seul composant sans Docker, voir le README de son dossier.
