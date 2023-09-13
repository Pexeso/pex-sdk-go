// Copyright 2020 Pexeso Inc. All rights reserved.

package pex

// #include <pex/sdk/lock.h>
// #include <pex/sdk/mockserver.h>
import "C"

type client interface {
	getCClient() *C.Pex_Client
}

// MockClient initializes the provided client to communicate with the mockserver.
func MockClient(c client) error {
	C.Pex_Lock()
	defer C.Pex_Unlock()

	cStatus := C.Pex_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.Pex_Status_Delete(&cStatus)

	C.Pex_Mockserver_InitClient(c.getCClient(), nil, cStatus)
	return statusToError(cStatus)
}
