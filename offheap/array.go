package offheap

import (
	"reflect"
	"syscall"
	"unsafe"
)

type array struct {
	sz  uintptr
	cap int
}

func (o array) Cap() int {
	return o.cap
}

func (o array) allocSlice(len int) reflect.SliceHeader {
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

	return reflect.SliceHeader{
		Data: data,
		Len:  0,
		Cap:  len,
	}
}

func (o array) deallocSlice(p unsafe.Pointer) {
	hdr := (*reflect.SliceHeader)(p)

	_, _, errno := syscall.Syscall(
		syscall.SYS_MUNMAP,
		hdr.Data,
		o.sz*uintptr(o.cap),
		0,
	)

	if errno != 0 {
		panic(errno)
	}
}
