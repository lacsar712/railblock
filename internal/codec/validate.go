package codec

import "strings"

// ValidateFrame performs semantic checks beyond wire-format parsing.
func ValidateFrame(f *Frame) error {
	if f == nil {
		return newDecodeError(ErrShortFrame, "nil frame")
	}
	if f.Version != FrameVersion {
		return ErrBadVersion
	}
	if f.BlockID > MaxBlockID {
		return ErrBadBlockID
	}
	if f.Occupied > 1 {
		return ErrBadOccupied
	}
	return nil
}

// DescribeFlags returns a human-readable summary of active flag bits.
func DescribeFlags(flags uint8) string {
	if flags == FlagNone {
		return "none"
	}
	var parts []string
	if flags&FlagForceClear != 0 {
		parts = append(parts, "FORCE_CLEAR")
	}
	if flags&FlagPriority != 0 {
		parts = append(parts, "PRIORITY")
	}
	if flags&FlagEmergency != 0 {
		parts = append(parts, "EMERGENCY")
	}
	return strings.Join(parts, "|")
}
