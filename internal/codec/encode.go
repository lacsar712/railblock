package codec

import "encoding/binary"

// Encode serializes a frame into the canonical 17-byte wire format.
// Byte layout:
//
//	[0:4)   magic uint32 BE
//	[4]     version uint8
//	[5]     flags uint8
//	[6:8)   blockID uint16 BE
//	[8]     occupied uint8
//	[9:13)  seq uint32 BE
//	[13:17) crc32 IEEE of bytes[0:13]
func Encode(f *Frame) ([]byte, error) {
	if f == nil {
		return nil, newDecodeError(ErrShortFrame, "nil frame")
	}
	if f.Version != FrameVersion {
		return nil, ErrBadVersion
	}
	if f.BlockID > MaxBlockID {
		return nil, ErrBadBlockID
	}
	if f.Occupied > 1 {
		return nil, ErrBadOccupied
	}

	buf := make([]byte, FrameSize)
	binary.BigEndian.PutUint32(buf[0:4], MagicValue)
	buf[4] = f.Version
	buf[5] = f.Flags
	binary.BigEndian.PutUint16(buf[6:8], f.BlockID)
	buf[8] = f.Occupied
	binary.BigEndian.PutUint32(buf[9:13], f.Seq)
	AppendCRC(buf, buf[:13])
	return buf, nil
}

// MustEncode is a convenience helper for tests and internal callers that
// guarantee valid frames. It panics on encode failure.
func MustEncode(f *Frame) []byte {
	b, err := Encode(f)
	if err != nil {
		panic(err)
	}
	return b
}
