package trie

import (
	"fmt"
	"regexp"
)

var maskPattern = regexp.MustCompile(`^[?\d]+$`)

// Mask contains digits and wildcards `?`
type Mask string

func ParseMask(text string) (Mask, error) {
	mask := Mask(text)
	if !maskPattern.MatchString(text) {
		return mask, fmt.Errorf("invalid mask: %s", text)
	}

	return mask, nil
}

func MustParseMask(text string) Mask {
	mask, err := ParseMask(text)
	if err != nil {
		panic(err)
	}

	return mask
}

func (mask Mask) Len() int {
	return len(string(mask))
}

func (mask Mask) Digit(i int) byte {
	return byte(string(mask)[i]) - 48
}

func (mask Mask) IsWildcard(i int) bool {
	return string(mask)[i] == byte('?')
}
