package app

import (
	"net/http"
	"strings"

	"github.com/lacsar712/railblock/internal/ingest"
	"github.com/lacsar712/railblock/internal/query"
	"github.com/lacsar712/railblock/internal/web"
)

// Router builds the HTTP routing table for the service.
type Router struct {
	mux *http.ServeMux
}

// NewRouter wires all HTTP endpoints.
func NewRouter(state *State) *Router {
	mux := http.NewServeMux()

	ingestProc := ingest.NewProcessor(state.Detector, ingest.DefaultParseOptions(state.Config.MaxBodyBytes))
	qh := query.NewHandler(query.Dependencies{
		Map:      state.BlockMap,
		Detector: state.Detector,
		Checker:  state.Checker,
		Policy:   state.Policy,
	})
	ui := web.NewHandler()

	mux.HandleFunc("/v1/frames/encode", ingest.HandleEncode)
	mux.HandleFunc("/v1/frames", ingestProc.HandleFrames)
	mux.HandleFunc("/v1/clearance/check", qh.HandleClearance)
	mux.HandleFunc("/v1/blocks/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/blocks/")
		if path == "" || path == "/" {
			qh.HandleSnapshot(w, r)
			return
		}
		qh.HandleBlock(w, r)
	})
	mux.HandleFunc("/v1/blocks", qh.HandleSnapshot)
	mux.HandleFunc("/healthz", qh.HandleHealth)
	RegisterDocs(mux)
	RegisterDiagnostics(mux, state)
	mux.HandleFunc("/v1/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		_ = WriteMetricsJSON(w, state.Snapshot())
	})
	mux.Handle("/", ui)

	return &Router{mux: mux}
}

// Handler returns the root HTTP handler with middleware applied.
func (rt *Router) Handler(state *State) http.Handler {
	h := http.Handler(rt.mux)
	h = ingest.CORSMiddleware(state.Config.EnableCORS, h)
	h = ingest.MaxBodyMiddleware(state.Config.MaxBodyBytes, h)
	h = ingest.LoggingMiddleware(state.Config.LogRequests, h)
	return h
}
