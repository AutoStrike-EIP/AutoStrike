# AutoStrike — Dashboard

Tout ce qui sert l'interface web d'AutoStrike, regroupé au même endroit :

| Dossier | Rôle | Stack |
|---|---|---|
| [`Frontend/`](Frontend/) | Interface de contrôle : lancement des campagnes, suivi en direct | Next.js + TypeScript |
| [`Backend/`](Backend/) | API REST + WebSocket, pont vers PostgreSQL et RabbitMQ | Spring Boot + Java |
| [`Database/`](Database/) | PostgreSQL local et documentation du schéma | PostgreSQL |

## Lancer l'ensemble en local

```bash
cd Database && docker compose up -d            # PostgreSQL
cd ../../Broker && docker compose up -d        # RabbitMQ (branche serveur, PR #23)
cd ../Dashboard/Backend && ./mvnw spring-boot:run   # API sur http://localhost:3001
cd ../Frontend && npm install && npm run dev        # UI sur http://localhost:3000
```

Chaque dossier a son propre README avec ses commandes de build, de test et de lint.
