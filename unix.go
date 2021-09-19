//go:build !windows
// +build !windows

package gohideparam

import (
	"os"
	"unsafe"
)

func Hide(position int) {
	if position > 0 && position < len(os.Args) {
		pwdstr := (*[2]uintptr)(unsafe.Pointer(&os.Args[position]))
		for i := 0; i < len(os.Args[position]); i++ {
			*(*uint8)(unsafe.Pointer((*pwdstr)[0] + uintptr(i))) = '*'
		}
	}
}
