package codec

import "fmt"

const (
	FrameSize     = 17
	MagicValue    = 0x52424C4B // "RBLK"
	FrameVersion  = 1
	MaxBlockID    = 65535
)

// Flag bits carried in the frame header.
const (
	FlagNone       uint8 = 0
	FlagForceClear uint8 = 1 << 0
	FlagPriority   uint8 = 1 << 1
	FlagEmergency  uint8 = 1 << 2
)

// Frame represents a fixed-length rail block occupancy report.
type Frame struct {
	Version  uint8
	Flags    uint8
	BlockID  uint16
	Occupied uint8
	Seq      uint32
}

// NewFrame constructs a validated frame ready for encoding.
func NewFrame(blockID uint16, occupied uint8, seq uint32, flags uint8) (*Frame, error) {
	if blockID > MaxBlockID {
		return nil, ErrBadBlockID
	}
	if occupied > 1 {
		return nil, ErrBadOccupied
	}
	return &Frame{
		Version:  FrameVersion,
		Flags:    flags,
		BlockID:  blockID,
		Occupied: occupied,
		Seq:      seq,
	}, nil
}

// HasForceClear reports whether the frame requests a forced clear operation.
func (f *Frame) HasForceClear() bool {
	return f.Flags&FlagForceClear != 0
}

// IsOccupied reports whether the frame declares the block occupied.
func (f *Frame) IsOccupied() bool {
	return f.Occupied == 1
}

// String returns a concise debug representation.
func (f *Frame) String() string {
	state := "free"
	if f.IsOccupied() {
		state = "occupied"
	}
	return fmt.Sprintf("Frame{block=%d state=%s seq=%d flags=0x%02x}", f.BlockID, state, f.Seq, f.Flags)
}
