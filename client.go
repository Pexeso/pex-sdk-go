// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #cgo pkg-config: pexae
// #include <pex/ae/sdk/init.h>
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
	// Initialized LicenseSearch struct that's using this client's
	// resources. This should be used instead of initializing the
	// struct directly.
	LicenseSearch *LicenseSearch

	// Initialized MetadataSearch struct that's using this client's
	// resources. This should be used instead of initializing the
	// struct directly.
	MetadataSearch *MetadataSearch

	Fingerprinter *Fingerprinter

	c *C.AE_Client
}

// NewClient initializes connections and authenticates with the
// backend service with the credentials provided as arguments.
func NewClient(clientID, clientSecret string) (*Client, error) {
	cClientID := C.CString(clientID)
	defer C.free(unsafe.Pointer(cClientID))

	cClientSecret := C.CString(clientSecret)
	defer C.free(unsafe.Pointer(cClientSecret))

	var cErrMsg *C.char

	C.AE_Init(cClientID, cClientSecret, &cErrMsg)
	if cErrMsg != nil {
		errMsg := C.GoString(cErrMsg)
		C.free(unsafe.Pointer(cErrMsg))

		return nil, &Error{
			Code:    StatusNotInitialized,
			Message: errMsg,
		}
	}

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

	client := &Client{
		c: cClient,
	}
	initClient(client)
	return client, nil
}

func initClient(client *Client) {
	// LicenseSearch
	if client.LicenseSearch != nil {
		C.AE_LicenseSearch_Delete(&client.LicenseSearch.c)
	}
	cLicenseSearch := C.AE_LicenseSearch_New(client.c)
	if cLicenseSearch == nil {
		panic("out of memory")
	}
	client.LicenseSearch = &LicenseSearch{
		embedded: true,
		c:        cLicenseSearch,
	}

	// MetadataSearch
	if client.MetadataSearch != nil {
		C.AE_MetadataSearch_Delete(&client.MetadataSearch.c)
	}
	cMetadataSearch := C.AE_MetadataSearch_New(client.c)
	if cMetadataSearch == nil {
		panic("out of memory")
	}
	client.MetadataSearch = &MetadataSearch{
		embedded: true,
		c:        cMetadataSearch,
	}

	// Fingerprinter
	client.Fingerprinter = &Fingerprinter{
		embedded: true,
	}
}

// Close closes all connections to the backend service and releases
// the memory manually allocated by the core library. The
// LicenseSearch and MetadataSearch fields must not be
// used after Close is called.
func (x *Client) Close() error {
	C.AE_LicenseSearch_Delete(&x.LicenseSearch.c)
	C.AE_MetadataSearch_Delete(&x.MetadataSearch.c)
	C.AE_Client_Delete(&x.c)
	C.AE_Cleanup()
	return nil
}
