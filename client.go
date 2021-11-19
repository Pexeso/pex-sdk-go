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
	cClientID := C.CString(clientID)
	defer C.free(unsafe.Pointer(cClientID))

	cClientSecret := C.CString(clientSecret)
	defer C.free(unsafe.Pointer(cClientSecret))

	var cStatusCode C.int
	var cStatusMessage *C.char
	defer C.free(unsafe.Pointer(cStatusMessage))

	C.AE_Init(cClientID, cClientSecret, &cStatusCode, &cStatusMessage)
	if StatusCode(cStatusCode) != StatusOK {
		return nil, &Error{
			Code:    StatusCode(cStatusCode),
			Message: C.GoString(cStatusMessage),
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

	C.AE_Client_Init(cClient, cClientID, cClientSecret, cStatus)
	if err := statusToError(cStatus); err != nil {
		C.free(unsafe.Pointer(cClient))
		return nil, err
	}

	return &Client{
		c: cClient,
	}, nil
}

// Close closes all connections to the backend service and releases
// the memory manually allocated by the core library. The
// LicenseSearch and MetadataSearch fields must not be
// used after Close is called.
func (x *Client) Close() error {
	C.AE_Client_Delete(&x.c)
	C.AE_Cleanup()
	return nil
}
