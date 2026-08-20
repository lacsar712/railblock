package clearance_test

import (
	"testing"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/clearance"
	"github.com/lacsar712/railblock/internal/codec"
	"github.com/lacsar712/railblock/internal/conflict"
)

func setupChecker(secret string) (*clearance.Checker, *bitmap.BlockMap, *conflict.Detector) {
	m := bitmap.NewBlockMap()
	d := conflict.NewDetector(m, secret)
	inspector := &clearance.ServiceInspector{Map: m, Detector: d}
	return clearance.NewChecker(inspector), m, d
}

func TestAllowWhenFree(t *testing.T) {
	chk, _, _ := setupChecker("test-secret-key-123")
	route := &clearance.Route{RouteID: "r1", Blocks: []uint16{1, 2, 3}}
	res := chk.Check(route)
	if !res.IsAllowed() {
		t.Fatalf("got %+v", res)
	}
}

func TestRejectWhenOccupied(t *testing.T) {
	chk, m, _ := setupChecker("test-secret-key-123")
	m.Set(12, true, "a", 1)
	route := &clearance.Route{RouteID: "r2", Blocks: []uint16{12}}
	res := chk.Check(route)
	if res.IsAllowed() {
		t.Fatal("expected reject")
	}
}

func TestRejectOnConflict(t *testing.T) {
	chk, _, d := setupChecker("test-secret-key-123")
	f1, _ := codec.NewFrame(12, 1, 1, 0)
	d.ApplyFrame(f1, "station-a", "")
	f2, _ := codec.NewFrame(12, 1, 2, 0)
	d.ApplyFrame(f2, "station-b", "")

	route := &clearance.Route{RouteID: "r3", Blocks: []uint16{12}}
	res := chk.Check(route)
	if res.IsAllowed() {
		t.Fatal("expected reject on conflict")
	}
}

func TestRejectEmptyRoute(t *testing.T) {
	chk, _, _ := setupChecker("test-secret-key-123")
	route := &clearance.Route{RouteID: "empty", Blocks: nil}
	res := chk.Check(route)
	if res.IsAllowed() {
		t.Fatal("expected reject")
	}
}

func TestRejectNilRoute(t *testing.T) {
	chk, _, _ := setupChecker("test-secret-key-123")
	res := chk.Check(nil)
	if res.IsAllowed() {
		t.Fatal("expected reject")
	}
}
