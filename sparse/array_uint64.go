package sparse

import (
	"github.com/andy722/structures/offheap"
	"sort"
)

type ArrayUint64Key = uint64

// arrayUint64 provides an off-heap map with numeric keys, internally represented as sparse array
type arrayUint64 struct {
	preallocate int
	grow        float64

	keys *offheap.ArrayUint64
}

func (s *arrayUint64) Size() int {
	return s.keys.Len()
}

func (s *arrayUint64) Close() {
	s.keys.Dealloc()
}

func (s *arrayUint64) idx(key ArrayUint64Key) int {
	return sort.Search(s.Size(), func(i int) bool { return s.keys.Get(i) >= key })
}

func (s *arrayUint64) cap() int {
	return s.keys.Cap()
}
