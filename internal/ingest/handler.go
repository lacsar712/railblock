package ingest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/lacsar712/railblock/internal/bitmap"
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
	OK        bool              `json:"ok"`
	Conflict  bool              `json:"conflict"`
	Block     bitmap.BlockReport  `json:"block"`
	ConflictDetail []conflict.Record `json:"conflict_detail,omitempty"`
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

	frame, err := DecodePayload(raw)
	if err != nil {
		status := http.StatusBadRequest
		if isCRCOrMagic(err) {
			status = http.StatusUnprocessableEntity
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	source := bitmap.SourceID(strings.TrimSpace(r.Header.Get(HeaderSourceID)))
	if source == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "X-Source-ID header required"})
		return
	}

	signature := strings.TrimSpace(r.Header.Get(HeaderForceSign))
	confResult, st := p.Detector.ApplyFrame(frame, source, signature)

	resp := IngestResponse{
		OK:       true,
		Conflict: confResult.HasConflict,
		Block:    bitmap.ToReport(st),
	}
	if confResult.HasConflict {
		resp.ConflictDetail = confResult.Records
	}
	writeJSON(w, http.StatusAccepted, resp)
}

func isCRCOrMagic(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "CRC") || strings.Contains(msg, "magic")
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
