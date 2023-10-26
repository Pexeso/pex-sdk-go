package pex

// This file is deliberately named like this so that it's first compiled and
// first executed because it performs version checks to make sure these
// bindings are compatible with the native library.

// #cgo pkg-config: pexsdk
// #cgo LDFLAGS: -Wl,-rpath,/usr/local/lib
//
// #define PEX_SDK_MAJOR_VERSION 4
// #define PEX_SDK_MINOR_VERSION 0
//
// #include <pex/sdk/version.h>
import "C"

func init() {
	compatible := C.Pex_Version_IsCompatible(C.PEX_SDK_MAJOR_VERSION, C.PEX_SDK_MINOR_VERSION)
	if !compatible {
		panic("bindings are not compatible with the native library")
	}
}
