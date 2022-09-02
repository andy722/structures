package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

const DefaultPreallocate = 17_000_000
const DefaultGrow = 1.25

type ArrayInterfaceKey = uint64
type ArrayUint32Key = uint32

// ArrayInterface provides an off-heap map with numeric keys, internally represented as sparse array
type ArrayInterface struct {
	arrayUint64

	values *offheap.ArrayInterface
}

func NewSparseArray(preallocate int, grow float64) *ArrayInterface {
	return &ArrayInterface{
		arrayUint64{
			preallocate,
			grow,
			offheap.NewArrayUint64(preallocate),
		},
		offheap.NewArrayInterface(preallocate),
	}
}

func (s *ArrayInterface) Close() {
	s.arrayUint64.Close()
	s.values.Dealloc()
}

func (s *ArrayInterface) Add(key ArrayInterfaceKey, val interface{}) {
	i := s.idx(key)
	if i < s.keys.Len() && s.keys.Get(i) == key {
		s.values.Set(i, val)
		return
	}

	s.growBackingArraysIfNeeded()

	s.keys.Insert(i, key)
	s.values.Insert(i, val)
}

func (s *ArrayInterface) Get(key ArrayInterfaceKey) interface{} {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		return s.values.Get(i)
	}
	return nil
}

func (s *ArrayInterface) Delete(key ArrayInterfaceKey) (prev interface{}) {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		prev = s.values.Get(i)
		s.values.Set(i, nil)
	}
	return
}

func (s *ArrayInterface) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.keys = s.keys.Grow(newSize)
	s.values = s.values.Grow(newSize)
}

func (s *ArrayInterface) cleanup() {
	for i := s.Size() - 1; i >= 0; i-- {
		if s.values.Get(i) == nil {
			s.keys.Remove(i)
			s.values.Remove(i)
		}
	}
}

func (s *ArrayInterface) shrink() {
	if size := s.Size(); size < s.cap() {
		s.keys = s.keys.TrimToSize()
		s.values = s.values.TrimToSize()
	}
}

type ArrayInterfaceBuilder struct {
	s             *ArrayInterface
	shouldSort    bool // Marks as containing non-sorted data, need to sort prior to lookups
	shouldCleanup bool // Marks as containing gaps, i.e. deleted entries
}

func NewArrayInterfaceBuilder() *ArrayInterfaceBuilder {
	return NewArrayInterfaceBuilder1(DefaultPreallocate, DefaultGrow)
}

func NewArrayInterfaceBuilder1(preallocate int, grow float64) *ArrayInterfaceBuilder {
	return &ArrayInterfaceBuilder{
		s: NewSparseArray(preallocate, grow),
	}
}

func (b *ArrayInterfaceBuilder) Add(key ArrayInterfaceKey, value interface{}) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.keys.Append(key)
	b.s.values.Append(value)
}

func (b *ArrayInterfaceBuilder) Delete(key ArrayInterfaceKey) {
	if b.shouldSort {
		b.sort()
	}

	if b.s.Delete(key) != nil {
		b.shouldCleanup = true
	}
}

func (b *ArrayInterfaceBuilder) Build() *ArrayInterface {
	if b.shouldCleanup {
		b.s.cleanup()
		b.shouldSort = true
		b.shouldCleanup = false
	}

	if b.shouldSort {
		b.sort()
	}

	b.s.shrink()

	return b.s
}

func (b *ArrayInterfaceBuilder) sort() {
	sort.Sort(sparseArraySorter(func() *ArrayInterface { return b.s }))
	b.shouldSort = false
}

type sparseArraySorter func() *ArrayInterface

func (s sparseArraySorter) Len() int {
	return s().Size()
}

func (s sparseArraySorter) Less(i, j int) bool {
	keys := s().keys
	return keys.Get(i) < keys.Get(j)
}

func (s sparseArraySorter) Swap(i, j int) {
	s().keys.Swap(i, j)
	s().values.Swap(i, j)
}
