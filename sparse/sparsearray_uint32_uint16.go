package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

// ArrayUint32Uint16 provides an off-heap map with numeric keys, internally represented as sparse array
type ArrayUint32Uint16 struct {
	arrayUint32

	values *offheap.ArrayUint16
	size   int
}

func NewArrayUint32Uint16(preallocate int, grow float64) *ArrayUint32Uint16 {
	return &ArrayUint32Uint16{
		arrayUint32{
			preallocate,
			grow,
			offheap.NewArrayUint32(preallocate),
		},
		offheap.NewArrayUint16(preallocate),
		0,
	}
}

func (s *ArrayUint32Uint16) Close() {
	s.arrayUint32.Close()
	s.values.Dealloc()
}

func (s *ArrayUint32Uint16) Get(key ArrayUint32Key) offheap.ArrayUint16Value {
	if i := s.idx(key); i < s.size && s.keys.Get(i) == key {
		return s.values.Get(i)
	}
	return ArrayUint16NoValue
}

func (s *ArrayUint32Uint16) Values(callback func(uint16)) {
	s.values.Values(callback)
}

func (s *ArrayUint32Uint16) idx(key ArrayUint32Key) int {
	return sort.Search(s.size, func(i int) bool { return s.keys.Get(i) >= key })
}

func (s *ArrayUint32Uint16) Delete(key ArrayUint32Key) (prev offheap.ArrayUint16Value) {
	if i := s.idx(key); i < s.size && s.keys.Get(i) == key {
		prev = s.values.Get(i)
		s.values.Set(i, ArrayUint16NoValue)
	}
	return
}

func (s *ArrayUint32Uint16) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.keys = s.keys.Grow(newSize)
	s.values = s.values.Grow(newSize)
}

func (s *ArrayUint32Uint16) cleanup() {
	for i := s.Size() - 1; i >= 0; i-- {
		if s.values.Get(i) == ArrayUint16NoValue {
			s.keys.Remove(i)
			s.values.Remove(i)
		}
	}
}

func (s *ArrayUint32Uint16) shrink() {
	if size := s.Size(); size < s.cap() {
		s.keys = s.keys.TrimToSize()
		s.values = s.values.TrimToSize()
	}
}

type ArrayUint32Uint16Builder struct {
	s             *ArrayUint32Uint16
	shouldSort    bool // Marks as containing non-sorted data, need to sort prior to lookups
	shouldCleanup bool // Marks as containing gaps, i.e. deleted entries
}

//goland:noinspection GoUnusedExportedFunction
func NewArrayUint32Uint16Builder() *ArrayUint32Uint16Builder {
	return NewArrayUint32Uint16Builder1(DefaultPreallocate, DefaultGrow)
}

func NewArrayUint32Uint16Builder1(preallocate int, grow float64) *ArrayUint32Uint16Builder {
	return &ArrayUint32Uint16Builder{
		s: NewArrayUint32Uint16(preallocate, grow),
	}
}

func (b *ArrayUint32Uint16Builder) Add(key ArrayUint32Key, value offheap.ArrayUint16Value) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.keys.Append(key)
	b.s.values.Append(value)
}

func (b *ArrayUint32Uint16Builder) Delete(key ArrayUint32Key) {
	if b.shouldSort {
		b.sort()
	}

	if b.s.Delete(key) != ArrayUint16NoValue {
		b.shouldCleanup = true
	}
}

func (b *ArrayUint32Uint16Builder) Build() *ArrayUint32Uint16 {
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

func (b *ArrayUint32Uint16Builder) sort() {
	sort.Sort(ArrayUint32Uint16Sorter(func() *ArrayUint32Uint16 { return b.s }))
	b.shouldSort = false
}

type ArrayUint32Uint16Sorter func() *ArrayUint32Uint16

func (s ArrayUint32Uint16Sorter) Len() int {
	return s().Size()
}

func (s ArrayUint32Uint16Sorter) Less(i, j int) bool {
	keys := s().keys
	return keys.Get(i) < keys.Get(j)
}

func (s ArrayUint32Uint16Sorter) Swap(i, j int) {
	s().keys.Swap(i, j)
	s().values.Swap(i, j)
}
