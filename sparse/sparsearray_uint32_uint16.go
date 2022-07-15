package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

// SparseArrayUint32Uint16 provides an off-heap map with numeric keys, internally represented as sparse array
type SparseArrayUint32Uint16 struct {
	sparseArrayUint32

	values *offheap.OffHeapArrayUint16
	size   int
}

func NewSparseArrayUint32Uint16(preallocate int, grow float64) *SparseArrayUint32Uint16 {
	return &SparseArrayUint32Uint16{
		sparseArrayUint32{
			preallocate,
			grow,
			offheap.NewOffHeapArrayUint32(preallocate),
		},
		offheap.NewOffHeapArrayUint16(preallocate),
		0,
	}
}

func (s *SparseArrayUint32Uint16) Close() {
	s.sparseArrayUint32.Close()
	s.values.Dealloc()
}

func (s *SparseArrayUint32Uint16) Get(key SparseArrayUint32Key) offheap.OffHeapArrayUint16Value {
	if i := s.idx(key); i < s.size && s.keys.Get(i) == key {
		return s.values.Get(i)
	}
	return SparseArrayUint16NoValue
}

func (s *SparseArrayUint32Uint16) Values(callback func(uint16)) {
	s.values.Values(callback)
}

func (s *SparseArrayUint32Uint16) idx(key SparseArrayUint32Key) int {
	return sort.Search(s.size, func(i int) bool { return s.keys.Get(i) >= key })
}

func (s *SparseArrayUint32Uint16) Delete(key SparseArrayUint32Key) (prev offheap.OffHeapArrayUint16Value) {
	if i := s.idx(key); i < s.size && s.keys.Get(i) == key {
		prev = s.values.Get(i)
		s.values.Set(i, SparseArrayUint16NoValue)
	}
	return
}

func (s *SparseArrayUint32Uint16) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.keys = s.keys.Grow(newSize)
	s.values = s.values.Grow(newSize)
}

func (s *SparseArrayUint32Uint16) cleanup() {
	for i := s.Size() - 1; i >= 0; i-- {
		if s.values.Get(i) == SparseArrayUint16NoValue {
			s.keys.Remove(i)
			s.values.Remove(i)
		}
	}
}

func (s *SparseArrayUint32Uint16) shrink() {
	if size := s.Size(); size < s.cap() {
		s.keys = s.keys.TrimToSize()
		s.values = s.values.TrimToSize()
	}
}

type SparseArrayUint32Uint16Builder struct {
	s             *SparseArrayUint32Uint16
	shouldSort    bool // Marks as containing non-sorted data, need to sort prior to lookups
	shouldCleanup bool // Marks as containing gaps, i.e. deleted entries
}

func NewSparseArrayUint32Uint16Builder() *SparseArrayUint32Uint16Builder {
	return NewSparseArrayUint32Uint16Builder1(DefaultPreallocate, DefaultGrow)
}

func NewSparseArrayUint32Uint16Builder1(preallocate int, grow float64) *SparseArrayUint32Uint16Builder {
	return &SparseArrayUint32Uint16Builder{
		s: NewSparseArrayUint32Uint16(preallocate, grow),
	}
}

func (b *SparseArrayUint32Uint16Builder) Add(key SparseArrayUint32Key, value offheap.OffHeapArrayUint16Value) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.keys.Append(key)
	b.s.values.Append(value)
}

func (b *SparseArrayUint32Uint16Builder) Delete(key SparseArrayUint32Key) {
	if b.shouldSort {
		b.sort()
	}

	if b.s.Delete(key) != SparseArrayUint16NoValue {
		b.shouldCleanup = true
	}
}

func (b *SparseArrayUint32Uint16Builder) Build() *SparseArrayUint32Uint16 {
	if b.shouldCleanup {
		b.s.cleanup()
		b.shouldCleanup = false
	}

	if b.shouldSort {
		b.sort()
	}

	b.s.shrink()

	b.s.size = b.s.Size()

	return b.s
}

func (b *SparseArrayUint32Uint16Builder) sort() {
	sort.Sort(SparseArrayUint32Uint16Sorter(func() *SparseArrayUint32Uint16 { return b.s }))
	b.shouldSort = false
}

type SparseArrayUint32Uint16Sorter func() *SparseArrayUint32Uint16

func (s SparseArrayUint32Uint16Sorter) Len() int {
	return s().Size()
}

func (s SparseArrayUint32Uint16Sorter) Less(i, j int) bool {
	keys := s().keys.slice
	return keys[i] < keys[j]
}

func (s SparseArrayUint32Uint16Sorter) Swap(i, j int) {
	keys := s().keys.slice
	values := s().values.slice

	keys[i], keys[j] = keys[j], keys[i]
	values[i], values[j] = values[j], values[i]
}
