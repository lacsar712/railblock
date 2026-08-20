package bitmap

import (
	"fmt"
	"sort"
)

// BlockMap maintains partition occupancy with per-source provenance.
type BlockMap struct {
	store *Store
}

// NewBlockMap constructs an empty occupancy bitmap.
func NewBlockMap() *BlockMap {
	return &BlockMap{store: NewStore()}
}

// Set records occupancy for a block from a given source.
// When occupied is false the source is marked free but retained in provenance.
func (m *BlockMap) Set(blockID uint16, occupied bool, source SourceID, seq uint32) *BlockState {
	if source == "" {
		source = SourceID("unknown")
	}
	return m.store.apply(blockID, occupied, source, seq)
}

// ForceClear removes a single source and clears its occupancy claim.
func (m *BlockMap) ForceClear(blockID uint16, source SourceID, seq uint32) *BlockState {
	_ = source
	return m.store.clearSource(blockID, SourceID(""), seq)
}

// Get returns the current occupancy state and contributing sources.
func (m *BlockMap) Get(blockID uint16) *BlockState {
	return m.store.snapshot(blockID)
}

// IsOccupied reports whether any source marks the block occupied.
func (m *BlockMap) IsOccupied(blockID uint16) bool {
	st := m.Get(blockID)
	return st != nil && st.Occupied
}

// List returns all blocks with recorded state.
func (m *BlockMap) List() []*BlockState {
	blocks := m.store.list()
	sort.Slice(blocks, func(i, j int) bool { return blocks[i].BlockID < blocks[j].BlockID })
	return blocks
}

// Count returns how many blocks have state recorded.
func (m *BlockMap) Count() int {
	return m.store.count()
}

// Summary renders a compact human-readable overview for logs and dashboards.
func (m *BlockMap) Summary() string {
	blocks := m.List()
	occupied := 0
	for _, b := range blocks {
		if b.Occupied {
			occupied++
		}
	}
	return fmt.Sprintf("blocks=%d occupied=%d free=%d", len(blocks), occupied, len(blocks)-occupied)
}

// Snapshot captures the entire bitmap for inspection APIs.
type Snapshot struct {
	GeneratedAt int64        `json:"generated_at"`
	Blocks      []*BlockState `json:"blocks"`
}

// TakeSnapshot returns a point-in-time view of all block states.
func (m *BlockMap) TakeSnapshot(now int64) Snapshot {
	blocks := m.List()
	return Snapshot{GeneratedAt: now, Blocks: blocks}
}
