/**
 * Color mapping for MITRE ATT&CK tactics.
 * Shared across components for consistent tactic styling.
 */

/**
 * Badge colors for tactic labels (used in tables and lists).
 */
export const tacticBadgeColors: Record<string, string> = {
  reconnaissance: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
  resource_development: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
  initial_access: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
  execution: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
  persistence: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
  privilege_escalation: 'bg-pink-100 text-pink-700 dark:bg-pink-900/30 dark:text-pink-400',
  defense_evasion: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
  credential_access: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400',
  discovery: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-900/30 dark:text-cyan-400',
  lateral_movement: 'bg-teal-100 text-teal-700 dark:bg-teal-900/30 dark:text-teal-400',
  collection: 'bg-lime-100 text-lime-700 dark:bg-lime-900/30 dark:text-lime-400',
  command_and_control: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400',
  exfiltration: 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-400',
  impact: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
};

/**
 * Get the badge color class for a tactic.
 * Normalizes the tactic string (replaces hyphens with underscores).
 */
export function getTacticBadgeColor(tactic: string): string {
  const normalizedTactic = String(tactic).replaceAll('-', '_');
  return tacticBadgeColors[normalizedTactic] || 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300';
}

/**
 * Format tactic name for display.
 * Replaces hyphens and underscores with spaces.
 */
export function formatTacticName(tactic: string): string {
  return String(tactic).replaceAll('-', '_').replaceAll('_', ' ');
}

/**
 * Progress bar colors for tactic coverage (used in CoverageReport).
 */
export const tacticBarColors: Record<string, string> = {
  reconnaissance: 'bg-purple-500',
  resource_development: 'bg-blue-500',
  initial_access: 'bg-red-500',
  execution: 'bg-orange-500',
  persistence: 'bg-yellow-500',
  privilege_escalation: 'bg-pink-500',
  defense_evasion: 'bg-green-500',
  credential_access: 'bg-indigo-500',
  discovery: 'bg-cyan-500',
  lateral_movement: 'bg-teal-500',
  collection: 'bg-lime-500',
  command_and_control: 'bg-amber-500',
  exfiltration: 'bg-rose-500',
  impact: 'bg-red-600',
};

/**
 * Get the progress bar color class for a tactic.
 * Normalizes the tactic string (replaces hyphens with underscores).
 */
export function getTacticBarColor(tactic: string): string {
  const normalizedTactic = String(tactic).replaceAll('-', '_');
  return tacticBarColors[normalizedTactic] || 'bg-gray-500';
}
