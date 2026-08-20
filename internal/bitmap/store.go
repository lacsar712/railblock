package bitmap

import (
	"sync"
	"time"
)

// Store is the low-level concurrent map of block states.
type Store struct {
	mu     sync.RWMutex
	blocks map[uint16]*BlockState
	clock  func() int64
}

// NewStore creates an empty block store.
func NewStore() *Store {
	return &Store{
		blocks: make(map[uint16]*BlockState),
		clock:  func() int64 { return time.Now().Unix() },
	}
}

// WithClock overrides the timestamp source (used in tests).
func (s *Store) WithClock(clock func() int64) *Store {
	s.clock = clock
	return s
}

func (s *Store) getOrCreate(id uint16) *BlockState {
	st, ok := s.blocks[id]
	if !ok {
		st = &BlockState{
			BlockID: id,
			Sources: make(map[SourceID]bool),
		}
		s.blocks[id] = st
	}
	return st
}

// apply writes occupancy for a source into an existing block entry.
func (s *Store) apply(id uint16, occupied bool, source SourceID, seq uint32) *BlockState {
	s.mu.Lock()
	defer s.mu.Unlock()

	st := s.getOrCreate(id)
	st.Sources[source] = occupied
	st.Seq = seq
	st.UpdatedAt = s.clock()

	anyOccupied := false
	for _, occ := range st.Sources {
		if occ {
			anyOccupied = true
			break
		}
	}
	st.Occupied = anyOccupied
	return st.Clone()
}

// clearSource removes a source contribution and recomputes aggregate occupancy.
func (s *Store) clearSource(id uint16, source SourceID, seq uint32) *BlockState {
	s.mu.Lock()
	defer s.mu.Unlock()

	st, ok := s.blocks[id]
	if !ok {
		return &BlockState{BlockID: id, Sources: map[SourceID]bool{}}
	}
	delete(st.Sources, source)
	st.Seq = seq
	st.UpdatedAt = s.clock()

	anyOccupied := false
	for _, occ := range st.Sources {
		if occ {
			anyOccupied = true
			break
		}
	}
	st.Occupied = anyOccupied
	if len(st.Sources) == 0 && !st.Occupied {
		delete(s.blocks, id)
		return &BlockState{BlockID: id, Sources: map[SourceID]bool{}}
	}
	return st.Clone()
}

// snapshot returns a clone of the block state or a default free state.
func (s *Store) snapshot(id uint16) *BlockState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.blocks[id]
	if !ok {
		return &BlockState{BlockID: id, Sources: map[SourceID]bool{}}
	}
	return st.Clone()
}

// list returns clones of all known blocks sorted by ID.
func (s *Store) list() []*BlockState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*BlockState, 0, len(s.blocks))
	for _, st := range s.blocks {
		out = append(out, st.Clone())
	}
	return out
}

// count returns the number of tracked blocks.
func (s *Store) count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blocks)
}
