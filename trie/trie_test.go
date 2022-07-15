package trie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie(t *testing.T) {
	s := NewTrie()

	s.Put(MustParseMask("7944???????"), "a")
	s.Put(MustParseMask("79440??????"), "b")
	s.Put(MustParseMask("79440001103"), "c")

	fmt.Println(s.String())

	assert.Equal(t, "a", s.Lookup(79441001101))
	assert.Equal(t, "b", s.Lookup(79440001101))
	assert.Equal(t, "b", s.Lookup(79440001102))
	assert.Equal(t, "b", s.Lookup(79440000001))
	assert.Equal(t, "c", s.Lookup(79440001103))
}

func BenchmarkTrie_Lookup(b *testing.B) {
	s := NewTrie()
	s.Put(MustParseMask("7944???????"), "a")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s.Lookup(79441001101)
	}
}

func BenchmarkTrie_Put(b *testing.B) {
	s := NewTrie()

	masks := make([]Mask, 0, b.N)
	for i := 0; i < b.N; i++ {
		masks = append(masks, MustParseMask(fmt.Sprintf("%d", 79000000000+i+(i*1000)+(i*100))))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s.Put(masks[i], "s")
	}
}
