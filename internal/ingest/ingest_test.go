package ingest_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/codec"
	"github.com/lacsar712/railblock/internal/conflict"
	"github.com/lacsar712/railblock/internal/ingest"
)

func TestBadMagicMapsToUnprocessable(t *testing.T) {
	m := bitmap.NewBlockMap()
	d := conflict.NewDetector(m, "test-secret-key-123")
	p := ingest.NewProcessor(d, ingest.DefaultParseOptions(1<<20))

	frame, _ := codec.NewFrame(1, 0, 1, 0)
	raw, _ := codec.Encode(frame)
	binary.BigEndian.PutUint32(raw[0:4], 0xDEADBEEF)

	req := httptest.NewRequest(http.MethodPost, "/v1/frames", bytes.NewReader(raw))
	req.Header.Set(ingest.HeaderContentType, ingest.MimeOctets)
	req.Header.Set(ingest.HeaderSourceID, "station-a")
	rr := httptest.NewRecorder()
	p.HandleFrames(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("bad magic status=%d want 422; body=%s", rr.Code, rr.Body.String())
	}
	var body map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &body)
	if body["error"] == "" {
		t.Fatal("expected error detail in response")
	}
}

func TestDecodePayloadRejectsMultiFrame(t *testing.T) {
	f1, _ := codec.NewFrame(1, 1, 1, 0)
	f2, _ := codec.NewFrame(2, 1, 2, 0)
	r1, _ := codec.Encode(f1)
	r2, _ := codec.Encode(f2)
	raw := append(append([]byte{}, r1...), r2...)

	_, err := ingest.DecodePayload(raw)
	if err == nil {
		t.Fatal("expected error for multi-frame payload")
	}
	if errors.Is(err, codec.ErrShortFrame) {
		t.Fatalf("unexpected short-frame error: %v", err)
	}
}

func TestDecodePayloadPreservesBadMagicUnwrap(t *testing.T) {
	frame, _ := codec.NewFrame(1, 0, 1, 0)
	raw, _ := codec.Encode(frame)
	binary.BigEndian.PutUint32(raw[0:4], 0xDEADBEEF)

	_, err := ingest.DecodePayload(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, codec.ErrBadMagic) {
		t.Fatalf("want unwrap ErrBadMagic, got %v", err)
	}
}
