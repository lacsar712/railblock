package clearance

// Policy tunes clearance behavior for operational modes.
type Policy struct {
	RejectOnEmptyRoute bool
	StrictConflict     bool
}

// DefaultPolicy returns production clearance defaults.
func DefaultPolicy() Policy {
	return Policy{
		RejectOnEmptyRoute: true,
		StrictConflict:     true,
	}
}

// Apply adjusts a clearance result according to policy (currently a pass-through hook).
func (p Policy) Apply(route *Route, result Result) Result {
	if p.RejectOnEmptyRoute && route.Len() == 0 {
		return Reject("empty route", nil)
	}
	return result
}

// Describe returns a short policy summary for diagnostics endpoints.
func (p Policy) Describe() map[string]bool {
	return map[string]bool{
		"reject_on_empty_route": p.RejectOnEmptyRoute,
		"strict_conflict":       p.StrictConflict,
	}
}
