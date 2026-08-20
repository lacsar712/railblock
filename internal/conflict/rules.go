package conflict

import "github.com/lacsar712/railblock/internal/bitmap"

// Rule identifiers for documentation and audit trails.
const (
	RuleR1DualOccupancy = "R1"
	RuleR2ForceClear    = "R2"
)

// Rule describes a conflict detection policy.
type Rule struct {
	ID          string
	Description string
}

// DefaultRules returns the built-in conflict rule catalog.
func DefaultRules() []Rule {
	return []Rule{
		{ID: RuleR1DualOccupancy, Description: "same block from >=2 sources both occupied"},
		{ID: RuleR2ForceClear, Description: "FORCE_CLEAR allows single-source clear when signed"},
	}
}

// EvaluateR1 checks whether two or more sources simultaneously claim occupancy.
func EvaluateR1(st *bitmap.BlockState) (Record, bool) {
	if st == nil {
		return Record{}, false
	}
	if st.OccupiedSourceCount() >= 2 {
		return DualOccupancyRecord(st), true
	}
	return Record{}, false
}

// BlocksInConflict scans a list of block states and returns R1 violations.
func BlocksInConflict(states []*bitmap.BlockState) Result {
	result := EmptyResult()
	for _, st := range states {
		if rec, ok := EvaluateR1(st); ok {
			result = result.WithRecord(rec)
		}
	}
	return result
}
