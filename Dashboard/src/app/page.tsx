"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useBackendSocket } from "@/hooks/use-backend-socket";
import type { ScenarioStatus } from "@/lib/types";

const STATUS_LABEL: Record<ScenarioStatus, string> = {
  CREATED: "Créé",
  RUNNING: "En cours",
  SUCCESS: "Succès",
  FAILED: "Échec",
};

const STATUS_VARIANT: Record<
  ScenarioStatus,
  "default" | "secondary" | "destructive" | "outline"
> = {
  CREATED: "outline",
  RUNNING: "secondary",
  SUCCESS: "default",
  FAILED: "destructive",
};

function ConnectionBadge({
  status,
}: {
  status: "connecting" | "connected" | "disconnected";
}) {
  const label =
    status === "connected"
      ? "Back-end connecté"
      : status === "connecting"
        ? "Connexion en cours…"
        : "Back-end déconnecté";
  const variant =
    status === "connected"
      ? "default"
      : status === "connecting"
        ? "secondary"
        : "destructive";
  return <Badge variant={variant}>{label}</Badge>;
}

export default function Home() {
  const { scenarios, status, launchScenario, launchError } =
    useBackendSocket();
  const [launching, setLaunching] = useState(false);

  async function handleLaunch() {
    setLaunching(true);
    await launchScenario(`Scénario ${new Date().toLocaleString("fr-FR")}`);
    setLaunching(false);
  }

  return (
    <main className="mx-auto flex max-w-4xl flex-col gap-6 p-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">AutoStrike — Dashboard</h1>
        <ConnectionBadge status={status} />
      </div>

      <div className="flex items-center gap-3">
        <Button onClick={handleLaunch} disabled={launching}>
          {launching ? "Lancement…" : "Lancer un scénario"}
        </Button>
        {launchError && (
          <p className="text-sm text-destructive">{launchError}</p>
        )}
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Nom</TableHead>
            <TableHead>Statut</TableHead>
            <TableHead>Mis à jour</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {scenarios.length === 0 ? (
            <TableRow>
              <TableCell colSpan={3} className="text-center text-muted-foreground">
                Aucun scénario pour le moment.
              </TableCell>
            </TableRow>
          ) : (
            scenarios.map((scenario) => (
              <TableRow key={scenario.id}>
                <TableCell>{scenario.name}</TableCell>
                <TableCell>
                  <Badge variant={STATUS_VARIANT[scenario.status]}>
                    {STATUS_LABEL[scenario.status]}
                  </Badge>
                </TableCell>
                <TableCell>
                  {new Date(scenario.updatedAt).toLocaleString("fr-FR")}
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </main>
  );
}
