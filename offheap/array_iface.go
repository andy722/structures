package offheap

import "unsafe"

// ArrayInterface is an off-heap array if interface{} values
type ArrayInterface struct {
	array
	slice []interface{}
}

func NewArrayInterface(size int) *ArrayInterface {
	var kEl interface{}

	base := array{unsafe.Sizeof(kEl), size}

	mapped := base.allocSlice(size)
	return &ArrayInterface{
		array: base,
		slice: *(*[]interface{})(unsafe.Pointer(&mapped)), //nolint:govet
	}
}

func (o *ArrayInterface) Len() int {
	return len(o.slice)
}

func (o *ArrayInterface) Dealloc() {
	o.deallocSlice(unsafe.Pointer(&o.slice))
	o.slice = nil
}

func (o *ArrayInterface) Insert(i int, v interface{}) {
	if i == o.Len() {
		o.slice = append(o.slice, v)
		return
	}

	o.slice = append(o.slice[:i+1], o.slice[i:]...)
	o.slice[i] = v
}

func (o *ArrayInterface) Get(i int) interface{} {
	return o.slice[i]
}

func (o *ArrayInterface) Set(i int, val interface{}) {
	o.slice[i] = val
}

func (o *ArrayInterface) Swap(i, j int) {
	slice := o.slice
	slice[i], slice[j] = slice[j], slice[i]
}

// Append add an element to the end. It is a caller's responsibility to Grow() underlying slice if needed.
func (o *ArrayInterface) Append(v interface{}) {
	o.slice = append(o.slice, v)
}

// Remove removes an element at index. It is a caller's responsibility to call TrimToSize() for reclaiming space.
func (o *ArrayInterface) Remove(i int) {
	o.slice[i] = o.slice[o.Len()-1]
	o.slice = o.slice[:o.Len()-1]
}

func (o *ArrayInterface) Grow(size int) *ArrayInterface {
	target := NewArrayInterface(size)
	target.slice = append(target.slice, o.slice...)

	o.Dealloc()
	return target
}

func (o *ArrayInterface) TrimToSize() *ArrayInterface {
	target := NewArrayInterface(len(o.slice))
	target.slice = append(target.slice, o.slice[:len(o.slice)]...)

	o.Dealloc()
	return target
}
