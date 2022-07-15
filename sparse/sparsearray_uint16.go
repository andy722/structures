package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

const SparseArrayUint16NoValue uint16 = 65535

// SparseArrayUint16 provides an off-heap map with numeric keys, internally represented as sparse array
type SparseArrayUint16 struct {
	sparseArray

	values *offheap.OffHeapArrayUint16
}

func NewSparseArrayUint16(preallocate int, grow float64) *SparseArrayUint16 {
	return &SparseArrayUint16{
		sparseArray{
			preallocate,
			grow,
			offheap.NewOffHeapArrayUint64(preallocate),
		},
		offheap.NewOffHeapArrayUint16(preallocate),
	}
}

func (s *SparseArrayUint16) Close() {
	s.sparseArray.Close()
	s.values.Dealloc()
}

func (s *SparseArrayUint16) Add(key SparseArrayKey, val offheap.OffHeapArrayUint16Value) {
	i := s.idx(key)
	if i < s.keys.Len() && s.keys.Get(i) == key {
		s.values.Set(i, val)
		return
	}

	s.growBackingArraysIfNeeded()

	s.keys.Insert(i, key)
	s.values.Insert(i, val)
}

func (s *SparseArrayUint16) Get(key SparseArrayKey) offheap.OffHeapArrayUint16Value {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		return s.values.Get(i)
	}
	return SparseArrayUint16NoValue
}

func (s *SparseArrayUint16) Delete(key SparseArrayKey) (prev offheap.OffHeapArrayUint16Value) {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		prev = s.values.Get(i)
		s.values.Set(i, SparseArrayUint16NoValue)
	}
	return
}

func (s *SparseArrayUint16) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.keys = s.keys.Grow(newSize)
	s.values = s.values.Grow(newSize)
}

func (s *SparseArrayUint16) cleanup() {
	for i := s.Size() - 1; i >= 0; i-- {
		if s.values.Get(i) == SparseArrayUint16NoValue {
			s.keys.Remove(i)
			s.values.Remove(i)
		}
	}
}

func (s *SparseArrayUint16) shrink() {
	if size := s.Size(); size < s.cap() {
		s.keys = s.keys.TrimToSize()
		s.values = s.values.TrimToSize()
	}
}

type SparseArrayUint16Builder struct {
	s             *SparseArrayUint16
	shouldSort    bool // Marks as containing non-sorted data, need to sort prior to lookups
	shouldCleanup bool // Marks as containing gaps, i.e. deleted entries
}

func NewSparseArrayUint16Builder() *SparseArrayUint16Builder {
	return NewSparseArrayUint16Builder1(DefaultPreallocate, DefaultGrow)
}

func NewSparseArrayUint16Builder1(preallocate int, grow float64) *SparseArrayUint16Builder {
	return &SparseArrayUint16Builder{
		s: NewSparseArrayUint16(preallocate, grow),
	}
}

func (b *SparseArrayUint16Builder) Add(key SparseArrayKey, value offheap.OffHeapArrayUint16Value) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.keys.Append(key)
	b.s.values.Append(value)
}

func (b *SparseArrayUint16Builder) Delete(key SparseArrayKey) {
	if b.shouldSort {
		b.sort()
	}

	if b.s.Delete(key) != SparseArrayUint16NoValue {
		b.shouldCleanup = true
	}
}

func (b *SparseArrayUint16Builder) Build() *SparseArrayUint16 {
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

func (b *SparseArrayUint16Builder) sort() {
	sort.Sort(sparseArrayUint16Sorter(func() *SparseArrayUint16 { return b.s }))
	b.shouldSort = false
}

type sparseArrayUint16Sorter func() *SparseArrayUint16

func (s sparseArrayUint16Sorter) Len() int {
	return s().Size()
}

func (s sparseArrayUint16Sorter) Less(i, j int) bool {
	keys := s().keys.slice
	return keys[i] < keys[j]
}

func (s sparseArrayUint16Sorter) Swap(i, j int) {
	keys := s().keys.slice
	values := s().values.slice

	keys[i], keys[j] = keys[j], keys[i]
	values[i], values[j] = values[j], values[i]
}
