package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

// arrayUint32 provides an off-heap map with numeric keys, internally represented as sparse array
type arrayUint32 struct {
	preallocate int
	grow        float64

	keys *offheap.ArrayUint32
}

func (s *arrayUint32) Size() int {
	return s.keys.Len()
}

func (s *arrayUint32) Close() {
	s.keys.Dealloc()
}

func (s *arrayUint32) idx(key ArrayUint32Key) int {
	return sort.Search(s.Size(), func(i int) bool { return s.keys.Get(i) >= key })
}

func (s *arrayUint32) cap() int {
	return s.keys.Cap()
}
