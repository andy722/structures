package sparse

import (
	"github.com/andy722/structures/offheap"
	"github.com/andy722/structures/range"
	"sort"
)

// SparseRangeStore maps inclusive range [fromIncl, toIncl] to value
type SparseRangeStore struct {
	grow float64

	from, end *offheap.OffHeapArrayUint64
	v1, v2    *offheap.OffHeapArrayUint16
}

func NewSparseRangeStore(initialSize int, grow float64) SparseRangeStore {
	return SparseRangeStore{
		grow: grow,
		from: offheap.NewOffHeapArrayUint64(initialSize),
		end:  offheap.NewOffHeapArrayUint64(initialSize),
		v1:   offheap.NewOffHeapArrayUint16(initialSize),
		v2:   offheap.NewOffHeapArrayUint16(initialSize),
	}
}

func (s *SparseRangeStore) Get(key SparseArrayKey) (v1 uint16, v2 uint16, exists bool) {
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

func (s *SparseRangeStore) ValuesV1(callback func(uint16)) {
	s.v1.Values(callback)
}

func (s *SparseRangeStore) ValuesV2(callback func(uint16)) {
	s.v2.Values(callback)
}

func (s *SparseRangeStore) checkMatch(key SparseArrayKey, idx int) (v1, v2 uint16, exists bool) {
	if rangeStart := s.from.Get(idx); rangeStart > key {
		return
	}

	if rangeEnd := s.end.Get(idx); rangeEnd < key {
		return
	}

	return s.v1.Get(idx), s.v2.Get(idx), true
}

func (s SparseRangeStore) Size() int {
	return s.from.Len()
}

func (s *SparseRangeStore) cap() int {
	return s.from.Cap()
}

func (s *SparseRangeStore) shrink() {
	if size := s.Size(); size < s.cap() {
		s.from = s.from.TrimToSize()
		s.end = s.end.TrimToSize()
		s.v1 = s.v1.TrimToSize()
		s.v2 = s.v2.TrimToSize()
	}
}

func (s *SparseRangeStore) growBackingArraysIfNeeded() {
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

func (s *SparseRangeStore) Close() {
	s.from.Dealloc()
	s.end.Dealloc()
	s.v1.Dealloc()
	s.v2.Dealloc()
}

type SparseRangeStoreBuilder struct {
	s          SparseRangeStore
	shouldSort bool
}

func NewSparseRangeStoreBuilder(initialSize int) SparseRangeStoreBuilder {
	return SparseRangeStoreBuilder{
		s: NewSparseRangeStore(initialSize, 1.25),
	}
}

func (b *SparseRangeStoreBuilder) Add(fromIncl, toIncl _range.RangePoint, v1, v2 uint16) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.from.Append(fromIncl)
	b.s.end.Append(toIncl)
	b.s.v1.Append(v1)
	b.s.v2.Append(v2)
}

func (b *SparseRangeStoreBuilder) Build() SparseRangeStore {
	if b.shouldSort {
		b.sort()
	}

	b.s.shrink()

	return b.s
}

func (b *SparseRangeStoreBuilder) sort() {
	sort.Sort(sparseRangeSorter(func() SparseRangeStore { return b.s }))
	b.shouldSort = false
}

type sparseRangeSorter func() SparseRangeStore

func (s sparseRangeSorter) Len() int {
	return s().Size()
}

func (s sparseRangeSorter) Less(i, j int) bool {
	keys := s().from.slice
	return keys[i] < keys[j]
}

func (s sparseRangeSorter) Swap(i, j int) {
	from := s().from.slice
	end := s().end.slice
	v1 := s().v1.slice
	v2 := s().v2.slice

	from[i], from[j] = from[j], from[i]
	end[i], end[j] = end[j], end[i]
	v1[i], v1[j] = v1[j], v1[i]
	v2[i], v2[j] = v2[j], v2[i]
}
