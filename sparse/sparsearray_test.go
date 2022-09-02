package sparse

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparseArray_Add_Sequential(t *testing.T) {
	s := NewSparseArray(DefaultPreallocate, DefaultGrow)

	s.Add(1, "1")
	s.Add(2, "2")
	s.Add(3, "3")

	assert.Equal(t, 3, s.Size())

	assert.Equal(t, "1", s.Get(1))
	assert.Equal(t, "2", s.Get(2))
	assert.Equal(t, "3", s.Get(3))
}

func TestSparseArray_Add_Simple(t *testing.T) {
	s := NewSparseArray(DefaultPreallocate, DefaultGrow)

	s.Add(1, "1")
	s.Add(5, "5")
	s.Add(2, "2")
	s.Add(100, "100")
	s.Add(3, "3")

	assert.Equal(t, "100", s.Get(100))
	assert.Equal(t, "1", s.Get(1))
}

func TestSparseArrayBuilder_Delete(t *testing.T) {
	b := NewArrayInterfaceBuilder()

	b.Add(1, "1")
	b.Add(5, "5")
	b.Add(2, "2")
	b.Add(100, "100")
	b.Add(3, "3")

	b.Delete(2)
	b.Delete(1)

	s := b.Build()

	assert.Nil(t, s.Get(2))
	assert.Nil(t, s.Get(1))
	assert.Equal(t, "5", s.Get(5))
}

func TestSparseArrayBuilder_Add(t *testing.T) {
	n := 5000

	items := pseudoRandomArray(n)
	b := NewArrayInterfaceBuilder()
	for i, v := range items {
		b.Add(ArrayInterfaceKey(v), i)
	}

	s := b.Build()
	for i, v := range items {
		assert.Equal(t, i, s.Get(ArrayInterfaceKey(v)))
	}
}

func BenchmarkSparseArrayBuilder_Add(b *testing.B) {
	s := NewArrayInterfaceBuilder()
	items := pseudoRandomArray(b.N)

	b.ResetTimer()
	b.ReportAllocs()

	for i, v := range items {
		s.Add(ArrayInterfaceKey(v), i)
	}

	_ = s.Build()
}

func BenchmarkSparseArrayBuilder_Delete(b *testing.B) {
	items := pseudoRandomArray(b.N + 1)

	s := NewArrayInterfaceBuilder()
	for i, v := range items {
		s.Add(ArrayInterfaceKey(v), i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for _, v := range items {
		s.Delete(ArrayInterfaceKey(v))
	}

	_ = s.Build()
}

func BenchmarkSparseArrayBuilder_Lookup(b *testing.B) {
	s := NewArrayInterfaceBuilder()
	items := pseudoRandomArray(b.N)

	for i, v := range items {
		s.Add(ArrayInterfaceKey(v), i)
	}

	a := s.Build()

	b.ResetTimer()
	b.ReportAllocs()

	for _, v := range items {
		a.Get(ArrayInterfaceKey(v))
	}
}

func pseudoRandomArray(size int) (rc []int) {
	rc = make([]int, size)
	for i := range rc {
		rc[i] = i
	}

	rand.Seed(42)
	rand.Shuffle(len(rc), func(i, j int) { rc[i], rc[j] = rc[j], rc[i] })
	return rc
}
