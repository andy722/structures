package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

const NoValue int = -1

// SparseArrayInt provides an off-heap map with numeric keys, internally represented as sparse array
type SparseArrayInt struct {
	sparseArray

	values *offheap.OffHeapArrayInt
}

func NewSparseArrayInt(preallocate int, grow float64) *SparseArrayInt {
	return &SparseArrayInt{
		sparseArray{
			preallocate,
			grow,
			offheap.NewOffHeapArrayUint64(preallocate),
		},
		offheap.NewOffHeapArrayInt(preallocate),
	}
}

func (s *SparseArrayInt) Close() {
	s.sparseArray.Close()
	s.values.Dealloc()
}

func (s *SparseArrayInt) Add(key SparseArrayKey, val int) {
	i := s.idx(key)
	if i < s.keys.Len() && s.keys.Get(i) == key {
		s.values.Set(i, val)
		return
	}

	s.growBackingArraysIfNeeded()

	s.keys.Insert(i, key)
	s.values.Insert(i, val)
}

func (s *SparseArrayInt) Get(key SparseArrayKey) int {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		return s.values.Get(i)
	}
	return NoValue
}

func (s *SparseArrayInt) Delete(key SparseArrayKey) (prev int) {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		prev = s.values.Get(i)
		s.values.Set(i, NoValue)
	}
	return
}

func (s *SparseArrayInt) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.keys = s.keys.Grow(newSize)
	s.values = s.values.Grow(newSize)
}

func (s *SparseArrayInt) cleanup() {
	for i := s.Size() - 1; i >= 0; i-- {
		if s.values.Get(i) == NoValue {
			s.keys.Remove(i)
			s.values.Remove(i)
		}
	}
}

func (s *SparseArrayInt) shrink() {
	if size := s.Size(); size < s.cap() {
		s.keys = s.keys.TrimToSize()
		s.values = s.values.TrimToSize()
	}
}

type SparseArrayIntBuilder struct {
	s             *SparseArrayInt
	shouldSort    bool // Marks as containing non-sorted data, need to sort prior to lookups
	shouldCleanup bool // Marks as containing gaps, i.e. deleted entries
}

func NewSparseArrayIntBuilder(preallocate int, grow float64) *SparseArrayIntBuilder {
	return &SparseArrayIntBuilder{
		s: NewSparseArrayInt(preallocate, grow),
	}
}

func (b *SparseArrayIntBuilder) Add(key SparseArrayKey, value int) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.keys.Append(key)
	b.s.values.Append(value)
}

func (b *SparseArrayIntBuilder) Delete(key SparseArrayKey) {
	if b.shouldSort {
		b.sort()
	}

	if b.s.Delete(key) != NoValue {
		b.shouldCleanup = true
	}
}

func (b *SparseArrayIntBuilder) Build() *SparseArrayInt {
	if b.shouldCleanup {
		b.s.cleanup()
		b.shouldCleanup = false
	}

	if b.shouldSort {
		b.sort()
	}

	b.s.shrink()

	return b.s
}

func (b *SparseArrayIntBuilder) sort() {
	sort.Sort(sparseArrayIntSorter(func() *SparseArrayInt { return b.s }))
	b.shouldSort = false
}

type sparseArrayIntSorter func() *SparseArrayInt

func (s sparseArrayIntSorter) Len() int {
	return s().Size()
}

func (s sparseArrayIntSorter) Less(i, j int) bool {
	keys := s().keys.slice
	return keys[i] < keys[j]
}

func (s sparseArrayIntSorter) Swap(i, j int) {
	keys := s().keys.slice
	values := s().values.slice

	keys[i], keys[j] = keys[j], keys[i]
	values[i], values[j] = values[j], values[i]
}
