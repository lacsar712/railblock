package ingest

import (
	"fmt"
	"strings"

	"github.com/lacsar712/railblock/internal/bitmap"
)

// ValidateSource ensures the reporting source identifier is acceptable.
func ValidateSource(source bitmap.SourceID) error {
	s := strings.TrimSpace(string(source))
	if s == "" {
		return fmt.Errorf("source id is required")
	}
	if len(s) > 128 {
		return fmt.Errorf("source id too long")
	}
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return fmt.Errorf("source id contains invalid character: %q", r)
	}
	return nil
}

// ValidatePayloadSize checks raw payload length before decode.
func ValidatePayloadSize(raw []byte, max int64) error {
	if int64(len(raw)) > max {
		return fmt.Errorf("payload exceeds max %d bytes", max)
	}
	if len(raw) == 0 {
		return fmt.Errorf("empty payload")
	}
	return nil
}

// NormalizeSource trims and validates a source header value.
func NormalizeSource(raw string) (bitmap.SourceID, error) {
	src := bitmap.SourceID(strings.TrimSpace(raw))
	if err := ValidateSource(src); err != nil {
		return "", err
	}
	return src, nil
}
