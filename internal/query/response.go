package query

import "github.com/lacsar712/railblock/internal/bitmap"
import "github.com/lacsar712/railblock/internal/conflict"

// BlockResponse is returned by GET /v1/blocks/{id}.
type BlockResponse struct {
	Block    bitmap.BlockReport `json:"block"`
	Conflict bool               `json:"conflict"`
	Records  []conflict.Record  `json:"records,omitempty"`
}

// ErrorResponse is a generic API error envelope.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// NewErrorResponse builds a structured error payload.
func NewErrorResponse(msg, code, details string) ErrorResponse {
	return ErrorResponse{Error: msg, Code: code, Details: details}
}
