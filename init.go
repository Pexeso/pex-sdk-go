// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/init.h>
import "C"

func init() {
	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	C.AE_Init(cStatus)
	if err := statusToError(cStatus); err != nil {
		panic(err)
	}
}
