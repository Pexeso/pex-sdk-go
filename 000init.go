package pexae

// This file is deliberately named like this so that it's first compiled and
// first executed because it performs version checks to make sure these
// bindings are compatible with the native library.

// #cgo pkg-config: pexae
//
// #define AE_SDK_MAJOR_VERSION 3
// #define AE_SDK_MINOR_VERSION 1
//
// #include <pex/ae/sdk/version.h>
import "C"

func init() {
	compatible := C.AE_Version_IsCompatible(C.AE_SDK_MAJOR_VERSION, C.AE_SDK_MINOR_VERSION)
	if !compatible {
		panic("bindings are not compatible with the native library")
	}
}
