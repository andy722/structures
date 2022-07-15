package offheap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOffHeapArrayInt(t *testing.T) {
	a := NewOffHeapArrayInt(4)

	assert.Equal(t, 0, a.Len())
	assert.Equal(t, 4, a.Cap())

	a.Append(1)
	assert.Equal(t, 1, a.Len())
	assert.Equal(t, 4, a.Cap())

	assert.Equal(t, 1, a.Get(0))

	a.Dealloc()
}
