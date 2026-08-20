package clearance

// Decision represents the outcome of a clearance evaluation.
type Decision string

const (
	DecisionAllow  Decision = "Allow"
	DecisionReject Decision = "Reject"
)

// Result is returned to interlocking systems querying route clearance.
type Result struct {
	Decision Decision `json:"decision"`
	Reason   string   `json:"reason,omitempty"`
	Blocks   []BlockStatus `json:"blocks,omitempty"`
}

// BlockStatus describes per-block clearance contribution.
type BlockStatus struct {
	BlockID   uint16 `json:"block_id"`
	Occupied  bool   `json:"occupied"`
	Conflict  bool   `json:"conflict"`
	Message   string `json:"message,omitempty"`
}

// Allow creates an Allow result for the given blocks.
func Allow(blocks []BlockStatus) Result {
	return Result{Decision: DecisionAllow, Reason: "all blocks free", Blocks: blocks}
}

// Reject creates a Reject result with an explanatory reason.
func Reject(reason string, blocks []BlockStatus) Result {
	return Result{Decision: DecisionReject, Reason: reason, Blocks: blocks}
}

// IsAllowed reports whether clearance was granted.
func (r Result) IsAllowed() bool { return r.Decision == DecisionAllow }
