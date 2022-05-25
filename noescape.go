package logger

import (
	"fmt"
	"io"
	"unsafe"
)

func noescape_writer(val *io.Writer) io.Writer {
	return *(*io.Writer)(noescape(unsafe.Pointer(val)))
}

func noescape_bytearray(val *[]byte) []byte {
	return *(*[]byte)(noescape(unsafe.Pointer(val)))
}

func noescape_string(val *string) string {
	return *(*string)(noescape(unsafe.Pointer(val)))
}

func noescape_stringer(val *fmt.Stringer) fmt.Stringer {
	return *(*fmt.Stringer)(noescape(unsafe.Pointer(val)))
}

func noescape_jsonvaluewriter(val *JSONValueWriter) JSONValueWriter {
	return *(*JSONValueWriter)(noescape(unsafe.Pointer(val)))
}

func noescape_stackwriterptr(val *StackWriter) *StackWriter {
	return (*StackWriter)(noescape(unsafe.Pointer(val)))
}

func noescape_interface(val *interface{}) interface{} {
	return *(*interface{})(noescape(unsafe.Pointer(val)))
}

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
