# AutoStrike — Frontend

Interface web d'AutoStrike (Next.js + TypeScript + Tailwind + shadcn/ui).

## Commandes

```bash
npm install     # registre npm via le miroir de .npmrc
npm run dev     # http://localhost:3000
npm run lint    # linter — à lancer avant toute PR
npm run build
```

Copier `.env.local.example` en `.env.local` pour pointer vers un autre Back-end que
`http://localhost:3001`.

Pas de `next/font/google` ni d'asset Google au build (bloqués sur certains réseaux) :
police système uniquement.
