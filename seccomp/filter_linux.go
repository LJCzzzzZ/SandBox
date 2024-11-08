package seccomp

import "syscall"

type Filter []syscall.SockFilter

func (f Filter) SockFprog() *syscall.SockFprog {
	b := []syscall.SockFilter(f)
	return &syscall.SockFprog{
		Len:    uint16(len(b)),
		Filter: &b[0],
	}
}
