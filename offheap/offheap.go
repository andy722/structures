package offheap

import (
	"reflect"
	"syscall"
	"unsafe"
)

type offHeapArray struct {
	sz  uintptr
	cap int
}

func (o offHeapArray) Cap() int {
	return o.cap
}

func (o offHeapArray) allocSlice(len int) reflect.SliceHeader {
	noFd := -1

	data, _, errno := syscall.Syscall6(
		syscall.SYS_MMAP,
		0,
		uintptr(len)*o.sz,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE,
		uintptr(noFd),
		0,
	)
	if errno != 0 {
		panic(errno)
	}

	//if log.IsLevelEnabled(log.TraceLevel) {
	//	log.Tracef("allocSlice(%d bytes): %d", uintptr(len)*o.sz, data)
	//}

	return reflect.SliceHeader{
		Data: data,
		Len:  0,
		Cap:  len,
	}
}

func (o offHeapArray) deallocSlice(p unsafe.Pointer) {
	hdr := (*reflect.SliceHeader)(p)

	_, _, errno := syscall.Syscall(
		syscall.SYS_MUNMAP,
		hdr.Data,
		o.sz*uintptr(o.cap),
		0,
	)

	//if log.IsLevelEnabled(log.TraceLevel) {
	//	log.Tracef("deallocSlice(%d bytes): %d", o.sz*uintptr(o.cap), hdr.Data)
	//}

	if errno != 0 {
		panic(errno)
	}
}
