package bitmap

import (
	"fmt"
	"sort"
)

// SourceID identifies the reporting endpoint (station gateway, track circuit, etc.).
type SourceID string

func (s SourceID) String() string {
	if s == "" {
		return "<anonymous>"
	}
	return string(s)
}

// BlockState captures occupancy and provenance for a single block partition.
type BlockState struct {
	BlockID   uint16
	Occupied  bool
	Sources   map[SourceID]bool
	UpdatedAt int64
	Seq       uint32
}

// Clone returns a deep copy safe for external callers.
func (bs *BlockState) Clone() *BlockState {
	if bs == nil {
		return nil
	}
	cp := &BlockState{
		BlockID:   bs.BlockID,
		Occupied:  bs.Occupied,
		UpdatedAt: bs.UpdatedAt,
		Seq:       bs.Seq,
		Sources:   make(map[SourceID]bool, len(bs.Sources)),
	}
	for k, v := range bs.Sources {
		cp.Sources[k] = v
	}
	return cp
}

// OccupiedSourceCount returns how many distinct sources report occupied.
func (bs *BlockState) OccupiedSourceCount() int {
	if bs == nil {
		return 0
	}
	count := 0
	for src, occ := range bs.Sources {
		if occ && src != "" {
			count++
		}
	}
	return count
}

// OccupiedSources lists sources currently reporting occupied.
func (bs *BlockState) OccupiedSources() []SourceID {
	if bs == nil {
		return nil
	}
	var out []SourceID
	for src, occ := range bs.Sources {
		if occ {
			out = append(out, src)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (bs *BlockState) String() string {
	if bs == nil {
		return "BlockState<nil>"
	}
	state := "free"
	if bs.Occupied {
		state = "occupied"
	}
	return fmt.Sprintf("BlockState{id=%d state=%s sources=%d seq=%d}", bs.BlockID, state, len(bs.Sources), bs.Seq)
}
