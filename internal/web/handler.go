package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed assets/*
var assets embed.FS

// Handler serves the embedded operator console and static assets.
type Handler struct {
	static http.Handler
	index  []byte
}

// NewHandler loads embedded web assets.
func NewHandler() *Handler {
	staticSub, err := fs.Sub(assets, "assets/static")
	if err != nil {
		panic(err)
	}
	index, err := assets.ReadFile("assets/templates/index.html")
	if err != nil {
		panic(err)
	}
	return &Handler{
		static: http.FileServer(http.FS(staticSub)),
		index:  index,
	}
}

// ServeHTTP implements http.Handler for the UI and static files.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" || path == "/index.html" {
		w.Header().Set("Content-Type", HTMLContentType)
		_, _ = w.Write(h.index)
		return
	}
	if strings.HasPrefix(path, "/static/") {
		http.StripPrefix("/static/", h.static).ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}
