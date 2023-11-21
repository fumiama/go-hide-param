//go:build !windows
// +build !windows

package gohideparam

// Hide replace arg at position with three `*`
//
// or less than three if len(os.Args[position]) < 3
func Hide(position int) {
	hideOSArg(position)
}
