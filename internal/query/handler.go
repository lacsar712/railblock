package query

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/clearance"
	"github.com/lacsar712/railblock/internal/conflict"
)

// Dependencies groups query-layer collaborators.
type Dependencies struct {
	Map      *bitmap.BlockMap
	Detector *conflict.Detector
	Checker  *clearance.Checker
	Policy   clearance.Policy
}

// Handler serves read and clearance endpoints.
type Handler struct {
	deps Dependencies
}

// NewHandler constructs query HTTP handlers.
func NewHandler(deps Dependencies) *Handler {
	return &Handler{deps: deps}
}

// HandleBlock serves GET /v1/blocks/{id}.
func (h *Handler) HandleBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/blocks/")
	idStr = strings.Trim(idStr, "/")
	id64, err := strconv.ParseUint(idStr, 10, 16)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid block id"})
		return
	}
	blockID := uint16(id64)

	st := h.deps.Map.Get(blockID)
	conf := h.deps.Detector.CheckBlock(blockID)

	resp := BlockResponse{
		Block:    bitmap.ToReport(st),
		Conflict: conf.HasConflict,
		Records:  conf.Records,
	}
	writeJSON(w, http.StatusOK, resp)
}

// HandleClearance serves POST /v1/clearance/check.
func (h *Handler) HandleClearance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var route clearance.Route
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result := h.deps.Checker.Check(&route)
	result = h.deps.Policy.Apply(&route, result)
	writeJSON(w, http.StatusOK, result)
}

// HandleSnapshot serves GET /v1/blocks for the operator console.
func (h *Handler) HandleSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snap := h.deps.Map.TakeSnapshot(time.Now().Unix())
	conf := h.deps.Detector.CheckAll()

	type blockView struct {
		bitmap.BlockReport
		Conflict bool `json:"conflict"`
	}
	var views []blockView
	for _, st := range snap.Blocks {
		views = append(views, blockView{
			BlockReport: bitmap.ToReport(st),
			Conflict:    conf.ForBlock(st.BlockID),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"generated_at": snap.GeneratedAt,
		"summary":      h.deps.Map.Summary(),
		"blocks":       views,
		"conflicts":    conf.Records,
	})
}

// HandleHealth serves GET /healthz.
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
