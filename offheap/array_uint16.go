package offheap

import (
	"golang.org/x/tools/container/intsets"
	"unsafe"
)

type ArrayUint16Value = uint16

// ArrayUint16 is an off-heap array if uint64 values
type ArrayUint16 struct {
	array
	slice []ArrayUint16Value
}

func NewArrayUint16(size int) *ArrayUint16 {
	var kEl ArrayUint16Value

	base := array{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &ArrayUint16{
		array: base,
		slice: *(*[]uint16)(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *ArrayUint16) Len() int {
	return len(o.slice)
}

func (o *ArrayUint16) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *ArrayUint16) Insert(i int, v ArrayUint16Value) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *ArrayUint16) Get(i int) ArrayUint16Value {
	return o.slice[i]
}

func (o *ArrayUint16) Set(i int, val ArrayUint16Value) {
	o.slice[i] = val
}

func (o *ArrayUint16) Swap(i, j int) {
	slice := o.slice
	slice[i], slice[j] = slice[j], slice[i]
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *ArrayUint16) Append(v ArrayUint16Value) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *ArrayUint16) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
}

func (o *ArrayUint16) Grow(size int) *ArrayUint16 {
	target := NewArrayUint16(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *ArrayUint16) TrimToSize() *ArrayUint16 {
	target := NewArrayUint16(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}

func (o *ArrayUint16) Values(callback func(ArrayUint16Value)) {
	var uniq intsets.Sparse
	for i := 0; i < o.Len(); i++ {
		v := o.Get(i)
		if uniq.Insert(int(v)) {
			callback(v)
		}
	}
}
