package query

// Endpoint catalog for documentation and OpenAPI-style introspection.

// Endpoint describes a public HTTP route.
type Endpoint struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

// Catalog returns the list of supported API endpoints.
func Catalog() []Endpoint {
	return []Endpoint{
		{Method: "POST", Path: "/v1/frames", Description: "Ingest a binary occupancy frame (raw octets or base64)"},
		{Method: "POST", Path: "/v1/frames/encode", Description: "Encode a demo frame to hex for testing"},
		{Method: "GET", Path: "/v1/blocks", Description: "Snapshot of all block states"},
		{Method: "GET", Path: "/v1/blocks/{id}", Description: "Query a single block partition"},
		{Method: "POST", Path: "/v1/clearance/check", Description: "Evaluate route clearance for interlocking"},
		{Method: "GET", Path: "/healthz", Description: "Health check"},
		{Method: "GET", Path: "/v1/metrics", Description: "Runtime metrics"},
		{Method: "GET", Path: "/", Description: "Embedded operator console"},
	}
}
