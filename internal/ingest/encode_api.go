package ingest

import (
	"encoding/json"
	"net/http"

	"github.com/lacsar712/railblock/internal/codec"
)

// EncodeRequest is used by the operator UI to build demo frames.
type EncodeRequest struct {
	BlockID  uint16 `json:"block_id"`
	Occupied uint8  `json:"occupied"`
	Seq      uint32 `json:"seq"`
	Flags    uint8  `json:"flags"`
}

// EncodeResponse returns a hex-encoded wire frame for client-side submission.
type EncodeResponse struct {
	Hex   string `json:"hex"`
	Bytes int    `json:"bytes"`
}

// HandleEncode serves POST /v1/frames/encode for the embedded UI.
func HandleEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req EncodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	frame, err := codec.NewFrame(req.BlockID, req.Occupied, req.Seq, req.Flags)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	hexStr, err := codec.HexEncode(frame)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, EncodeResponse{Hex: hexStr, Bytes: codec.FrameSize})
}
