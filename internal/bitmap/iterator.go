package bitmap

// Iterator walks block states without exposing internal locks for long periods.
type Iterator struct {
	mapRef *BlockMap
	index  int
	items  []*BlockState
}

// NewIterator snapshots the current block list for sequential traversal.
func (m *BlockMap) NewIterator() *Iterator {
	return &Iterator{mapRef: m, items: m.List()}
}

// Next advances the iterator and returns the next block, or nil when exhausted.
func (it *Iterator) Next() *BlockState {
	if it.index >= len(it.items) {
		return nil
	}
	st := it.items[it.index]
	it.index++
	return st
}

// Remaining reports how many blocks are left to visit.
func (it *Iterator) Remaining() int {
	if it.index >= len(it.items) {
		return 0
	}
	return len(it.items) - it.index
}

// Reset restarts iteration from the first block.
func (it *Iterator) Reset() {
	it.index = 0
	it.items = it.mapRef.List()
}

// FilterOccupied returns occupied blocks from the iterator's snapshot.
func (it *Iterator) FilterOccupied() []*BlockState {
	var out []*BlockState
	for _, st := range it.items {
		if st.Occupied {
			out = append(out, st)
		}
	}
	return out
}

// FilterBySource returns blocks where the given source reports occupied.
func (it *Iterator) FilterBySource(source SourceID) []*BlockState {
	var out []*BlockState
	for _, st := range it.items {
		if st.Sources[source] {
			out = append(out, st)
		}
	}
	return out
}

// ForEach invokes fn for each block in the snapshot until fn returns false.
func (it *Iterator) ForEach(fn func(*BlockState) bool) {
	for _, st := range it.items {
		if !fn(st) {
			return
		}
	}
}
