// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #cgo pkg-config: pexae
// #include <pex/ae/sdk/c/client.h>
// #include <pex/ae/sdk/c/license_search.h>
// #include <pex/ae/sdk/c/metadata_search.h>
// #include <pex/ae/sdk/c/asset_library.h>
// #include <stdlib.h>
import "C"
import "unsafe"

// Client serves as an entry point to all operations that
// communicate with the Attribution Engine backend service. It
// automatically handles the connection and authentication with the
// service.
type Client struct {
	// Initialized AssetLibrary struct that's using this client's
	// resources. This should be used instead of initializing the
	// struct directly.
	AssetLibrary *AssetLibrary

	// Initialized LicenseSearch struct that's using this client's
	// resources. This should be used instead of initializing the
	// struct directly.
	LicenseSearch *LicenseSearch

	// Initialized MetadataSearch struct that's using this client's
	// resources. This should be used instead of initializing the
	// struct directly.
	MetadataSearch *MetadataSearch

	c *C.AE_Client
}

// NewClient initializes connections and authenticates with the
// backend service with the credentials provided as arguments.
func NewClient(clientID, clientSecret string) (*Client, error) {
	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cClientID := C.CString(clientID)
	defer C.free(unsafe.Pointer(cClientID))

	cClientSecret := C.CString(clientSecret)
	defer C.free(unsafe.Pointer(cClientSecret))

	cClient := C.AE_Client_New()
	if cClient == nil {
		panic("out of memory")
	}

	C.AE_Client_Init(cClient, cClientID, cClientSecret, cStatus)
	if err := statusToError(cStatus); err != nil {
		C.free(unsafe.Pointer(cClient))
		return nil, err
	}

	cAssetLibrary := C.AE_AssetLibrary_New(cClient)
	if cAssetLibrary == nil {
		panic("out of memory")
	}

	cLicenseSearch := C.AE_LicenseSearch_New(cClient)
	if cLicenseSearch == nil {
		panic("out of memory")
	}

	cMetadataSearch := C.AE_MetadataSearch_New(cClient)
	if cMetadataSearch == nil {
		panic("out of memory")
	}

	return &Client{
		c: cClient,
		AssetLibrary: &AssetLibrary{
			c: cAssetLibrary,
		},
		LicenseSearch: &LicenseSearch{
			c: cLicenseSearch,
		},
		MetadataSearch: &MetadataSearch{
			c: cMetadataSearch,
		},
	}, nil
}

// Close closes all connections to the backend service and releases
// the memory manually allocated by the core library. The
// AssetLibrary, LicenseSearch and MetadataSearch fields must not be
// used after Close is called.
func (x *Client) Close() error {
	C.AE_AssetLibrary_Delete(&x.AssetLibrary.c)
	C.AE_MetadataSearch_Delete(&x.MetadataSearch.c)
	C.AE_Client_Delete(&x.c)
	return nil
}
