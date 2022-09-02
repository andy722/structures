package offheap

import (
	"golang.org/x/tools/container/intsets"
	"unsafe"
)

type ArrayUint32Value = uint32

// ArrayUint32 is an off-heap array if uint64 values
type ArrayUint32 struct {
	array
	slice []ArrayUint32Value
}

func NewArrayUint32(size int) *ArrayUint32 {
	var kEl ArrayUint32Value

	base := array{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &ArrayUint32{
		array: base,
		slice: *(*[]uint32)(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *ArrayUint32) Len() int {
	return len(o.slice)
}

func (o *ArrayUint32) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *ArrayUint32) Insert(i int, v ArrayUint32Value) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *ArrayUint32) Get(i int) ArrayUint32Value {
	return o.slice[i]
}

func (o *ArrayUint32) Set(i int, val ArrayUint32Value) {
	o.slice[i] = val
}

func (o *ArrayUint32) Swap(i, j int) {
	slice := o.slice
	slice[i], slice[j] = slice[j], slice[i]
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *ArrayUint32) Append(v ArrayUint32Value) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *ArrayUint32) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
}

func (o *ArrayUint32) Grow(size int) *ArrayUint32 {
	target := NewArrayUint32(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *ArrayUint32) TrimToSize() *ArrayUint32 {
	target := NewArrayUint32(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}

func (o *ArrayUint32) Values(callback func(ArrayUint32Value)) {
	var uniq intsets.Sparse
	for i := 0; i < o.Len(); i++ {
		v := o.Get(i)
		if uniq.Insert(int(v)) {
			callback(v)
		}
	}
}
