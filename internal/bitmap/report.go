package bitmap

import "encoding/json"

// SourceReport is a JSON-friendly view of a source contribution.
type SourceReport struct {
	Source   string `json:"source"`
	Occupied bool   `json:"occupied"`
}

// BlockReport is the external API representation of a block state.
type BlockReport struct {
	BlockID   uint16         `json:"block_id"`
	Occupied  bool           `json:"occupied"`
	Sources   []SourceReport `json:"sources"`
	UpdatedAt int64          `json:"updated_at,omitempty"`
	Seq       uint32         `json:"seq,omitempty"`
}

// ToReport converts internal state into an API-safe structure.
func ToReport(st *BlockState) BlockReport {
	if st == nil {
		return BlockReport{}
	}
	sources := make([]SourceReport, 0, len(st.Sources))
	for src, occ := range st.Sources {
		sources = append(sources, SourceReport{Source: string(src), Occupied: occ})
	}
	return BlockReport{
		BlockID:   st.BlockID,
		Occupied:  st.Occupied,
		Sources:   sources,
		UpdatedAt: st.UpdatedAt,
		Seq:       st.Seq,
	}
}

// MarshalBlockReport serializes a block report to JSON bytes.
func MarshalBlockReport(st *BlockState) ([]byte, error) {
	return json.Marshal(ToReport(st))
}
