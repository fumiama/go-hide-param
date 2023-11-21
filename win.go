//go:build windows
// +build windows

package gohideparam

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// next splits command line string cmd into next
// argument and command line remainder.
func (cmd *commandSlice) next(is2erase bool) {
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

func (cmd *commandSlice) erase(pos uint) {
	var p uint
	for len(*cmd) > 0 && p <= pos {
		if (*cmd)[0] == uint16(' ') || (*cmd)[0] == uint16('\t') {
			(*cmd) = (*cmd)[1:]
			continue
		}
		cmd.next(p == pos)
		p++
	}
}

type commandSlice []uint16

func utf16PtrToCommandSlice(p *uint16) commandSlice {
	if p == nil {
		return nil
	}
	// Find NUL terminator.
	end := unsafe.Pointer(p)
	start := end
	n := 0
	for *(*uint16)(end) != 0 {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
		n++
	}
	return (commandSlice)(uint16Slice(start, n))
}

// Hide replace arg at position with three `*`
//
// or less than three if len(os.Args[position]) < 3
func Hide(position int) {
	if position < 0 || position >= len(os.Args) {
		panic("invalid gohideparam position" + strconv.Itoa(position))
	}
	utf16PtrToCommandSlice(syscall.GetCommandLine()).erase(uint(position))
	hideOSArg(position)
}
