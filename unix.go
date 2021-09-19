//go:build !windows
// +build !windows

package gohideparam

import (
	"os"
	"unsafe"
)

func Hide(position int) {
	if position > 0 && position < len(os.Args) {
		p := *(*unsafe.Pointer)(unsafe.Pointer(&os.Args[position]))
		for i := 0; i < len(os.Args[position]); i++ {
			*(*uint8)(unsafe.Pointer(uintptr(p) + uintptr(i))) = '*'
		}
	}
}
