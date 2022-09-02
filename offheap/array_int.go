package offheap

import (
	"golang.org/x/tools/container/intsets"
	"unsafe"
)

type ArrayIntValue = int

// ArrayInt is an off-heap array if uint64 values
type ArrayInt struct {
	array
	slice []ArrayIntValue
}

func NewArrayInt(size int) *ArrayInt {
	var kEl ArrayIntValue

	base := array{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &ArrayInt{
		array: base,
		slice: *(*[]int)(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *ArrayInt) Len() int {
	return len(o.slice)
}

func (o *ArrayInt) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *ArrayInt) Insert(i int, v ArrayIntValue) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *ArrayInt) Get(i int) ArrayIntValue {
	return o.slice[i]
}

func (o *ArrayInt) Set(i int, val ArrayIntValue) {
	o.slice[i] = val
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *ArrayInt) Append(v ArrayIntValue) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *ArrayInt) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
	// o.slice = append(o.slice[:i], o.slice[i+1:]...)
}

func (o *ArrayInt) Grow(size int) *ArrayInt {
	target := NewArrayInt(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *ArrayInt) TrimToSize() *ArrayInt {
	target := NewArrayInt(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}

func (o *ArrayInt) Values(callback func(ArrayIntValue)) {
	var uniq intsets.Sparse
	for i := 0; i < o.Len(); i++ {
		v := o.Get(i)
		if uniq.Insert(v) {
			callback(v)
		}
	}
}
