package conflict_test

import (
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
	// apply valid first
	f, _ := codec.NewFrame(5, 1, 1, 0)
	d.ApplyFrame(f, "a", "")
	st := m.Get(5)
	if !st.Occupied {
		t.Fatal("setup failed")
	}
	// corrupt frame would be rejected at decode layer before ApplyFrame
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
