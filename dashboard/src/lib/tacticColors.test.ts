import { describe, it, expect } from 'vitest';
import { tacticBadgeColors, getTacticBadgeColor, formatTacticName } from './tacticColors';

describe('tacticColors', () => {
  describe('tacticBadgeColors', () => {
    it('contains all 14 MITRE ATT&CK tactics', () => {
      const tactics = [
        'reconnaissance',
        'resource_development',
        'initial_access',
        'execution',
        'persistence',
        'privilege_escalation',
        'defense_evasion',
        'credential_access',
        'discovery',
        'lateral_movement',
        'collection',
        'command_and_control',
        'exfiltration',
        'impact',
      ];

      tactics.forEach((tactic) => {
        expect(tacticBadgeColors[tactic]).toBeDefined();
        expect(tacticBadgeColors[tactic]).toContain('bg-');
        expect(tacticBadgeColors[tactic]).toContain('text-');
      });
    });

    it('has correct color for reconnaissance', () => {
      expect(tacticBadgeColors.reconnaissance).toContain('bg-purple-100 text-purple-700');
      expect(tacticBadgeColors.reconnaissance).toContain('dark:bg-purple-900/30 dark:text-purple-400');
    });

    it('has correct color for discovery', () => {
      expect(tacticBadgeColors.discovery).toContain('bg-cyan-100 text-cyan-700');
      expect(tacticBadgeColors.discovery).toContain('dark:bg-cyan-900/30 dark:text-cyan-400');
    });

    it('has correct color for execution', () => {
      expect(tacticBadgeColors.execution).toContain('bg-orange-100 text-orange-700');
      expect(tacticBadgeColors.execution).toContain('dark:bg-orange-900/30 dark:text-orange-400');
    });
  });

  describe('getTacticBadgeColor', () => {
    it('returns correct color for known tactic', () => {
      expect(getTacticBadgeColor('discovery')).toContain('bg-cyan-100 text-cyan-700');
      expect(getTacticBadgeColor('discovery')).toContain('dark:');
    });

    it('returns correct color for tactic with underscores', () => {
      expect(getTacticBadgeColor('lateral_movement')).toContain('bg-teal-100 text-teal-700');
      expect(getTacticBadgeColor('lateral_movement')).toContain('dark:');
    });

    it('normalizes hyphens to underscores', () => {
      expect(getTacticBadgeColor('lateral-movement')).toContain('bg-teal-100 text-teal-700');
      expect(getTacticBadgeColor('lateral-movement')).toContain('dark:');
    });

    it('returns default gray for unknown tactic', () => {
      expect(getTacticBadgeColor('unknown_tactic')).toContain('bg-gray-100 text-gray-700');
      expect(getTacticBadgeColor('unknown_tactic')).toContain('dark:bg-gray-700 dark:text-gray-300');
    });

    it('returns default gray for empty string', () => {
      expect(getTacticBadgeColor('')).toContain('bg-gray-100 text-gray-700');
      expect(getTacticBadgeColor('')).toContain('dark:bg-gray-700 dark:text-gray-300');
    });

    it('handles mixed case input by not matching', () => {
      // tacticBadgeColors keys are lowercase, so mixed case won't match
      expect(getTacticBadgeColor('Discovery')).toContain('bg-gray-100 text-gray-700');
      expect(getTacticBadgeColor('Discovery')).toContain('dark:bg-gray-700 dark:text-gray-300');
    });
  });

  describe('formatTacticName', () => {
    it('replaces underscores with spaces', () => {
      expect(formatTacticName('lateral_movement')).toBe('lateral movement');
    });

    it('replaces hyphens with spaces', () => {
      expect(formatTacticName('lateral-movement')).toBe('lateral movement');
    });

    it('handles mixed hyphens and underscores', () => {
      expect(formatTacticName('command-and_control')).toBe('command and control');
    });

    it('returns same string if no separators', () => {
      expect(formatTacticName('discovery')).toBe('discovery');
    });

    it('handles empty string', () => {
      expect(formatTacticName('')).toBe('');
    });

    it('handles multiple consecutive separators', () => {
      expect(formatTacticName('test__name')).toBe('test  name');
    });
  });
});
