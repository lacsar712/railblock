package clearance

import (
	"fmt"
	"strings"
)

// RouteBuilder constructs validated clearance routes fluently.
type RouteBuilder struct {
	id     string
	blocks []uint16
	err    error
}

// NewRouteBuilder starts building a route with the given identifier.
func NewRouteBuilder(id string) *RouteBuilder {
	return &RouteBuilder{id: id}
}

// Add appends a block ID unless a duplicate is detected.
func (rb *RouteBuilder) Add(id uint16) *RouteBuilder {
	if rb.err != nil {
		return rb
	}
	for _, existing := range rb.blocks {
		if existing == id {
			rb.err = fmt.Errorf("duplicate block %d", id)
			return rb
		}
	}
	rb.blocks = append(rb.blocks, id)
	return rb
}

// AddMany appends multiple block IDs in order.
func (rb *RouteBuilder) AddMany(ids ...uint16) *RouteBuilder {
	for _, id := range ids {
		rb.Add(id)
		if rb.err != nil {
			return rb
		}
	}
	return rb
}

// Build returns the constructed route or a validation error.
func (rb *RouteBuilder) Build() (*Route, error) {
	if rb.err != nil {
		return nil, rb.err
	}
	route := &Route{RouteID: rb.id, Blocks: append([]uint16(nil), rb.blocks...)}
	if err := route.Validate(); err != nil {
		return nil, err
	}
	return route, nil
}

// Describe renders a concise route summary for logs.
func (r *Route) Describe() string {
	if r == nil {
		return "<nil route>"
	}
	parts := make([]string, len(r.Blocks))
	for i, id := range r.Blocks {
		parts[i] = fmt.Sprintf("%d", id)
	}
	return fmt.Sprintf("route %s: [%s]", r.RouteID, strings.Join(parts, ","))
}

// Split divides a long route into segments of at most size blocks.
func (r *Route) Split(size int) ([]*Route, error) {
	if r == nil {
		return nil, fmt.Errorf("route is nil")
	}
	if size <= 0 {
		return nil, fmt.Errorf("segment size must be positive")
	}
	if len(r.Blocks) == 0 {
		return nil, fmt.Errorf("empty route")
	}
	var segments []*Route
	for i := 0; i < len(r.Blocks); i += size {
		end := i + size
		if end > len(r.Blocks) {
			end = len(r.Blocks)
		}
		segments = append(segments, &Route{
			RouteID: fmt.Sprintf("%s-%d", r.RouteID, len(segments)),
			Blocks:  append([]uint16(nil), r.Blocks[i:end]...),
		})
	}
	return segments, nil
}

// Merge combines two routes when the tail of a meets the head of b.
func Merge(a, b *Route) (*Route, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("cannot merge nil routes")
	}
	combined := append(append([]uint16(nil), a.Blocks...), b.Blocks...)
	route := &Route{RouteID: a.RouteID + "+" + b.RouteID, Blocks: combined}
	if err := route.Validate(); err != nil {
		return nil, err
	}
	return route, nil
}
