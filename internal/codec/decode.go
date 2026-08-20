package codec

import "encoding/binary"

const payloadLen = 13

// Decode parses a 17-byte wire frame and validates magic, version, and CRC.
// On any validation failure the bitmap must not be updated by callers.
func Decode(data []byte) (*Frame, error) {
	if len(data) < FrameSize {
		return nil, newDecodeError(ErrShortFrame, "need 17 bytes")
	}

	magic := binary.BigEndian.Uint32(data[0:4])
	if magic != MagicValue {
		return nil, newDecodeError(ErrBadMagic, "expected RBLK")
	}

	version := data[4]
	if version != FrameVersion {
		return nil, newDecodeError(ErrBadVersion, "only version 1 supported")
	}

	if !VerifyCRC(data, payloadLen) {
		return nil, newDecodeError(ErrBadCRC, "checksum mismatch")
	}

	blockID := binary.BigEndian.Uint16(data[6:8])
	occupied := data[8]
	if occupied > 1 {
		return nil, newDecodeError(ErrBadOccupied, "invalid occupied value")
	}

	seq := binary.BigEndian.Uint32(data[9:13])

	return &Frame{
		Version:  version,
		Flags:    data[5],
		BlockID:  blockID,
		Occupied: occupied,
		Seq:      seq,
	}, nil
}

// DecodeMany parses a concatenated stream of frames, stopping at the first error.
func DecodeMany(data []byte) ([]*Frame, error) {
	var frames []*Frame
	for offset := 0; offset+FrameSize <= len(data); offset += FrameSize {
		f, err := Decode(data[offset : offset+FrameSize])
		if err != nil {
			return frames, err
		}
		frames = append(frames, f)
	}
	if len(data)%FrameSize != 0 {
		return frames, newDecodeError(ErrShortFrame, "trailing partial frame")
	}
	return frames, nil
}
