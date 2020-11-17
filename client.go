// Copyright 2020 Pexeso Inc. All rights reserved.

package ae

// #cgo pkg-config: ae
// #include <pex/aesdk/c/client.h>
// #include <pex/aesdk/c/metadata_search.h>
// #include <pex/aesdk/c/asset_library.h>
// #include <stdlib.h>
import "C"
import "unsafe"

type Client struct {
	client *C.AE_Client
}

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

	return &Client{
		client: cClient,
	}, nil
}

func (x *Client) Close() error {
	C.AE_Client_Delete(&x.client)
	return nil
}

func (x *Client) MetadataSearch() *MetadataSearch {
	search := C.AE_MetadataSearch_New(x.client)
	if search == nil {
		panic("out of memory")
	}
	return &MetadataSearch{
		search: search,
	}
}

func (x *Client) AssetLibrary() *AssetLibrary {
	library := C.AE_AssetLibrary_New(x.client)
	if library == nil {
		panic("out of memory")
	}
	return &AssetLibrary{
		library: library,
	}
}
