package conflict

import "github.com/lacsar712/railblock/internal/bitmap"

// Kind identifies the conflict rule that fired.
type Kind string

const (
	KindDualOccupancy Kind = "dual_occupancy"
	KindForceClearDenied Kind = "force_clear_denied"
	KindStaleSequence Kind = "stale_sequence"
)

// Record captures a detected conflict for a block partition.
type Record struct {
	BlockID uint16     `json:"block_id"`
	Kind    Kind       `json:"kind"`
	Message string     `json:"message"`
	Sources []string   `json:"sources,omitempty"`
}

// Result aggregates conflict evaluation output.
type Result struct {
	HasConflict bool      `json:"has_conflict"`
	Records     []Record  `json:"records,omitempty"`
}

// EmptyResult returns a non-conflicting result.
func EmptyResult() Result {
	return Result{HasConflict: false}
}

// WithRecord adds a conflict record and marks the result as conflicting.
func (r Result) WithRecord(rec Record) Result {
	r.HasConflict = true
	r.Records = append(r.Records, rec)
	return r
}

// BlocksAffected returns the set of block IDs implicated in conflicts.
func (r Result) BlocksAffected() []uint16 {
	seen := make(map[uint16]struct{})
	var out []uint16
	for _, rec := range r.Records {
		if _, ok := seen[rec.BlockID]; ok {
			continue
		}
		seen[rec.BlockID] = struct{}{}
		out = append(out, rec.BlockID)
	}
	return out
}

// ForBlock returns whether a specific block appears in the conflict set.
func (r Result) ForBlock(id uint16) bool {
	for _, rec := range r.Records {
		if rec.BlockID == id {
			return true
		}
	}
	return false
}

// FromBlockState builds a dual-occupancy record from a block snapshot.
func DualOccupancyRecord(st *bitmap.BlockState) Record {
	sources := st.OccupiedSources()
	names := make([]string, len(sources))
	for i, s := range sources {
		names[i] = s.String()
	}
	return Record{
		BlockID: st.BlockID,
		Kind:    KindDualOccupancy,
		Message: "multiple sources report occupied for the same block",
		Sources: names,
	}
}
