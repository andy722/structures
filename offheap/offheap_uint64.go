package offheap

import "unsafe"

// OffHeapArrayUint64 is an off-heap array if uint64 values
type OffHeapArrayUint64 struct {
	offHeapArray
	slice []uint64
}

func NewOffHeapArrayUint64(size int) *OffHeapArrayUint64 {
	var kEl uint64

	base := offHeapArray{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &OffHeapArrayUint64{
		offHeapArray: base,
		slice:        *(*[]uint64)(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *OffHeapArrayUint64) Len() int {
	return len(o.slice)
}

func (o *OffHeapArrayUint64) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *OffHeapArrayUint64) Insert(i int, v uint64) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *OffHeapArrayUint64) Get(i int) uint64 {
	return o.slice[i]
}

func (o *OffHeapArrayUint64) Set(i int, val uint64) {
	o.slice[i] = val
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *OffHeapArrayUint64) Append(v uint64) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *OffHeapArrayUint64) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
}

func (o *OffHeapArrayUint64) Grow(size int) *OffHeapArrayUint64 {
	target := NewOffHeapArrayUint64(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *OffHeapArrayUint64) TrimToSize() *OffHeapArrayUint64 {
	target := NewOffHeapArrayUint64(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}
