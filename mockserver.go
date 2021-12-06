// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/mockserver.h>
import "C"

// MockClient initializes the provided client to communicate with the mockserver.
func MockClient(client *Client) error {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	C.AE_Mockserver_InitClient(client.c, nil, cStatus)
	return statusToError(cStatus)
}
