package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

const DefaultPreallocate = 17_000_000
const DefaultGrow = 1.25

type SparseArrayKey = uint64
type SparseArrayUint32Key = uint32

// SparseArray provides an off-heap map with numeric keys, internally represented as sparse array
type sparseArray struct {
	preallocate int
	grow        float64

	keys *offheap.OffHeapArrayUint64
}

type sparseArrayUint32 struct {
	preallocate int
	grow        float64

	keys *offheap.OffHeapArrayUint32
}

// SparseArray provides an off-heap map with numeric keys, internally represented as sparse array
type SparseArray struct {
	sparseArray

	values *offheap.OffHeapArrayInterface
}

func NewSparseArray(preallocate int, grow float64) *SparseArray {
	return &SparseArray{
		sparseArray{
			preallocate,
			grow,
			offheap.NewOffHeapArrayUint64(preallocate),
		},
		offheap.NewOffHeapArrayInterface(preallocate),
	}
}

func (s *sparseArray) Size() int {
	return s.keys.Len()
}

func (s *sparseArray) Close() {
	s.keys.Dealloc()
}

func (s *sparseArrayUint32) Size() int {
	return s.keys.Len()
}

func (s *sparseArrayUint32) Close() {
	s.keys.Dealloc()
}

func (s *SparseArray) Close() {
	s.sparseArray.Close()
	s.values.Dealloc()
}

func (s *sparseArray) cap() int {
	return s.keys.Cap()
}

func (s *sparseArrayUint32) cap() int {
	return s.keys.Cap()
}

func (s *SparseArray) Add(key SparseArrayKey, val interface{}) {
	i := s.idx(key)
	if i < s.keys.Len() && s.keys.Get(i) == key {
		s.values.Set(i, val)
		return
	}

	s.growBackingArraysIfNeeded()

	s.keys.Insert(i, key)
	s.values.Insert(i, val)
}

func (s *SparseArray) Get(key SparseArrayKey) interface{} {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		return s.values.Get(i)
	}
	return nil
}

func (s *SparseArray) Delete(key SparseArrayKey) (prev interface{}) {
	if i := s.idx(key); i < s.Size() && s.keys.Get(i) == key {
		prev = s.values.Get(i)
		s.values.Set(i, nil)
	}
	return
}

func (s *sparseArray) idx(key SparseArrayKey) int {
	return sort.Search(s.Size(), func(i int) bool { return s.keys.Get(i) >= key })
}

func (s *sparseArrayUint32) idx(key SparseArrayUint32Key) int {
	return sort.Search(s.Size(), func(i int) bool { return s.keys.Get(i) >= key })
}

func (s *SparseArray) growBackingArraysIfNeeded() {
	size := s.Size()
	if size < s.cap() {
		return
	}

	newSize := int(s.grow * float64(size))

	s.keys = s.keys.Grow(newSize)
	s.values = s.values.Grow(newSize)
}

func (s *SparseArray) cleanup() {
	for i := s.Size() - 1; i >= 0; i-- {
		if s.values.Get(i) == nil {
			s.keys.Remove(i)
			s.values.Remove(i)
		}
	}
}

func (s *SparseArray) shrink() {
	if size := s.Size(); size < s.cap() {
		s.keys = s.keys.TrimToSize()
		s.values = s.values.TrimToSize()
	}
}

type SparseArrayBuilder struct {
	s             *SparseArray
	shouldSort    bool // Marks as containing non-sorted data, need to sort prior to lookups
	shouldCleanup bool // Marks as containing gaps, i.e. deleted entries
}

func NewSparseArrayBuilder() *SparseArrayBuilder {
	return NewSparseArrayBuilder1(DefaultPreallocate, DefaultGrow)
}

func NewSparseArrayBuilder1(preallocate int, grow float64) *SparseArrayBuilder {
	return &SparseArrayBuilder{
		s: NewSparseArray(preallocate, grow),
	}
}

func (b *SparseArrayBuilder) Add(key SparseArrayKey, value interface{}) {
	b.shouldSort = true

	b.s.growBackingArraysIfNeeded()

	b.s.keys.Append(key)
	b.s.values.Append(value)
}

func (b *SparseArrayBuilder) Delete(key SparseArrayKey) {
	if b.shouldSort {
		b.sort()
	}

	if b.s.Delete(key) != nil {
		b.shouldCleanup = true
	}
}

func (b *SparseArrayBuilder) Build() *SparseArray {
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

func (b *SparseArrayBuilder) sort() {
	sort.Sort(sparseArraySorter(func() *SparseArray { return b.s }))
	b.shouldSort = false
}

type sparseArraySorter func() *SparseArray

func (s sparseArraySorter) Len() int {
	return s().Size()
}

func (s sparseArraySorter) Less(i, j int) bool {
	keys := s().keys.slice
	return keys[i] < keys[j]
}

func (s sparseArraySorter) Swap(i, j int) {
	keys := s().keys.slice
	values := s().values.slice

	keys[i], keys[j] = keys[j], keys[i]
	values[i], values[j] = values[j], values[i]
}
