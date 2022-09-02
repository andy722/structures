package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMnoMask_ParseMask(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		masks := []string{
			"79?????????",
			"7??????????",
			"79441112233",
			"1",
			"?",
		}
		for _, v := range masks {
			_, err := ParseMask(v)
			assert.NoError(t, err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		masks := []string{
			"abc",
			"",
			"def",
			"??a?",
		}
		for _, v := range masks {
			_, err := ParseMask(v)
			assert.Error(t, err)
		}

		assert.Panics(t, func() {
			MustParseMask("abc")
		})
	})
}

func TestMnoMask_Digit(t *testing.T) {
	mask := MustParseMask("79?????????")

	assert.Equal(t, byte(7), mask.Digit(0))
	assert.Equal(t, byte(9), mask.Digit(1))
	assert.True(t, mask.IsWildcard(10))
}

func BenchmarkMnoMask_Digit(b *testing.B) {
	mask := MustParseMask("79?????????")

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		mask.Digit(1)
	}
}
