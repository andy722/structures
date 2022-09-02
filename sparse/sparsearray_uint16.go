package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

const ArrayUint16NoValue uint16 = 65535

// ArrayUint16 provides an off-heap map with numeric keys, internally represented as sparse array
type ArrayUint16 struct {
	arrayUint64

	values *offheap.ArrayUint16
}

func NewSparseArrayUint16(preallocate int, grow float64) *ArrayUint16 {
	return &ArrayUint16{
		arrayUint64{
			preallocate,
			grow,
			offheap.NewArrayUint64(preallocate),
		},
		offheap.NewArrayUint16(preallocate),
	}
}

func (s *ArrayUint16) Close() {
	s.arrayUint64.Close()
	s.values.Dealloc()
}

func (s *ArrayUint16) Add(key ArrayInterfaceKey, val offheap.ArrayUint16Value) {
	i := s.idx(key)
	if i < s.keys.Len() && s.keys.Get(i) == key {
		s.values.Set(i, val)
		return
	}

	s.growBackingArraysIfNeeded()

	s.keys.Insert(i, key)
	s.values.Insert(i, val)
}

func (s *ArrayUint16) Get(key ArrayInterfaceKey) offheap.ArrayUint16Value {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		return s.values.Get(i)
	}
	return ArrayUint16NoValue
}

func (s *ArrayUint16) Delete(key ArrayInterfaceKey) (prev offheap.ArrayUint16Value) {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		prev = s.values.Get(i)
		s.values.Set(i, ArrayUint16NoValue)
	}
	return
}

func (s *ArrayUint16) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.keys = s.keys.Grow(newSize)
	s.values = s.values.Grow(newSize)
}

func (s *ArrayUint16) cleanup() {
	for i := s.Size() - 1; i >= 0; i-- {
		if s.values.Get(i) == ArrayUint16NoValue {
			s.keys.Remove(i)
			s.values.Remove(i)
		}
	}
}

func (s *ArrayUint16) shrink() {
	if size := s.Size(); size < s.cap() {
		s.keys = s.keys.TrimToSize()
		s.values = s.values.TrimToSize()
	}
}

type ArrayUint16Builder struct {
	s             *ArrayUint16
	shouldSort    bool // Marks as containing non-sorted data, need to sort prior to lookups
	shouldCleanup bool // Marks as containing gaps, i.e. deleted entries
}

//goland:noinspection GoUnusedExportedFunction
func NewArrayUint16Builder() *ArrayUint16Builder {
	return NewArrayUint16Builder1(DefaultPreallocate, DefaultGrow)
}

func NewArrayUint16Builder1(preallocate int, grow float64) *ArrayUint16Builder {
	return &ArrayUint16Builder{
		s: NewSparseArrayUint16(preallocate, grow),
	}
}

func (b *ArrayUint16Builder) Add(key ArrayInterfaceKey, value offheap.ArrayUint16Value) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.keys.Append(key)
	b.s.values.Append(value)
}

func (b *ArrayUint16Builder) Delete(key ArrayInterfaceKey) {
	if b.shouldSort {
		b.sort()
	}

	if b.s.Delete(key) != ArrayUint16NoValue {
		b.shouldCleanup = true
	}
}

func (b *ArrayUint16Builder) Build() *ArrayUint16 {
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

func (b *ArrayUint16Builder) sort() {
	sort.Sort(sparseArrayUint16Sorter(func() *ArrayUint16 { return b.s }))
	b.shouldSort = false
}

type sparseArrayUint16Sorter func() *ArrayUint16

func (s sparseArrayUint16Sorter) Len() int {
	return s().Size()
}

func (s sparseArrayUint16Sorter) Less(i, j int) bool {
	keys := s().keys
	return keys.Get(i) < keys.Get(j)
}

func (s sparseArrayUint16Sorter) Swap(i, j int) {
	s().keys.Swap(i, j)
	s().values.Swap(i, j)
}
