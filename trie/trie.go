package trie

import (
	"fmt"
	"math"
	"strings"
)

type TrieKey = Mask
type TrieLookupKey = uint64

// Digits returns number of digits
func Digits(val TrieLookupKey) (count int) {
	n := val
	for n > 0 {
		n = n / 10
		count++
	}
	return count
}

// Digit returns `n'-th digit of the value where `n' is in range [ 0, value.Digits() )
func Digit(val TrieLookupKey, place int) byte {
	r := val % uint64(math.Pow10(place+1))
	return byte(r / uint64(math.Pow10(place)))
}

// Trie is a simple prefix tree implementation
type Trie struct {
	value    interface{}
	children [10]*Trie // For each digit
	wildcard *Trie
}

func NewTrie() *Trie {
	return new(Trie)
}

func (trie *Trie) String() string {
	var buffer strings.Builder
	trie.string(0, &buffer)
	return buffer.String()
}

func (trie *Trie) string(level int, buf *strings.Builder) {
	orNil := func(p interface{}) interface{} {
		if p != nil {
			return p
		}
		return "(nil)"
	}

	pad := fmt.Sprintf("%*s", level, "")

	if trie.value != nil {
		_, _ = fmt.Fprintf(buf, "%s  %v\n", pad, orNil(trie.value))
	}
	for i, e := range trie.children {
		if e != nil {
			_, _ = fmt.Fprintf(buf, "%s ↳ %d\n", pad, i)
			e.string(level+1, buf)
		}
	}
}

func (trie *Trie) Lookup(key TrieLookupKey) interface{} {
	digits := Digits(key)

	// Try to force stack allocation
	var stack []*Trie
	if digits < 20 {
		stack = make([]*Trie, 0, 2*20)
	} else {
		stack = make([]*Trie, 0, digits)
	}

	stack = append(stack, trie)

	for i := digits - 1; i >= 0; i-- {
		r := Digit(key, i)

		stackLen := len(stack)

		for j := 0; j < stackLen; j++ {
			node := stack[j]

			if byMask := node.children[r]; byMask != nil {
				stack = append(stack, byMask)
			}

			if byWildcard := node.wildcard; byWildcard != nil {
				stack = append(stack, byWildcard)
			}
		}

		stack = stack[stackLen:]
	}

	for _, c := range stack {
		if c.value != nil {
			return c.value
		}
	}

	return nil
}

func (trie *Trie) Put(key TrieKey, value interface{}) bool {
	digits := key.Len()

	node := trie
	for i := 0; i < digits; i++ {

		if key.IsWildcard(i) {
			child := NewTrie()
			node.wildcard, node = child, child

		} else {
			r := key.Digit(i)

			child := node.children[r]
			if child == nil {
				child = NewTrie()
				node.children[r] = child
			}

			node = child
		}
	}

	isNewVal := node.value == nil
	node.value = value
	return isNewVal
}
