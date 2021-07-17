package neofs

import "syscall"

func copyBuffer(dst []byte, src []byte) (uint32, syscall.Errno) {
	n := copy(dst, src)
	if n != len(src) {
		return uint32(len(src)), syscall.ERANGE
	}
	return uint32(n), 0
}
