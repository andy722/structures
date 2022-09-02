package sparse

import (
	"github.com/andy722/structures/offheap"
	"github.com/andy722/structures/range"
	"sort"
)

// RangeStore maps inclusive range [fromIncl, toIncl] to value
type RangeStore struct {
	grow float64

	from, end *offheap.ArrayUint64
	v1, v2    *offheap.ArrayUint16
}

func NewSparseRangeStore(initialSize int, grow float64) RangeStore {
	return RangeStore{
		grow: grow,
		from: offheap.NewArrayUint64(initialSize),
		end:  offheap.NewArrayUint64(initialSize),
		v1:   offheap.NewArrayUint16(initialSize),
		v2:   offheap.NewArrayUint16(initialSize),
	}
}

func (s *RangeStore) Get(key ArrayUint64Key) (v1 uint16, v2 uint16, exists bool) {
	idx := sort.Search(s.Size(), func(i int) bool { return s.from.Get(i) >= key })
	if idx >= s.Size() {
		// Check if the last element matches
		return s.checkMatch(key, idx-1)
	}

	if v1, v2, exists = s.checkMatch(key, idx); exists {
		return
	}

	if idx > 0 {
		return s.checkMatch(key, idx-1)
	}

	return
}

func (s *RangeStore) ValuesV1(callback func(uint16)) {
	s.v1.Values(callback)
}

func (s *RangeStore) ValuesV2(callback func(uint16)) {
	s.v2.Values(callback)
}

func (s *RangeStore) checkMatch(key ArrayUint64Key, idx int) (v1, v2 uint16, exists bool) {
	if rangeStart := s.from.Get(idx); rangeStart > key {
		return
	}

	if rangeEnd := s.end.Get(idx); rangeEnd < key {
		return
	}

	return s.v1.Get(idx), s.v2.Get(idx), true
}

func (s *RangeStore) Size() int {
	return s.from.Len()
}

func (s *RangeStore) cap() int {
	return s.from.Cap()
}

func (s *RangeStore) shrink() {
	if size := s.Size(); size < s.cap() {
		s.from = s.from.TrimToSize()
		s.end = s.end.TrimToSize()
		s.v1 = s.v1.TrimToSize()
		s.v2 = s.v2.TrimToSize()
	}
}

func (s *RangeStore) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.from = s.from.Grow(newSize)
	s.end = s.end.Grow(newSize)
	s.v1 = s.v1.Grow(newSize)
	s.v2 = s.v2.Grow(newSize)
}

func (s *RangeStore) Close() {
	s.from.Dealloc()
	s.end.Dealloc()
	s.v1.Dealloc()
	s.v2.Dealloc()
}

type RangeStoreBuilder struct {
	s          RangeStore
	shouldSort bool
}

//goland:noinspection GoUnusedExportedFunction
func NewRangeStoreBuilder(initialSize int) RangeStoreBuilder {
	return RangeStoreBuilder{
		s: NewSparseRangeStore(initialSize, 1.25),
	}
}

func (b *RangeStoreBuilder) Add(fromIncl, toIncl _range.RangePoint, v1, v2 uint16) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.from.Append(fromIncl)
	b.s.end.Append(toIncl)
	b.s.v1.Append(v1)
	b.s.v2.Append(v2)
}

func (b *RangeStoreBuilder) Build() RangeStore {
	if b.shouldSort {
		b.sort()
	}

	b.s.shrink()

	return b.s
}

func (b *RangeStoreBuilder) sort() {
	sort.Sort(sparseRangeSorter(func() *RangeStore { return &b.s }))
	b.shouldSort = false
}

type sparseRangeSorter func() *RangeStore

func (s sparseRangeSorter) Len() int {
	return s().Size()
}

func (s sparseRangeSorter) Less(i, j int) bool {
	keys := s().from
	return keys.Get(i) < keys.Get(j)
}

func (s sparseRangeSorter) Swap(i, j int) {
	s().from.Swap(i, j)
	s().end.Swap(i, j)
	s().v1.Swap(i, j)
	s().v2.Swap(i, j)
}
