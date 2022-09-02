package offheap

import (
	"unsafe"
)

type ArrayUint64Value = uint64

// ArrayUint64 is an off-heap array if uint64 values
type ArrayUint64 struct {
	array
	slice []ArrayUint64Value
}

func NewArrayUint64(size int) *ArrayUint64 {
	var kEl ArrayUint64Value

	base := array{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &ArrayUint64{
		array: base,
		slice: *(*[]uint64)(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *ArrayUint64) Len() int {
	return len(o.slice)
}

func (o *ArrayUint64) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *ArrayUint64) Insert(i int, v ArrayUint64Value) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *ArrayUint64) Get(i int) ArrayUint64Value {
	return o.slice[i]
}

func (o *ArrayUint64) Set(i int, val ArrayUint64Value) {
	o.slice[i] = val
}

func (o *ArrayUint64) Swap(i, j int) {
	slice := o.slice
	slice[i], slice[j] = slice[j], slice[i]
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *ArrayUint64) Append(v ArrayUint64Value) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *ArrayUint64) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
}

func (o *ArrayUint64) Grow(size int) *ArrayUint64 {
	target := NewArrayUint64(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *ArrayUint64) TrimToSize() *ArrayUint64 {
	target := NewArrayUint64(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}
