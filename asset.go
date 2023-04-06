// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/asset.h>
// #include <stdlib.h>
import "C"

// Asset contains all information about a particular asset. Searches performed
// through the SDK match against assets representing copyrighted works.
type Asset struct {
	// The ID of the asset.
	ID string

	// The type of the asset.
	Type string

	// The title of the asset.
	Title string

	// The artist who contributed to the asset.
	Artist string
}

func newAssetFromC(cAsset *C.AE_Asset) *Asset {
	return &Asset{
		ID:     C.GoString(C.AE_Asset_GetID(cAsset)),
		Type:   C.GoString(C.AE_Asset_GetTypeStr(cAsset)),
		Title:  C.GoString(C.AE_Asset_GetTitle(cAsset)),
		Artist: C.GoString(C.AE_Asset_GetArtist(cAsset)),
	}
}
