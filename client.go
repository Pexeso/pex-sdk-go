// Copyright 2020 Pexeso Inc. All rights reserved.

package ae

// #cgo pkg-config: pexae
// #include <pex/ae/client.h>
// #include <pex/ae/metadata_search.h>
// #include <pex/ae/asset_library.h>
// #include <stdlib.h>
import "C"
import "unsafe"

type Client struct {
	client *C.AE_Client
}

func NewClient(clientID, clientSecret string) (*Client, error) {
	cClientID := C.CString(clientID)
	defer C.free(unsafe.Pointer(cClientID))

	cClientSecret := C.CString(clientSecret)
	defer C.free(unsafe.Pointer(cClientSecret))

	client := C.AE_Client_New(cClientID, cClientSecret)
	if client == nil {
		panic("out of memory")
	}
	return &Client{
		client: client,
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
