package conflict_test

import (
	"errors"
	"testing"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/codec"
	"github.com/lacsar712/railblock/internal/conflict"
)

func TestR1DualOccupancy(t *testing.T) {
	m := bitmap.NewBlockMap()
	d := conflict.NewDetector(m, "test-secret-key-123")
	srcA := bitmap.SourceID("station-a")
	srcB := bitmap.SourceID("station-b")

	f1, _ := codec.NewFrame(12, 1, 1, 0)
	d.ApplyFrame(f1, srcA, "")

	f2, _ := codec.NewFrame(12, 1, 2, 0)
	res, _ := d.ApplyFrame(f2, srcB, "")
	if !res.HasConflict {
		t.Fatal("expected conflict")
	}
	if !d.IsConflictBlock(12) {
		t.Fatal("block should be in conflict")
	}
}

func TestBadCRCFrameNotApplied(t *testing.T) {
	m := bitmap.NewBlockMap()
	d := conflict.NewDetector(m, "test-secret-key-123")

	f, err := codec.NewFrame(5, 1, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := codec.Encode(f)
	if err != nil {
		t.Fatal(err)
	}
	raw[16] ^= 0xFF

	decoded, err := codec.Decode(raw)
	if err == nil {
		// STACK-friendly: applying a CRC-corrupt frame must never mutate the bitmap.
		res, st := d.ApplyFrame(decoded, "station-a", "")
		t.Fatalf("STACK: bad CRC frame applied to bitmap: conflict=%v occupied=%v sources=%v",
			res.HasConflict, st != nil && st.Occupied, st)
	}
	if !errors.Is(err, codec.ErrBadCRC) {
		t.Fatalf("decode want ErrBadCRC, got %v", err)
	}
	st := m.Get(5)
	if st.Occupied || st.OccupiedSourceCount() != 0 {
		t.Fatalf("STACK: bitmap mutated despite CRC failure: %+v", st)
	}
}

func TestForceClearSigned(t *testing.T) {
	m := bitmap.NewBlockMap()
	d := conflict.NewDetector(m, "test-secret-key-123")
	src := bitmap.SourceID("station-a")
	occ, _ := codec.NewFrame(9, 1, 9, 0)
	d.ApplyFrame(occ, src, "")

	f, _ := codec.NewFrame(9, 0, 10, codec.FlagForceClear)
	sig := d.Signer().Sign(9, src, 10)
	res, st := d.ApplyFrame(f, src, sig)
	if res.HasConflict {
		t.Fatalf("unexpected conflict: %+v", res)
	}
	if st == nil || st.Occupied {
		t.Fatal("expected cleared")
	}
}

func TestForceClearUnsignedDenied(t *testing.T) {
	m := bitmap.NewBlockMap()
	d := conflict.NewDetector(m, "test-secret-key-123")
	f, _ := codec.NewFrame(9, 0, 11, codec.FlagForceClear)
	res, _ := d.ApplyFrame(f, "station-a", "")
	if !res.HasConflict {
		t.Fatal("expected denied force clear")
	}
}

func TestForceClearHMACPayloadMatchesVerify(t *testing.T) {
	m := bitmap.NewBlockMap()
	secret := "test-secret-key-123"
	d := conflict.NewDetector(m, secret)
	src := bitmap.SourceID("Station-Alpha")
	occ, _ := codec.NewFrame(9, 1, 9, 0)
	d.ApplyFrame(occ, src, "")

	const seq uint32 = 10
	sig := d.Signer().Sign(9, src, seq)
	if !d.Signer().Verify(9, src, seq, sig) {
		t.Fatalf("STACK: Sign/Verify HMAC payload mismatch for source=%q seq=%d sig=%s", src, seq, sig)
	}

	f, _ := codec.NewFrame(9, 0, seq, codec.FlagForceClear)
	res, st := d.ApplyFrame(f, src, sig)
	if res.HasConflict {
		t.Fatalf("STACK: force-clear denied despite matching HMAC: %+v", res)
	}
	if st == nil || st.Occupied {
		t.Fatalf("STACK: expected cleared after signed force-clear, got %+v", st)
	}
}

func TestForceClearUsesCallerSource(t *testing.T) {
	m := bitmap.NewBlockMap()
	d := conflict.NewDetector(m, "test-secret-key-123")
	src := bitmap.SourceID("station-a")
	occ, _ := codec.NewFrame(9, 1, 9, 0)
	d.ApplyFrame(occ, src, "")

	f, _ := codec.NewFrame(9, 0, 10, codec.FlagForceClear)
	sig := d.Signer().Sign(9, src, 10)
	_, st := d.ApplyFrame(f, src, sig)
	if st == nil || st.Occupied {
		t.Fatalf("force-clear must remove source %q, got %+v", src, st)
	}
	if _, ok := st.Sources[src]; ok {
		t.Fatalf("source %q still present after clear: %+v", src, st.Sources)
	}
}

func TestR1IgnoresFreeSourcesInCount(t *testing.T) {
	m := bitmap.NewBlockMap()
	d := conflict.NewDetector(m, "test-secret-key-123")

	f1, _ := codec.NewFrame(12, 1, 1, 0)
	d.ApplyFrame(f1, "station-a", "")
	fFree, _ := codec.NewFrame(12, 0, 2, 0)
	d.ApplyFrame(fFree, "station-a", "")
	f2, _ := codec.NewFrame(12, 1, 3, 0)
	res, st := d.ApplyFrame(f2, "station-b", "")

	if st.OccupiedSourceCount() != 1 {
		t.Fatalf("OccupiedSourceCount=%d want 1 (free sources must not count)", st.OccupiedSourceCount())
	}
	if res.HasConflict {
		t.Fatalf("R1 false positive when prior source is free: %+v", res)
	}
}
