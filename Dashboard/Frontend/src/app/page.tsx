"use client";

import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { LaunchScenarioForm } from "@/components/launch-scenario-form";
import { useBackendSocket } from "@/hooks/use-backend-socket";
import { useScenarioTemplates } from "@/hooks/use-scenario-templates";
import type { ScenarioResult, ScenarioStatus } from "@/lib/types";

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

/** Ex. « 3 hôtes, dont 1 en production ». */
function describeTargets({ targets }: ScenarioResult): string {
  const production = targets.filter((t) => t.criticality === "PRODUCTION");
  const hosts = targets.length === 1 ? "1 hôte" : `${targets.length} hôtes`;
  return production.length === 0
    ? hosts
    : `${hosts}, dont ${production.length} en production`;
}

export default function Home() {
  const { scenarios, status, launchScenario, launchError } =
    useBackendSocket();
  const templates = useScenarioTemplates();

  return (
    <main className="mx-auto flex max-w-4xl flex-col gap-6 p-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">AutoStrike — Dashboard</h1>
        <ConnectionBadge status={status} />
      </div>

      <LaunchScenarioForm
        templates={templates.templates}
        templatesLoading={templates.loading}
        templatesError={templates.error}
        launchError={launchError}
        onLaunch={launchScenario}
      />

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Nom</TableHead>
            <TableHead>Cibles</TableHead>
            <TableHead>Statut</TableHead>
            <TableHead>Mis à jour</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {scenarios.length === 0 ? (
            <TableRow>
              <TableCell colSpan={4} className="text-center text-muted-foreground">
                Aucun scénario pour le moment.
              </TableCell>
            </TableRow>
          ) : (
            scenarios.map((scenario) => (
              <TableRow key={scenario.id}>
                <TableCell>{scenario.name}</TableCell>
                <TableCell>{describeTargets(scenario)}</TableCell>
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
