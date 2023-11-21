package gohideparam

import (
	"os"
	"strconv"
	"unsafe"
)

// slice is the runtime representation of a slice.
// It cannot be used safely or portably and its representation may
// change in a later release.
//
// Unlike reflect.SliceHeader, its Data field is sufficient to guarantee the
// data it references will not be garbage collected.
type slice struct {
	data unsafe.Pointer
	len  uintptr
	cap  uintptr
}

// stringToBytes 没有内存开销的转换
func stringToBytes(s string) (b []byte) {
	bh := (*slice)(unsafe.Pointer(&b))
	sh := (*slice)(unsafe.Pointer(&s))
	bh.data = sh.data
	bh.len = sh.len
	bh.cap = sh.len
	return b
}

func replaceStringPointerLength(s *string, length uintptr) {
	sh := (*slice)(unsafe.Pointer(s))
	sh.len = length
}

func uint16Slice(ptr unsafe.Pointer, n uintptr) (s []uint16) {
	hdr := (*slice)(unsafe.Pointer(&s))
	hdr.data = ptr
	hdr.cap = n
	hdr.len = n
	return
}

func hideOSArg(position int) {
	if position < 0 || position >= len(os.Args) {
		panic("invalid gohideparam position" + strconv.Itoa(position))
	}
	if len(os.Args[position]) == 0 {
		return
	}
	argp := stringToBytes(os.Args[position])
	for i := 0; i < len(os.Args[position]); i++ {
		argp[i] = '*'
	}
	if len(os.Args[position]) <= 3 {
		return
	}
	argp[3] = 0
	replaceStringPointerLength(&os.Args[position], 3)
}
