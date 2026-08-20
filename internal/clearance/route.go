package clearance

import "fmt"

// Route describes a train movement path as an ordered list of block IDs.
type Route struct {
	RouteID string   `json:"route_id"`
	Blocks  []uint16 `json:"blocks"`
}

// Validate ensures the route request is well-formed.
func (r *Route) Validate() error {
	if r == nil {
		return fmt.Errorf("route is nil")
	}
	if len(r.Blocks) == 0 {
		return fmt.Errorf("route must include at least one block")
	}
	seen := make(map[uint16]struct{}, len(r.Blocks))
	for _, id := range r.Blocks {
		if _, dup := seen[id]; dup {
			return fmt.Errorf("duplicate block %d in route", id)
		}
		seen[id] = struct{}{}
	}
	return nil
}

// Contains reports whether a block ID is part of the route.
func (r *Route) Contains(id uint16) bool {
	for _, b := range r.Blocks {
		if b == id {
			return true
		}
	}
	return false
}

// Len returns the number of blocks in the route.
func (r *Route) Len() int {
	if r == nil {
		return 0
	}
	return len(r.Blocks)
}
