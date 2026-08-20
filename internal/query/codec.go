package query

import (
	"encoding/json"
	"net/http"

	"github.com/lacsar712/railblock/internal/clearance"
)

// EncodeResult writes a clearance result as JSON.
func EncodeResult(w http.ResponseWriter, result clearance.Result) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(result)
}

// DecodeRoute parses a clearance route from JSON bytes.
func DecodeRoute(data []byte) (*clearance.Route, error) {
	var route clearance.Route
	if err := json.Unmarshal(data, &route); err != nil {
		return nil, err
	}
	if err := route.Validate(); err != nil {
		return nil, err
	}
	return &route, nil
}

// NotFound writes a 404 JSON response.
func NotFound(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
