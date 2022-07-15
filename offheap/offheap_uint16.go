package offheap

import (
	"golang.org/x/tools/container/intsets"
	"unsafe"
)

type OffHeapArrayUint16Value = uint16

// OffHeapArrayUint16 is an off-heap array if uint64 values
type OffHeapArrayUint16 struct {
	offHeapArray
	slice []OffHeapArrayUint16Value
}

func NewOffHeapArrayUint16(size int) *OffHeapArrayUint16 {
	var kEl OffHeapArrayUint16Value

	base := offHeapArray{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &OffHeapArrayUint16{
		offHeapArray: base,
		slice:        *(*[]uint16)(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *OffHeapArrayUint16) Len() int {
	return len(o.slice)
}

func (o *OffHeapArrayUint16) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *OffHeapArrayUint16) Insert(i int, v OffHeapArrayUint16Value) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *OffHeapArrayUint16) Get(i int) OffHeapArrayUint16Value {
	return o.slice[i]
}

func (o *OffHeapArrayUint16) Set(i int, val OffHeapArrayUint16Value) {
	o.slice[i] = val
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *OffHeapArrayUint16) Append(v OffHeapArrayUint16Value) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *OffHeapArrayUint16) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
}

func (o *OffHeapArrayUint16) Grow(size int) *OffHeapArrayUint16 {
	target := NewOffHeapArrayUint16(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *OffHeapArrayUint16) TrimToSize() *OffHeapArrayUint16 {
	target := NewOffHeapArrayUint16(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}

func (o *OffHeapArrayUint16) Values(callback func(uint16)) {
	var uniq intsets.Sparse
	for i := 0; i < o.Len(); i++ {
		v := o.Get(i)
		if uniq.Insert(int(v)) {
			callback(v)
		}
	}
}
