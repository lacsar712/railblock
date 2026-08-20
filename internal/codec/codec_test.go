package codec_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/lacsar712/railblock/internal/codec"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	frame, err := codec.NewFrame(12, 1, 42, codec.FlagNone)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := codec.Encode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) != codec.FrameSize {
		t.Fatalf("size=%d want %d", len(raw), codec.FrameSize)
	}
	got, err := codec.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.BlockID != 12 || got.Occupied != 1 || got.Seq != 42 {
		t.Fatalf("decode mismatch: %+v", got)
	}
}

func TestBadMagic(t *testing.T) {
	frame, _ := codec.NewFrame(1, 0, 1, 0)
	raw, _ := codec.Encode(frame)
	binary.BigEndian.PutUint32(raw[0:4], 0xDEADBEEF)
	_, err := codec.Decode(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	var dec *codec.DecodeError
	if !errors.As(err, &dec) {
		t.Fatalf("expected DecodeError, got %T", err)
	}
	if !errors.Is(err, codec.ErrBadMagic) {
		t.Fatalf("got %v", err)
	}
}

func TestBadCRC(t *testing.T) {
	frame, _ := codec.NewFrame(5, 1, 9, 0)
	raw, _ := codec.Encode(frame)
	raw[16] ^= 0xFF
	_, err := codec.Decode(raw)
	if !errors.Is(err, codec.ErrBadCRC) {
		t.Fatalf("got %v", err)
	}
}

func TestBadVersion(t *testing.T) {
	frame, _ := codec.NewFrame(5, 1, 9, 0)
	raw, _ := codec.Encode(frame)
	raw[4] = 99
	// recompute crc
	copy(raw[13:17], []byte{0, 0, 0, 0})
	sum := codec.CRC32IEEE(raw[:13])
	buf := bytes.NewBuffer(nil)
	_ = binary.Write(buf, binary.BigEndian, sum)
	// manual crc write
	raw[13] = byte(sum >> 24)
	raw[14] = byte(sum >> 16)
	raw[15] = byte(sum >> 8)
	raw[16] = byte(sum)
	_, err := codec.Decode(raw)
	if !errors.Is(err, codec.ErrBadVersion) {
		t.Fatalf("got %v", err)
	}
}

func TestBadFrameDoesNotIncludeOccupiedTwo(t *testing.T) {
	frame := &codec.Frame{Version: 1, BlockID: 1, Occupied: 2, Seq: 1}
	_, err := codec.Encode(frame)
	if !errors.Is(err, codec.ErrBadOccupied) {
		t.Fatalf("got %v", err)
	}
}

func TestHexRoundTrip(t *testing.T) {
	f, _ := codec.NewFrame(100, 0, 7, 0)
	hex, err := codec.HexEncode(f)
	if err != nil {
		t.Fatal(err)
	}
	got, err := codec.HexDecode(hex)
	if err != nil {
		t.Fatal(err)
	}
	if got.BlockID != 100 {
		t.Fatalf("block=%d", got.BlockID)
	}
}

func TestAppendCRCMatchesVerify(t *testing.T) {
	prefix := []byte{
		0x52, 0x42, 0x4C, 0x4B, // magic
		0x01, 0x00, // ver, flags
		0x00, 0x0C, // block 12 BE
		0x01,                   // occupied
		0x00, 0x00, 0x00, 0x2A, // seq 42
	}
	buf := make([]byte, len(prefix)+4)
	copy(buf, prefix)
	codec.AppendCRC(buf, prefix)
	if !codec.VerifyCRC(buf, len(prefix)) {
		t.Fatal("AppendCRC output must VerifyCRC with matching payload length and endian")
	}
}

func TestDecodeManyRejectsTrailingPartial(t *testing.T) {
	f, _ := codec.NewFrame(1, 1, 1, 0)
	raw, _ := codec.Encode(f)
	raw = append(raw, 0x00, 0x01, 0x02) // trailing partial frame
	_, err := codec.DecodeMany(raw)
	if !errors.Is(err, codec.ErrShortFrame) {
		t.Fatalf("trailing partial: got %v want ErrShortFrame", err)
	}
}

func TestBlockIDBigEndianRoundTrip(t *testing.T) {
	const want uint16 = 0x1234
	f, err := codec.NewFrame(want, 1, 9, 0)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := codec.Encode(f)
	if err != nil {
		t.Fatal(err)
	}
	if got := binary.BigEndian.Uint16(raw[6:8]); got != want {
		t.Fatalf("wire block_id=%#04x want %#04x (big-endian)", got, want)
	}
	decoded, err := codec.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.BlockID != want {
		t.Fatalf("decoded block_id=%#04x want %#04x", decoded.BlockID, want)
	}
}
