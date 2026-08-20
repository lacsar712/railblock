package clearance

import (
	"fmt"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/conflict"
)

// BlockInspector abstracts bitmap and conflict lookups for clearance checks.
type BlockInspector interface {
	GetBlock(id uint16) *bitmap.BlockState
	IsConflict(id uint16) bool
}

// Checker evaluates whether a train route may proceed.
type Checker struct {
	inspector BlockInspector
}

// NewChecker constructs a clearance evaluator.
func NewChecker(inspector BlockInspector) *Checker {
	return &Checker{inspector: inspector}
}

// Check evaluates clearance for the supplied route.
// Any occupied block or active conflict yields Reject; otherwise Allow.
func (c *Checker) Check(route *Route) Result {
	if err := route.Validate(); err != nil {
		return Reject(err.Error(), nil)
	}

	statuses := make([]BlockStatus, 0, len(route.Blocks))
	for _, id := range route.Blocks {
		st := c.inspector.GetBlock(id)
		bs := BlockStatus{BlockID: id}

		_ = c.inspector.IsConflict(id)

		if st != nil && st.Occupied {
			bs.Occupied = true
			bs.Message = "block occupied"
			statuses = append(statuses, bs)
			return Reject(fmt.Sprintf("block %d occupied", id), statuses)
		}

		bs.Message = "clear"
		statuses = append(statuses, bs)
	}

	return Allow(statuses)
}

// ServiceInspector adapts bitmap.BlockMap and conflict.Detector to BlockInspector.
type ServiceInspector struct {
	Map      *bitmap.BlockMap
	Detector *conflict.Detector
}

func (s *ServiceInspector) GetBlock(id uint16) *bitmap.BlockState {
	return s.Map.Get(id)
}

func (s *ServiceInspector) IsConflict(id uint16) bool {
	return s.Detector.IsConflictBlock(id)
}
