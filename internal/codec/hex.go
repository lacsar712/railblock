package codec

import (
	"encoding/hex"
	"fmt"
)

// HexEncode returns a lowercase hex string of a serialized frame.
func HexEncode(f *Frame) (string, error) {
	raw, err := Encode(f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

// HexDecode parses a hex-encoded wire frame.
func HexDecode(s string) (*Frame, error) {
	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("hex decode: %w", err)
	}
	return Decode(raw)
}

// FormatWire returns a spaced hex dump suitable for logs.
func FormatWire(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	const group = 4
	var parts []string
	for i := 0; i < len(data); i += group {
		end := i + group
		if end > len(data) {
			end = len(data)
		}
		parts = append(parts, hex.EncodeToString(data[i:end]))
	}
	return fmt.Sprintf("%s (%d bytes)", joinParts(parts), len(data))
}

func joinParts(parts []string) string {
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += " " + parts[i]
	}
	return out
}
