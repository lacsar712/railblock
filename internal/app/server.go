package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server wraps http.Server with graceful lifecycle helpers.
type Server struct {
	cfg    *configWrapper
	state  *State
	router *Router
	http   *http.Server
}

type configWrapper struct {
	addr         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

// NewServer constructs a listening HTTP server without starting it.
func NewServer(state *State) *Server {
	rt := NewRouter(state)
	return &Server{
		cfg: &configWrapper{
			addr:         state.Config.Addr,
			readTimeout:  state.Config.ReadTimeout,
			writeTimeout: state.Config.WriteTimeout,
			idleTimeout:  state.Config.IdleTimeout,
		},
		state:  state,
		router: rt,
		http: &http.Server{
			Addr:         state.Config.Addr,
			Handler:      rt.Handler(state),
			ReadTimeout:  state.Config.ReadTimeout,
			WriteTimeout: state.Config.WriteTimeout,
			IdleTimeout:  state.Config.IdleTimeout,
		},
	}
}

// ListenAndServe starts the HTTP server and blocks until it exits.
func (s *Server) ListenAndServe() error {
	log.Printf("railblock listening on %s", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("railblock shutting down")
	return s.http.Shutdown(ctx)
}

// Addr returns the configured listen address.
func (s *Server) Addr() string { return s.http.Addr }

// WriteMetricsJSON encodes metrics for the /v1/metrics endpoint.
func WriteMetricsJSON(w http.ResponseWriter, m Metrics) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(m)
}
