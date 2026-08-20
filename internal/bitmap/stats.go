package bitmap

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Stats summarizes occupancy across the entire bitmap.
type Stats struct {
	TotalBlocks    int `json:"total_blocks"`
	OccupiedBlocks int `json:"occupied_blocks"`
	FreeBlocks     int `json:"free_blocks"`
	SourceCount    int `json:"source_count"`
}

// ComputeStats derives aggregate statistics from the current bitmap.
func (m *BlockMap) ComputeStats() Stats {
	blocks := m.List()
	occupied := 0
	sources := make(map[SourceID]struct{})
	for _, st := range blocks {
		if st.Occupied {
			occupied++
		}
		for src, occ := range st.Sources {
			if occ {
				sources[src] = struct{}{}
			}
		}
	}
	return Stats{
		TotalBlocks:    len(blocks),
		OccupiedBlocks: occupied,
		FreeBlocks:     len(blocks) - occupied,
		SourceCount:    len(sources),
	}
}

// ExportJSON serializes all block states to JSON for archival.
func (m *BlockMap) ExportJSON() ([]byte, error) {
	snap := m.TakeSnapshot(0)
	type export struct {
		Blocks []BlockReport `json:"blocks"`
		Stats  Stats         `json:"stats"`
	}
	reports := make([]BlockReport, 0, len(snap.Blocks))
	for _, st := range snap.Blocks {
		reports = append(reports, ToReport(st))
	}
	payload := export{Blocks: reports, Stats: m.ComputeStats()}
	return json.Marshal(payload)
}

// ImportJSON restores block states from JSON (testing and recovery helper).
func (m *BlockMap) ImportJSON(data []byte) error {
	var payload struct {
		Blocks []struct {
			BlockID  uint16         `json:"block_id"`
			Occupied bool           `json:"occupied"`
			Sources  []SourceReport `json:"sources"`
			Seq      uint32         `json:"seq"`
		} `json:"blocks"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	for _, b := range payload.Blocks {
		for _, s := range b.Sources {
			m.Set(b.BlockID, s.Occupied, SourceID(s.Source), b.Seq)
		}
		if len(b.Sources) == 0 && b.Occupied {
			m.Set(b.BlockID, true, "import", b.Seq)
		}
	}
	return nil
}

// FormatTable renders a fixed-width text table of block states.
func (m *BlockMap) FormatTable() string {
	blocks := m.List()
	if len(blocks) == 0 {
		return "no blocks recorded"
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-8s %-10s %-30s\n", "BLOCK", "STATE", "SOURCES"))
	b.WriteString(strings.Repeat("-", 52) + "\n")
	for _, st := range blocks {
		state := "free"
		if st.Occupied {
			state = "occupied"
		}
		srcs := st.OccupiedSources()
		names := make([]string, len(srcs))
		for i, s := range srcs {
			names[i] = s.String()
		}
		sort.Strings(names)
		b.WriteString(fmt.Sprintf("%-8d %-10s %-30s\n", st.BlockID, state, strings.Join(names, ",")))
	}
	return b.String()
}

// FindOccupiedInRange returns occupied blocks within [from, to] inclusive.
func (m *BlockMap) FindOccupiedInRange(from, to uint16) []*BlockState {
	var out []*BlockState
	for _, st := range m.List() {
		if st.BlockID >= from && st.BlockID <= to && st.Occupied {
			out = append(out, st)
		}
	}
	return out
}
