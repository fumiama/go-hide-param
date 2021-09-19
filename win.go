//go:build windows
// +build windows

package gohideparam

import (
	"os"
	"syscall"
	"unsafe"
)

// Slice is the runtime representation of a slice.
// It cannot be used safely or portably and its representation may
// change in a later release.
//
// Unlike reflect.SliceHeader, its Data field is sufficient to guarantee the
// data it references will not be garbage collected.
type Slice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// readNextArg splits command line string cmd into next
// argument and command line remainder.
func readNextArg(cmd *[]uint16, is2erase bool) {
	var inquote bool
	var nslash int
	for ; len(*cmd) > 0; (*cmd) = (*cmd)[1:] {
		switch (*cmd)[0] {
		case uint16(' '), uint16('\t'):
			if !inquote {
				return
			}
		case uint16('"'):
			if nslash%2 == 0 {
				// use "Prior to 2008" rule from
				// http://daviddeley.com/autohotkey/parameters/parameters.htm
				// section 5.2 to deal with double double quotes
				if inquote && len(*cmd) > 1 && (*cmd)[1] == uint16('"') {
					*cmd = (*cmd)[1:]
				}
				inquote = !inquote
			}
			nslash = 0
			continue
		case uint16('\\'):
			nslash++
			if is2erase {
				(*cmd)[0] = uint16('*')
			}
			continue
		default:
			if is2erase {
				(*cmd)[0] = uint16('*')
			}
		}
		nslash = 0
	}
}

func eraseCommandLine(cmd *[]uint16, pos uint) {
	var p uint
	for len(*cmd) > 0 && p <= pos {
		if (*cmd)[0] == uint16(' ') || (*cmd)[0] == uint16('\t') {
			(*cmd) = (*cmd)[1:]
			continue
		}
		readNextArg(cmd, p == pos)
		p++
	}
}

func utf16PtrToSlice(p *uint16) []uint16 {
	if p == nil {
		return nil
	}
	// Find NUL terminator.
	end := unsafe.Pointer(p)
	n := 0
	for *(*uint16)(end) != 0 {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
		n++
	}
	// Turn *uint16 into []uint16.
	var s []uint16
	hdr := (*Slice)(unsafe.Pointer(&s))
	hdr.Data = unsafe.Pointer(p)
	hdr.Cap = n
	hdr.Len = n
	return s
}

func Hide(position int) {
	if position > 0 && position < len(os.Args) {
		cmd := utf16PtrToSlice(syscall.GetCommandLine())
		eraseCommandLine(&cmd, uint(position))
	}
}
