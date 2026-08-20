package web

import "net/http"

// Register mounts the UI handler on a ServeMux (optional helper).
func Register(mux *http.ServeMux, h *Handler) {
	mux.Handle("/", h)
}

// ContentTypes used by the embedded UI.
const (
	HTMLContentType = "text/html; charset=utf-8"
	CSSContentType  = "text/css; charset=utf-8"
	JSContentType   = "application/javascript; charset=utf-8"
)
