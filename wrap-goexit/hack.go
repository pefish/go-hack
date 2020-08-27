package wrap_goexit

import (
	get_gid "github.com/pefish/go-hack/get-gid"
	"github.com/pefish/go-hack/pkg/g"
	"github.com/pefish/go-hack/pkg/syscall"
	"runtime"
	"unsafe"
)

var fns = make(map[int64]func(), 0)

func hackedGoexit()

func hackedGoexit1() {
	gid := get_gid.GetGid()
	fn, ok := fns[gid]
	if ok {
		fn()
	}
	runtime.Goexit()
}

type eface struct {  // 空interface类型实际上是这个struct
	_type unsafe.Pointer
	data  unsafe.Pointer  // 指向源struct数据的指针
}

type _func struct {
	entry   uintptr // start pc
	nameoff int32   // function name

	args int32 // in/out args size
	_    int32 // previously legacy frame size; kept for layout compatibility

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	nfuncdata int32
}

func funcPC(f interface{}) uintptr {
	return *(*uintptr)((*eface)(unsafe.Pointer(&f)).data)
}

var (
	originalGoexitFnPC uintptr
	hackedGoexitFnPC  uintptr
)

func init() {
	// 获取原始的goexit函数的入口
	ch := make(chan uintptr, 1)
	go func() {
		pc := make([]uintptr, 16)
		sz := runtime.Callers(0, pc)
		ch <- pc[sz-1]
	}()
	originalGoexitFnPC = <-ch   // 栈中保存的函数的跳转pc要比函数入口大PCQuantum(amd64下是1)个字节

	// 将函数的pcsp设置为0
	hackedGoexitFnPC = funcPC(hackedGoexit)
	fnHacked := runtime.FuncForPC(hackedGoexitFnPC)
	fnSymtab := (*_func)(unsafe.Pointer(fnHacked))
	funcSymbolSize := unsafe.Sizeof(_func{})
	syscall.Mprotect(unsafe.Pointer(fnSymtab), funcSymbolSize, syscall.ProtectWrite)  // 要修改hackedGoexit函数的pcsp信息，需要去掉写保护
	fnSymtab.pcsp = 0
	syscall.Mprotect(unsafe.Pointer(fnSymtab), funcSymbolSize, syscall.ProtectRead)
}

func hackGoexit(from, to uintptr) (success bool) {
	s := g.GetG().Stack
	for offset := s.Lo; offset <= s.Hi; offset += 4 {  // 替换栈中所有指向from函数的指针值
		val := (*uintptr)(unsafe.Pointer(offset))

		if *val == from {
			*val = to
			success = true
		}
	}

	return
}

func WrapGoexit(f func()) bool {
	gid := get_gid.GetGid()
	fns[gid] = f
	return hackGoexit(originalGoexitFnPC, hackedGoexitFnPC)
}

