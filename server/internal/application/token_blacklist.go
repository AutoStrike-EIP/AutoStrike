package application

import (
	"sync"
	"time"
)

// TokenBlacklist maintains a set of revoked JWT tokens in memory.
// Tokens are automatically cleaned up after they would have expired.
type TokenBlacklist struct {
	mu     sync.RWMutex
	tokens map[string]time.Time // token -> expiry time
}

// NewTokenBlacklist creates a new token blacklist with background cleanup.
func NewTokenBlacklist() *TokenBlacklist {
	bl := &TokenBlacklist{
		tokens: make(map[string]time.Time),
	}
	go bl.cleanupLoop()
	return bl
}

// Revoke adds a token to the blacklist until its expiry time.
func (bl *TokenBlacklist) Revoke(token string, expiresAt time.Time) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	bl.tokens[token] = expiresAt
}

// IsRevoked checks if a token has been revoked and is still within its expiry window.
func (bl *TokenBlacklist) IsRevoked(token string) bool {
	bl.mu.RLock()
	defer bl.mu.RUnlock()
	expiresAt, exists := bl.tokens[token]
	if !exists {
		return false
	}
	return time.Now().Before(expiresAt)
}

// cleanupLoop removes expired tokens every minute to prevent memory growth.
func (bl *TokenBlacklist) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		bl.cleanup()
	}
}

func (bl *TokenBlacklist) cleanup() {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	now := time.Now()
	for token, expiresAt := range bl.tokens {
		if now.After(expiresAt) {
			delete(bl.tokens, token)
		}
	}
}
