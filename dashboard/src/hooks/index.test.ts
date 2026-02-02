import { describe, it, expect } from 'vitest';
import { useWebSocket } from './index';
import type { WebSocketMessage } from './index';

describe('hooks/index exports', () => {
  it('exports useWebSocket hook', () => {
    expect(useWebSocket).toBeDefined();
    expect(typeof useWebSocket).toBe('function');
  });

  it('exports WebSocketMessage type correctly', () => {
    // Type check - this validates the type export
    const message: WebSocketMessage = {
      type: 'test',
      payload: { data: 'test' },
    };
    expect(message.type).toBe('test');
  });
});
