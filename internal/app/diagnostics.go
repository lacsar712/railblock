package app

import (
	"encoding/json"
	"net/http"

	"github.com/lacsar712/railblock/internal/query"
)

// RegisterDocs mounts a JSON catalog of API endpoints when enabled.
func RegisterDocs(mux *http.ServeMux) {
	mux.HandleFunc("/v1/catalog", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"service":   "railblock",
			"version":   Version,
			"endpoints": query.Catalog(),
		})
	})
}

// RegisterDiagnostics exposes bitmap table output for operators.
func RegisterDiagnostics(mux *http.ServeMux, state *State) {
	mux.HandleFunc("/v1/diagnostics/bitmap", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(state.BlockMap.FormatTable()))
	})
}
