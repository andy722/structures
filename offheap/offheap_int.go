package offheap

import "unsafe"

// OffHeapArrayInt is an off-heap array if uint64 values
type OffHeapArrayInt struct {
	offHeapArray
	slice []int
}

func NewOffHeapArrayInt(size int) *OffHeapArrayInt {
	var kEl int

	base := offHeapArray{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &OffHeapArrayInt{
		offHeapArray: base,
		slice:        *(*[]int)(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *OffHeapArrayInt) Len() int {
	return len(o.slice)
}

func (o *OffHeapArrayInt) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *OffHeapArrayInt) Insert(i int, v int) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *OffHeapArrayInt) Get(i int) int {
	return o.slice[i]
}

func (o *OffHeapArrayInt) Set(i int, val int) {
	o.slice[i] = val
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *OffHeapArrayInt) Append(v int) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *OffHeapArrayInt) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
	// o.slice = append(o.slice[:i], o.slice[i+1:]...)
}

func (o *OffHeapArrayInt) Grow(size int) *OffHeapArrayInt {
	target := NewOffHeapArrayInt(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *OffHeapArrayInt) TrimToSize() *OffHeapArrayInt {
	target := NewOffHeapArrayInt(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}
