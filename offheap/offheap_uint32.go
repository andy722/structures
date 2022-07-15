package offheap

import "unsafe"

type OffHeapArrayUint32Value = uint32

// OffHeapArrayUint32 is an off-heap array if uint64 values
type OffHeapArrayUint32 struct {
	offHeapArray
	slice []OffHeapArrayUint32Value
}

func NewOffHeapArrayUint32(size int) *OffHeapArrayUint32 {
	var kEl OffHeapArrayUint32Value

	base := offHeapArray{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &OffHeapArrayUint32{
		offHeapArray: base,
		slice:        *(*[]uint32)(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *OffHeapArrayUint32) Len() int {
	return len(o.slice)
}

func (o *OffHeapArrayUint32) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *OffHeapArrayUint32) Insert(i int, v OffHeapArrayUint32Value) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *OffHeapArrayUint32) Get(i int) OffHeapArrayUint32Value {
	return o.slice[i]
}

func (o *OffHeapArrayUint32) Set(i int, val OffHeapArrayUint32Value) {
	o.slice[i] = val
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *OffHeapArrayUint32) Append(v OffHeapArrayUint32Value) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *OffHeapArrayUint32) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
}

func (o *OffHeapArrayUint32) Grow(size int) *OffHeapArrayUint32 {
	target := NewOffHeapArrayUint32(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *OffHeapArrayUint32) TrimToSize() *OffHeapArrayUint32 {
	target := NewOffHeapArrayUint32(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}
