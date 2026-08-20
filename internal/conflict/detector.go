package conflict

import (
	"sync"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/codec"
)

// Detector evaluates occupancy updates against conflict rules.
type Detector struct {
	mu        sync.RWMutex
	mapRef    *bitmap.BlockMap
	signer    *Signer
	forceEval *ForceClearEvaluator
	history   []Record
	maxHist   int
}

// NewDetector wires a bitmap with signing support for FORCE_CLEAR frames.
func NewDetector(m *bitmap.BlockMap, secret string) *Detector {
	signer := NewSigner(secret)
	return &Detector{
		mapRef:    m,
		signer:    signer,
		forceEval: NewForceClearEvaluator(signer),
		maxHist:   256,
	}
}

// Signer exposes the HMAC signer for clients that need to mint signatures.
func (d *Detector) Signer() *Signer { return d.signer }

// CheckBlock evaluates R1 for a single block ID.
func (d *Detector) CheckBlock(blockID uint16) Result {
	st := d.mapRef.Get(blockID)
	return BlocksInConflict([]*bitmap.BlockState{st})
}

// CheckAll scans every tracked block for conflicts.
func (d *Detector) CheckAll() Result {
	return BlocksInConflict(d.mapRef.List())
}

// ApplyFrame validates and applies a decoded frame to the bitmap.
// CRC and wire-format validation must happen before calling this method.
func (d *Detector) ApplyFrame(frame *codec.Frame, source bitmap.SourceID, signature string) (Result, *bitmap.BlockState) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if frame.HasForceClear() {
		if err := d.forceEval.AllowClear(frame, source, signature); err != nil {
			rec := DeniedRecord(frame.BlockID, err.Error())
			d.appendHistory(rec)
			return Result{HasConflict: true, Records: []Record{rec}}, d.mapRef.Get(frame.BlockID)
		}
		st := d.mapRef.ForceClear(frame.BlockID, source, frame.Seq)
		return EmptyResult(), st
	}

	occupied := frame.IsOccupied()
	st := d.mapRef.Set(frame.BlockID, occupied, source, frame.Seq)

	if rec, ok := EvaluateR1(st); ok {
		d.appendHistory(rec)
		return Result{HasConflict: true, Records: []Record{rec}}, st
	}
	return EmptyResult(), st
}

func (d *Detector) appendHistory(rec Record) {
	d.history = append(d.history, rec)
	if len(d.history) > d.maxHist {
		d.history = d.history[len(d.history)-d.maxHist:]
	}
}

// History returns recent conflict records (newest last).
func (d *Detector) History() []Record {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]Record, len(d.history))
	copy(out, d.history)
	return out
}

// IsConflictBlock reports whether a block currently violates R1.
func (d *Detector) IsConflictBlock(blockID uint16) bool {
	res := d.CheckBlock(blockID)
	return res.HasConflict
}
