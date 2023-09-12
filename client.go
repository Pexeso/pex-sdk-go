// Copyright 2020 Pexeso Inc. All rights reserved.

package pex

// #include <pex/sdk/init.h>
// #include <pex/sdk/lock.h>
// #include <pex/sdk/client.h>
// #include <stdlib.h>
import "C"
import "unsafe"

func newClient(typ C.Pex_ClientType, clientID, clientSecret string) (*C.Pex_Client, error) {
	cClientID := C.CString(clientID)
	defer C.free(unsafe.Pointer(cClientID))

	cClientSecret := C.CString(clientSecret)
	defer C.free(unsafe.Pointer(cClientSecret))

	var cStatusCode C.int
	cStatusMessage := make([]C.char, 100)
	cStatusMessageSize := C.size_t(len(cStatusMessage))

	C.AE_Init(cClientID, cClientSecret, &cStatusCode, &cStatusMessage[0], cStatusMessageSize)
	if StatusCode(cStatusCode) != StatusOK {
		return nil, &Error{
			Code:    StatusCode(cStatusCode),
			Message: C.GoString(&cStatusMessage[0]),
		}
	}

	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cClient := C.Pex_Client_New()
	if cClient == nil {
		panic("out of memory")
	}

	C.Pex_Client_Init(cClient, typ, cClientID, cClientSecret, cStatus)
	if err := statusToError(cStatus); err != nil {
		// TODO: if this fails, run AE_Cleanup
		C.free(unsafe.Pointer(cClient))
		return nil, err
	}
	return cClient, nil
}

func closeClient(c **C.Pex_Client) error {
	C.AE_Lock()
	C.Pex_Client_Delete(c)
	C.AE_Unlock()

	C.AE_Cleanup()
	return nil
}
