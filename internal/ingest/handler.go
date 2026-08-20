package ingest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/codec"
	"github.com/lacsar712/railblock/internal/conflict"
)

// Service applies validated frames to the live system state.
type Service interface {
	ApplyFrame(frame interface {
		BlockID() uint16
	}, source bitmap.SourceID, signature string) (conflict.Result, *bitmap.BlockState)
}

// FrameApplier is the minimal frame surface used by ingest.
type FrameApplier interface {
	HasForceClear() bool
	IsOccupied() bool
}

// Processor coordinates decode and apply steps.
type Processor struct {
	Detector *conflict.Detector
	Options  ParseOptions
}

// NewProcessor constructs an ingest processor.
func NewProcessor(det *conflict.Detector, opts ParseOptions) *Processor {
	return &Processor{Detector: det, Options: opts}
}

// IngestResponse is returned after a successful frame submission.
type IngestResponse struct {
	OK             bool               `json:"ok"`
	Conflict       bool               `json:"conflict"`
	Block          bitmap.BlockReport `json:"block"`
	Applied        int                `json:"applied"`
	ConflictDetail []conflict.Record  `json:"conflict_detail,omitempty"`
}

// HandleFrames serves POST /v1/frames.
func (p *Processor) HandleFrames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, p.Options.MaxBytes)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	kind := DetectKind(r.Header.Get(HeaderContentType))
	raw, err := ParseBody(body, kind)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Decode the full frame stream up front. DecodeMany stops at the first
	// malformed frame and rejects any trailing partial bytes, so a corrupt
	// frame never leaves the bitmap half-applied and trailing data is never
	// silently dropped.
	frames, err := codec.DecodeMany(raw)
	if err != nil {
		status := http.StatusBadRequest
		if isCRCOrMagic(err) {
			status = http.StatusUnprocessableEntity
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	if len(frames) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty payload"})
		return
	}

	source := bitmap.SourceID(strings.TrimSpace(r.Header.Get(HeaderSourceID)))
	if source == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "X-Source-ID header required"})
		return
	}

	signature := strings.TrimSpace(r.Header.Get(HeaderForceSign))

	// Apply every decoded frame so a multi-frame packet updates all of its
	// blocks instead of only the first. The response reflects the aggregate
	// outcome: any conflict is surfaced, and Applied reports how many frames
	// were processed so callers can confirm the whole stream was handled.
	resp := IngestResponse{OK: true, Applied: len(frames)}
	var lastState *bitmap.BlockState
	for _, frame := range frames {
		confResult, st := p.Detector.ApplyFrame(frame, source, signature)
		if confResult.HasConflict {
			resp.Conflict = true
			resp.ConflictDetail = append(resp.ConflictDetail, confResult.Records...)
		}
		lastState = st
	}
	resp.Block = bitmap.ToReport(lastState)
	writeJSON(w, http.StatusAccepted, resp)
}

func isCRCOrMagic(err error) bool {
	return errors.Is(err, codec.ErrBadCRC) || errors.Is(err, codec.ErrBadMagic)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
