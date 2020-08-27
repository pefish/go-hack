// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.
// +build !windows

package syscall

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	pageSize = 4096
)

const (
	ProtectNone  = syscall.PROT_NONE
	ProtectRead  = syscall.PROT_READ
	ProtectWrite = syscall.PROT_WRITE
	ProtectExec  = syscall.PROT_EXEC
)

var (
	goexitCode = make([]byte, pageSize*2)
)

func Mprotect(ptr unsafe.Pointer, size, prot uintptr) {
	addr := uintptr(ptr)
	aligned := addr &^ (pageSize - 1)
	_, _, errno := syscall.Syscall(syscall.SYS_MPROTECT, aligned, addr-aligned+size, prot)

	if errno != 0 {
		panic(fmt.Errorf("tls: fail to call mprotect(addr=0x%x, size=%v, prot=0x%x) with error %v", addr, size, prot, errno))
	}
}
