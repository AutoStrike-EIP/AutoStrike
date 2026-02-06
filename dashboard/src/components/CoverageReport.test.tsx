import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { CoverageReport } from './CoverageReport';

describe('CoverageReport', () => {
  const mockCoverage = {
    discovery: 9,
    execution: 3,
    persistence: 4,
    defense_evasion: 3,
    credential_access: 4,
    collection: 4,
    lateral_movement: 3,
  };

  it('renders total technique count', () => {
    render(<CoverageReport coverage={mockCoverage} totalTechniques={30} />);

    expect(screen.getByText('30')).toBeInTheDocument();
    expect(screen.getByText(/techniques/)).toBeInTheDocument();
  });

  it('renders all provided tactics', () => {
    render(<CoverageReport coverage={mockCoverage} totalTechniques={30} />);

    expect(screen.getByText('discovery')).toBeInTheDocument();
    expect(screen.getByText('execution')).toBeInTheDocument();
    expect(screen.getByText('persistence')).toBeInTheDocument();
    expect(screen.getByText('defense evasion')).toBeInTheDocument();
    expect(screen.getByText('credential access')).toBeInTheDocument();
    expect(screen.getByText('collection')).toBeInTheDocument();
    expect(screen.getByText('lateral movement')).toBeInTheDocument();
  });

  it('displays correct counts for each tactic', () => {
    render(<CoverageReport coverage={mockCoverage} totalTechniques={30} />);

    expect(screen.getByText('9')).toBeInTheDocument();
    // Multiple tactics have count 3 (execution, defense_evasion, lateral_movement)
    expect(screen.getAllByText('3').length).toBe(3);
    // Multiple tactics have count 4 (persistence, credential_access, collection)
    expect(screen.getAllByText('4').length).toBe(3);
  });

  it('handles empty coverage data', () => {
    render(<CoverageReport coverage={{}} totalTechniques={0} />);

    expect(screen.getByText('No coverage data available')).toBeInTheDocument();
  });

  it('supports compact variant (default)', () => {
    const { container } = render(
      <CoverageReport coverage={mockCoverage} totalTechniques={30} variant="compact" />
    );

    // Compact variant should have progress bars
    expect(container.querySelectorAll('.h-2').length).toBeGreaterThan(0);
  });

  it('supports detailed variant', () => {
    render(
      <CoverageReport coverage={mockCoverage} totalTechniques={30} variant="detailed" />
    );

    // Detailed variant shows percentages for each tactic
    expect(screen.getAllByText(/%/).length).toBeGreaterThan(0);
    // Multiple elements contain "technique" (header and cards)
    expect(screen.getAllByText(/technique/).length).toBeGreaterThan(0);
  });

  it('formats tactic names correctly (replaces underscores with spaces)', () => {
    const coverage = {
      defense_evasion: 3,
      lateral_movement: 2,
    };

    render(<CoverageReport coverage={coverage} totalTechniques={5} />);

    expect(screen.getByText('defense evasion')).toBeInTheDocument();
    expect(screen.getByText('lateral movement')).toBeInTheDocument();
  });

  it('sorts tactics by count descending', () => {
    const coverage = {
      low_count: 1,
      high_count: 10,
      mid_count: 5,
    };

    render(<CoverageReport coverage={coverage} totalTechniques={16} />);

    const tacticNames = screen.getAllByText(/count/);
    expect(tacticNames[0]).toHaveTextContent('high count');
    expect(tacticNames[1]).toHaveTextContent('mid count');
    expect(tacticNames[2]).toHaveTextContent('low count');
  });

  it('applies custom className', () => {
    const { container } = render(
      <CoverageReport coverage={mockCoverage} totalTechniques={30} className="custom-class" />
    );

    expect(container.firstChild).toHaveClass('custom-class');
  });

  it('renders progress bars with correct colors', () => {
    const coverage = { discovery: 5 };
    const { container } = render(<CoverageReport coverage={coverage} totalTechniques={5} />);

    // Discovery should have cyan color
    expect(container.querySelector('.bg-cyan-500')).toBeInTheDocument();
  });

  it('displays tactic count in detailed variant', () => {
    const coverage = { discovery: 9 };

    render(<CoverageReport coverage={coverage} totalTechniques={9} variant="detailed" />);

    expect(screen.getByText('9 techniques')).toBeInTheDocument();
  });

  it('handles singular technique correctly in detailed variant', () => {
    const coverage = { discovery: 1 };

    render(<CoverageReport coverage={coverage} totalTechniques={1} variant="detailed" />);

    expect(screen.getByText('1 technique')).toBeInTheDocument();
  });

  it('uses fallback color for unknown tactic', () => {
    const coverage = { unknown_tactic: 3 };
    const { container } = render(<CoverageReport coverage={coverage} totalTechniques={3} />);

    // Unknown tactic should use gray fallback
    expect(container.querySelector('.bg-gray-500')).toBeInTheDocument();
  });

  it('handles detailed variant with zero total correctly', () => {
    // total (sum of coverage values) is 0 when coverage has entries with 0
    const coverage = { discovery: 0 };
    render(<CoverageReport coverage={coverage} totalTechniques={0} variant="detailed" />);

    // Percentage should be 0.0% when total is 0
    expect(screen.getByText('0.0%')).toBeInTheDocument();
  });

  it('applies custom className to empty state', () => {
    const { container } = render(
      <CoverageReport coverage={{}} totalTechniques={0} className="custom-empty" />
    );

    expect(container.firstChild).toHaveClass('custom-empty');
  });

  it('uses fallback color for unknown tactic in detailed variant', () => {
    const coverage = { unknown_tactic: 5 };
    const { container } = render(
      <CoverageReport coverage={coverage} totalTechniques={5} variant="detailed" />
    );

    // Unknown tactic should use gray fallback in the detailed card
    expect(container.querySelector('.bg-gray-500')).toBeInTheDocument();
    // Should show percentage
    expect(screen.getByText('100.0%')).toBeInTheDocument();
    expect(screen.getByText('5 techniques')).toBeInTheDocument();
  });

  it('handles hyphenated tactic names through normalization', () => {
    const coverage = { 'defense-evasion': 3 };
    const { container } = render(
      <CoverageReport coverage={coverage} totalTechniques={3} />
    );

    // Hyphenated tactic should be normalized to underscore and matched to green color
    expect(container.querySelector('.bg-green-500')).toBeInTheDocument();
    expect(screen.getByText('defense evasion')).toBeInTheDocument();
  });
});
