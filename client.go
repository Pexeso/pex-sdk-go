// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #cgo pkg-config: pexae
// #include <pex/ae/sdk/init.h>
// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/client.h>
// #include <pex/ae/sdk/license_search.h>
// #include <pex/ae/sdk/metadata_search.h>
// #include <stdlib.h>
import "C"
import "unsafe"

// Client serves as an entry point to all operations that
// communicate with the Attribution Engine backend service. It
// automatically handles the connection and authentication with the
// service.
type Client struct {
	c *C.AE_Client
}

// NewClient initializes connections and authenticates with the
// backend service with the credentials provided as arguments.
func NewClient(clientID, clientSecret string) (*Client, error) {
	cClient, err := newClient(C.AE_LICENSE_SEARCH, clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	return &Client{
		c: cClient,
	}, nil
}

func newClient(typ C.AE_ClientType, clientID, clientSecret string) (*C.AE_Client, error) {
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

	cClient := C.AE_Client_New()
	if cClient == nil {
		panic("out of memory")
	}

	C.AE_Client_InitType(cClient, typ, cClientID, cClientSecret, cStatus)
	if err := statusToError(cStatus); err != nil {
		// TODO: if this fails, run AE_Cleanup
		C.free(unsafe.Pointer(cClient))
		return nil, err
	}
	return cClient, nil
}

// Close closes all connections to the backend service and releases
// the memory manually allocated by the core library. The
// LicenseSearch and MetadataSearch fields must not be
// used after Close is called.
func (x *Client) Close() error {
	C.AE_Lock()
	C.AE_Client_Delete(&x.c)
	C.AE_Unlock()

	C.AE_Cleanup()
	return nil
}
