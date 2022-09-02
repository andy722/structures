package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

const NoValue int = -1

// ArrayInt provides an off-heap map with numeric keys, internally represented as sparse array
type ArrayInt struct {
	arrayUint64

	values *offheap.ArrayInt
}

func NewSparseArrayInt(preallocate int, grow float64) *ArrayInt {
	return &ArrayInt{
		arrayUint64{
			preallocate,
			grow,
			offheap.NewArrayUint64(preallocate),
		},
		offheap.NewArrayInt(preallocate),
	}
}

func (s *ArrayInt) Close() {
	s.arrayUint64.Close()
	s.values.Dealloc()
}

func (s *ArrayInt) Add(key ArrayUint64Key, val int) {
	i := s.idx(key)
	if i < s.keys.Len() && s.keys.Get(i) == key {
		s.values.Set(i, val)
		return
	}

	s.growBackingArraysIfNeeded()

	s.keys.Insert(i, key)
	s.values.Insert(i, val)
}

func (s *ArrayInt) Get(key ArrayUint64Key) int {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		return s.values.Get(i)
	}
	return NoValue
}

func (s *ArrayInt) Delete(key ArrayUint64Key) (prev int) {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		prev = s.values.Get(i)
		s.values.Set(i, NoValue)
	}
	return
}

func (s *ArrayInt) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.keys = s.keys.Grow(newSize)
	s.values = s.values.Grow(newSize)
}

func (s *ArrayInt) cleanup() {
	for i := s.Size() - 1; i >= 0; i-- {
		if s.values.Get(i) == NoValue {
			s.keys.Remove(i)
			s.values.Remove(i)
		}
	}
}

func (s *ArrayInt) shrink() {
	if size := s.Size(); size < s.cap() {
		s.keys = s.keys.TrimToSize()
		s.values = s.values.TrimToSize()
	}
}

type ArrayIntBuilder struct {
	s             *ArrayInt
	shouldSort    bool // Marks as containing non-sorted data, need to sort prior to lookups
	shouldCleanup bool // Marks as containing gaps, i.e. deleted entries
}

//goland:noinspection GoUnusedExportedFunction
func NewArrayIntBuilder(preallocate int, grow float64) *ArrayIntBuilder {
	return &ArrayIntBuilder{
		s: NewSparseArrayInt(preallocate, grow),
	}
}

func (b *ArrayIntBuilder) Add(key ArrayUint64Key, value int) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.keys.Append(key)
	b.s.values.Append(value)
}

func (b *ArrayIntBuilder) Delete(key ArrayUint64Key) {
	if b.shouldSort {
		b.sort()
	}

	if b.s.Delete(key) != NoValue {
		b.shouldCleanup = true
	}
}

func (b *ArrayIntBuilder) Build() *ArrayInt {
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

func (b *ArrayIntBuilder) sort() {
	sort.Sort(sparseArrayIntSorter(func() *ArrayInt { return b.s }))
	b.shouldSort = false
}

type sparseArrayIntSorter func() *ArrayInt

func (s sparseArrayIntSorter) Len() int {
	return s().Size()
}

func (s sparseArrayIntSorter) Less(i, j int) bool {
	keys := s().keys
	return keys.Get(i) < keys.Get(j)
}

func (s sparseArrayIntSorter) Swap(i, j int) {
	keys := s().keys
	tmp := keys.Get(i)
	keys.Set(i, keys.Get(j))
	keys.Set(j, tmp)

	values := s().values
	tmp1 := values.Get(i)
	values.Set(i, values.Get(j))
	values.Set(j, tmp1)
}
