// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

// Package g exposes goroutine struct g to user space.
package g

import (
	"unsafe"
)

func getg() unsafe.Pointer

// G returns current g (the goroutine struct) to user space.
func GetG() *G {
	return (*G)(getg())
}

type stack struct {
	Lo uintptr
	Hi uintptr
}

type gobuf struct {
	sp   uintptr
	pc   uintptr
	g    uintptr
	ctxt unsafe.Pointer
	ret  uint64
	lr   uintptr
	bp   uintptr // for GOEXPERIMENT=framepointer
}

type G struct {
	Stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	_panic       uintptr // innermost panic - offset known to liblink
	_defer       uintptr // innermost defer
	m            uintptr      // current m; offset known to arm liblink
	sched        gobuf
	syscallsp    uintptr        // if status==Gsyscall, syscallsp = sched.sp to use during gc  用于gc。g处于系统调用状态时，这个值被设置成进入系统调用前的sp
	syscallpc    uintptr        // if status==Gsyscall, syscallpc = sched.pc to use during gc  同上
	stktopsp     uintptr        // expected sp at top of stack, to check in traceback
	param        unsafe.Pointer // passed parameter on wakeup
	atomicstatus uint32
	stackLock    uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	Id         int64
	schedlink    uintptr
	waitsince    int64      // approx time when the g become blocked
}

