//go:build !windows
// +build !windows

package gohideparam

import (
	"os"
	"unsafe"
)

func Hide(position uint) {
	pwdstr := (*[2]uintptr)(unsafe.Pointer(&os.Args[position]))
	for i := 0; i < len(os.Args[3]); i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(pwdstr)) + uintptr(i))) = '*'
	}
}
