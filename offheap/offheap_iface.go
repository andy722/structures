package offheap

import "unsafe"

// OffHeapArrayInterface is an off-heap array if interface{} values
type OffHeapArrayInterface struct {
	offHeapArray
	slice []interface{}
}

func NewOffHeapArrayInterface(size int) *OffHeapArrayInterface {
	var kEl interface{}

	base := offHeapArray{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &OffHeapArrayInterface{
		offHeapArray: base,
		slice:        *(*[]interface{})(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *OffHeapArrayInterface) Len() int {
	return len(o.slice)
}

func (o *OffHeapArrayInterface) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *OffHeapArrayInterface) Insert(i int, v interface{}) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *OffHeapArrayInterface) Get(i int) interface{} {
	return o.slice[i]
}

func (o *OffHeapArrayInterface) Set(i int, val interface{}) {
	o.slice[i] = val
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *OffHeapArrayInterface) Append(v interface{}) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *OffHeapArrayInterface) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
}

func (o *OffHeapArrayInterface) Grow(size int) *OffHeapArrayInterface {
	target := NewOffHeapArrayInterface(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *OffHeapArrayInterface) TrimToSize() *OffHeapArrayInterface {
	target := NewOffHeapArrayInterface(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}
