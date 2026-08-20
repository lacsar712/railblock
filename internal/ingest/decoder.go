package ingest

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/lacsar712/railblock/internal/codec"
)

// ContentKind identifies how frame bytes were transported.
type ContentKind int

const (
	KindRaw ContentKind = iota
	KindBase64
)

const (
	HeaderContentType = "Content-Type"
	HeaderSourceID    = "X-Source-ID"
	HeaderForceSign   = "X-Force-Sign"
	MimeOctets        = "application/octet-stream"
	MimeBase64        = "application/x-railblock-base64"
)

// ParseOptions controls frame body parsing.
type ParseOptions struct {
	MaxBytes int64
}

// DefaultParseOptions returns safe defaults for HTTP ingestion.
func DefaultParseOptions(max int64) ParseOptions {
	return ParseOptions{MaxBytes: max}
}

// DetectKind chooses raw vs base64 parsing based on content type.
func DetectKind(contentType string) ContentKind {
	ct := strings.ToLower(strings.TrimSpace(contentType))
	if strings.Contains(ct, "base64") {
		return KindBase64
	}
	return KindRaw
}

// ParseBody extracts frame bytes from an HTTP request body.
func ParseBody(body []byte, kind ContentKind) ([]byte, error) {
	switch kind {
	case KindRaw:
		return body, nil
	case KindBase64:
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(body)))
		if err != nil {
			return nil, fmt.Errorf("base64 decode: %w", err)
		}
		return decoded, nil
	default:
		return nil, fmt.Errorf("unknown content kind")
	}
}

// DecodePayload parses one or more frames from raw bytes.
func DecodePayload(raw []byte) (*codec.Frame, error) {
	if len(raw) == codec.FrameSize {
		f, err := codec.Decode(raw)
		if err != nil {
			return nil, fmt.Errorf("decode payload: %w", err)
		}
		return f, nil
	}
	frames, err := codec.DecodeMany(raw)
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	if len(frames) == 0 {
		return nil, codec.ErrShortFrame
	}
	return frames[0], nil
}
