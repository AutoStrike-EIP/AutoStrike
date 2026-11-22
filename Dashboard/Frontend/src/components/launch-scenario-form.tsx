"use client";

import { useId, useState, type FormEvent } from "react";
import { Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import type {
  HostCriticality,
  LaunchScenarioRequest,
  ScenarioTemplate,
} from "@/lib/types";

/** Mêmes caractères autorisés que le Back-end : nom DNS ou adresse IP. */
const HOST_PATTERN = /^[A-Za-z0-9.:-]+$/;
const MAX_TARGETS = 50;

const FIELD_CLASS =
  "h-8 rounded-lg border border-input bg-transparent px-2.5 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 aria-invalid:border-destructive";

interface TargetRow {
  key: number;
  host: string;
  criticality: HostCriticality;
}

interface LaunchScenarioFormProps {
  templates: ScenarioTemplate[];
  templatesLoading: boolean;
  templatesError: boolean;
  launchError: string | null;
  onLaunch: (request: LaunchScenarioRequest) => Promise<boolean>;
}

let nextRowKey = 0;
function newRow(): TargetRow {
  nextRowKey += 1;
  return { key: nextRowKey, host: "", criticality: "TEST" };
}

export function LaunchScenarioForm({
  templates,
  templatesLoading,
  templatesError,
  launchError,
  onLaunch,
}: LaunchScenarioFormProps) {
  const ids = useId();
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [name, setName] = useState("");
  const [targets, setTargets] = useState<TargetRow[]>(() => [newRow()]);
  const [launching, setLaunching] = useState(false);

  const template =
    templates.find((t) => t.id === selectedId) ?? templates[0] ?? null;
  const hostsValid = targets.every((t) => HOST_PATTERN.test(t.host.trim()));
  const canLaunch = template !== null && hostsValid && !launching;

  function updateTarget(key: number, patch: Partial<TargetRow>) {
    setTargets((rows) =>
      rows.map((row) => (row.key === key ? { ...row, ...patch } : row)),
    );
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!canLaunch) return;
    setLaunching(true);
    const launched = await onLaunch({
      templateId: template.id,
      name: name.trim() || undefined,
      targets: targets.map(({ host, criticality }) => ({
        host: host.trim(),
        criticality,
      })),
    });
    setLaunching(false);
    if (launched) {
      setName("");
      setTargets([newRow()]);
    }
  }

  if (templatesLoading) {
    return (
      <p className="text-sm text-muted-foreground">
        Chargement des scénarios…
      </p>
    );
  }
  if (templatesError || template === null) {
    return (
      <p className="text-sm text-destructive">
        Impossible de charger les scénarios disponibles : le back-end est
        injoignable.
      </p>
    );
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="flex flex-col gap-4 rounded-xl border p-4"
    >
      <h2 className="text-lg font-medium">Lancer un scénario</h2>

      <div className="grid gap-4 sm:grid-cols-2">
        <label className="flex flex-col gap-1.5 text-sm" htmlFor={`${ids}-template`}>
          Scénario
          <select
            id={`${ids}-template`}
            className={FIELD_CLASS}
            value={template.id}
            onChange={(e) => setSelectedId(e.target.value)}
          >
            {templates.map((t) => (
              <option key={t.id} value={t.id}>
                {t.name}
              </option>
            ))}
          </select>
        </label>
        <label className="flex flex-col gap-1.5 text-sm" htmlFor={`${ids}-name`}>
          Nom (facultatif)
          <input
            id={`${ids}-name`}
            className={FIELD_CLASS}
            value={name}
            maxLength={255}
            placeholder={template.name}
            onChange={(e) => setName(e.target.value)}
          />
        </label>
      </div>

      <div className="flex flex-col gap-1.5 text-sm">
        <p className="text-muted-foreground">{template.description}</p>
        <ol className="flex flex-col gap-1">
          {template.steps.map((step) => (
            <li key={step.order}>
              <span className="font-mono">{step.technique}</span>{" "}
              {step.techniqueName}{" "}
              <span className="text-muted-foreground">({step.tactic})</span>
            </li>
          ))}
        </ol>
      </div>

      <fieldset className="flex flex-col gap-2">
        <legend className="mb-1.5 text-sm">Hôtes cibles</legend>
        {targets.map((row, index) => {
          const invalid =
            row.host.trim() !== "" && !HOST_PATTERN.test(row.host.trim());
          return (
            <div key={row.key} className="flex items-center gap-2">
              <input
                aria-label={`Hôte ${index + 1}`}
                aria-invalid={invalid}
                className={`${FIELD_CLASS} min-w-0 flex-1 font-mono`}
                value={row.host}
                maxLength={255}
                placeholder="10.0.0.12 ou srv-web.local"
                onChange={(e) => updateTarget(row.key, { host: e.target.value })}
              />
              <select
                aria-label={`Criticité de l'hôte ${index + 1}`}
                className={FIELD_CLASS}
                value={row.criticality}
                onChange={(e) =>
                  updateTarget(row.key, {
                    criticality: e.target.value as HostCriticality,
                  })
                }
              >
                <option value="PRODUCTION">Production</option>
                <option value="TEST">Test</option>
              </select>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label={`Retirer l'hôte ${index + 1}`}
                disabled={targets.length === 1}
                onClick={() =>
                  setTargets((rows) => rows.filter((r) => r.key !== row.key))
                }
              >
                <Trash2 />
              </Button>
            </div>
          );
        })}
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="self-start"
          disabled={targets.length >= MAX_TARGETS}
          onClick={() => setTargets((rows) => [...rows, newRow()])}
        >
          <Plus />
          Ajouter un hôte
        </Button>
      </fieldset>

      <div className="flex items-center gap-3">
        <Button type="submit" disabled={!canLaunch}>
          {launching ? "Lancement…" : "Lancer le scénario"}
        </Button>
        {launchError && (
          <p className="text-sm text-destructive" role="alert">
            {launchError}
          </p>
        )}
      </div>
    </form>
  );
}
