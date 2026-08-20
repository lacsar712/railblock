package conflict

import (
	"sync"
	"time"
)

// AuditLog retains recent conflict-related events for operator review.
type AuditLog struct {
	mu      sync.RWMutex
	entries []AuditEntry
	max     int
}

// AuditEntry describes a recorded conflict event.
type AuditEntry struct {
	Timestamp int64  `json:"timestamp"`
	BlockID   uint16 `json:"block_id"`
	Kind      Kind   `json:"kind"`
	Message   string `json:"message"`
	Source    string `json:"source,omitempty"`
}

// NewAuditLog creates an audit log with a maximum entry count.
func NewAuditLog(max int) *AuditLog {
	if max <= 0 {
		max = 128
	}
	return &AuditLog{max: max}
}

// Record appends an audit entry with the current unix timestamp.
func (a *AuditLog) Record(blockID uint16, kind Kind, message, source string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, AuditEntry{
		Timestamp: time.Now().Unix(),
		BlockID:   blockID,
		Kind:      kind,
		Message:   message,
		Source:    source,
	})
	if len(a.entries) > a.max {
		a.entries = a.entries[len(a.entries)-a.max:]
	}
}

// Entries returns a copy of recent audit records.
func (a *AuditLog) Entries() []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]AuditEntry, len(a.entries))
	copy(out, a.entries)
	return out
}

// Clear removes all audit entries.
func (a *AuditLog) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = nil
}

// Count returns the number of stored entries.
func (a *AuditLog) Count() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}
