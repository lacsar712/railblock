package bitmap_test

import (
	"testing"

	"github.com/lacsar712/railblock/internal/bitmap"
)

func TestSetGet(t *testing.T) {
	m := bitmap.NewBlockMap()
	m.Set(12, true, "station-a", 1)
	st := m.Get(12)
	if !st.Occupied {
		t.Fatal("expected occupied")
	}
	if !st.Sources["station-a"] {
		t.Fatal("expected source occupied")
	}
}

func TestDualSourceOccupancy(t *testing.T) {
	m := bitmap.NewBlockMap()
	m.Set(12, true, "station-a", 1)
	m.Set(12, true, "station-b", 2)
	st := m.Get(12)
	if st.OccupiedSourceCount() < 2 {
		t.Fatalf("sources=%d", st.OccupiedSourceCount())
	}
}

func TestForceClear(t *testing.T) {
	m := bitmap.NewBlockMap()
	m.Set(7, true, "a", 1)
	m.ForceClear(7, "a", 2)
	st := m.Get(7)
	if st.Occupied {
		t.Fatal("expected free after clear")
	}
}

func TestListSorted(t *testing.T) {
	m := bitmap.NewBlockMap()
	m.Set(3, false, "x", 1)
	m.Set(1, true, "x", 2)
	blocks := m.List()
	if len(blocks) != 2 {
		t.Fatalf("len=%d", len(blocks))
	}
	if blocks[0].BlockID > blocks[1].BlockID {
		t.Fatal("not sorted")
	}
}

func TestOccupiedSourceCountSkipsFree(t *testing.T) {
	m := bitmap.NewBlockMap()
	m.Set(12, true, "station-a", 1)
	m.Set(12, false, "station-a", 2)
	m.Set(12, true, "station-b", 3)
	st := m.Get(12)
	if st.OccupiedSourceCount() != 1 {
		t.Fatalf("OccupiedSourceCount=%d want 1", st.OccupiedSourceCount())
	}
}

func TestForceClearRecomputesOccupied(t *testing.T) {
	m := bitmap.NewBlockMap()
	m.Set(7, true, "a", 1)
	m.Set(7, true, "b", 2)
	st := m.ForceClear(7, "a", 3)
	if !st.Occupied {
		t.Fatal("block should stay occupied while source b remains")
	}
	st = m.ForceClear(7, "b", 4)
	if st.Occupied {
		t.Fatalf("Occupied aggregate must recompute after last source cleared, got %+v", st)
	}
}
