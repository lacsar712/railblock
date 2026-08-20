package codec

import (
	"fmt"
	"strings"
)

// WireLayout documents the on-the-wire frame layout for tooling and diagnostics.
type WireLayout struct {
	Offset  int
	Length  int
	Name    string
	Type    string
	Notes   string
}

// FrameLayout returns the canonical byte layout table.
func FrameLayout() []WireLayout {
	return []WireLayout{
		{Offset: 0, Length: 4, Name: "magic", Type: "uint32 BE", Notes: fmt.Sprintf("0x%08X (RBLK)", MagicValue)},
		{Offset: 4, Length: 1, Name: "version", Type: "uint8", Notes: fmt.Sprintf("must be %d", FrameVersion)},
		{Offset: 5, Length: 1, Name: "flags", Type: "uint8", Notes: DescribeFlags(FlagForceClear | FlagPriority | FlagEmergency)},
		{Offset: 6, Length: 2, Name: "block_id", Type: "uint16 BE", Notes: "partition identifier"},
		{Offset: 8, Length: 1, Name: "occupied", Type: "uint8", Notes: "0=free, 1=occupied"},
		{Offset: 9, Length: 4, Name: "seq", Type: "uint32 BE", Notes: "monotonic sequence from source"},
		{Offset: 13, Length: 4, Name: "crc32", Type: "uint32 BE", Notes: "IEEE CRC-32 of bytes[0:13]"},
	}
}

// FormatLayout renders the layout as a human-readable table.
func FormatLayout() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Railblock frame (%d bytes)\n", FrameSize))
	b.WriteString("Offset  Len  Field      Type        Notes\n")
	b.WriteString("------  ---  ---------  ----------  -----\n")
	for _, row := range FrameLayout() {
		b.WriteString(fmt.Sprintf("0x%02X    %-3d  %-10s %-11s %s\n", row.Offset, row.Length, row.Name, row.Type, row.Notes))
	}
	return b.String()
}

// InspectWire returns a parsed summary of raw frame bytes for debugging.
func InspectWire(data []byte) (string, error) {
	f, err := Decode(data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s flags=%s crc=ok", f.String(), DescribeFlags(f.Flags)), nil
}

// CompareFrames reports whether two frames carry identical semantic content.
func CompareFrames(a, b *Frame) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Version == b.Version &&
		a.Flags == b.Flags &&
		a.BlockID == b.BlockID &&
		a.Occupied == b.Occupied &&
		a.Seq == b.Seq
}

// CloneFrame returns a shallow copy of a frame.
func CloneFrame(f *Frame) *Frame {
	if f == nil {
		return nil
	}
	cp := *f
	return &cp
}
