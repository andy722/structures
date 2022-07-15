package _range

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	s := NewRangeStore()

	s.Add(1, 2, "1")
	s.Add(3, 4, "3")
	s.Add(5, 6, "4")

	assert.Equal(t, "1", s.Lookup(2))
	assert.Equal(t, "3", s.Lookup(3))
	assert.Equal(t, "3", s.Lookup(4))
	assert.Nil(t, s.Lookup(7))
}

func BenchmarkAddSequential(b *testing.B) {
	a := []int{100, 1000, 10_000, 1000_000}
	for _, r := range a {
		b.Run(fmt.Sprintf("range=%d", r), func(b *testing.B) {
			s := NewRangeStore()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				s.Add(RangePoint(r*i), RangePoint(r*i+r), "s")
			}
		})

	}
}

func BenchmarkLookup(b *testing.B) {
	b.ReportAllocs()

	s := NewRangeStore()

	r := 100
	for i := 0; i < 100_000; i++ {
		s.Add(RangePoint(r*i), RangePoint(r*i+r), "s")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Lookup(RangePoint(1000))
	}
}
