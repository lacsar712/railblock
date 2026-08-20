package app

import (
	"sync"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/clearance"
	"github.com/lacsar712/railblock/internal/conflict"
	"github.com/lacsar712/railblock/internal/config"
)

// State holds the shared in-memory service state.
type State struct {
	mu       sync.RWMutex
	Config   *config.Config
	BlockMap *bitmap.BlockMap
	Detector *conflict.Detector
	Checker  *clearance.Checker
	Policy   clearance.Policy
}

// NewState constructs initialized application state from configuration.
func NewState(cfg *config.Config) *State {
	m := bitmap.NewBlockMap()
	det := conflict.NewDetector(m, cfg.ForceClearSecret)
	inspector := &clearance.ServiceInspector{Map: m, Detector: det}
	checker := clearance.NewChecker(inspector)

	return &State{
		Config:   cfg,
		BlockMap: m,
		Detector: det,
		Checker:  checker,
		Policy:   clearance.DefaultPolicy(),
	}
}

// Snapshot returns a copy of operational counters for metrics.
func (s *State) Snapshot() Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conf := s.Detector.CheckAll()
	return Metrics{
		BlockCount:     s.BlockMap.Count(),
		ConflictCount:  len(conf.Records),
		ConflictActive: conf.HasConflict,
	}
}

// Metrics exposes lightweight runtime counters.
type Metrics struct {
	BlockCount     int  `json:"block_count"`
	ConflictCount  int  `json:"conflict_count"`
	ConflictActive bool `json:"conflict_active"`
}
