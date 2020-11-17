// Copyright 2020 Pexeso Inc. All rights reserved.

package ae

// #include <pex/ae/asset_library.h>
// #include <stdlib.h>
import "C"

type AssetType int

const (
	AssetTypeRecording   = AssetType(0)
	AssetTypeComposition = AssetType(1)
	AssetTypeVideo       = AssetType(2)
	AssetTypeImage       = AssetType(3)
	AssetTypeText        = AssetType(4)
)

type Asset struct {
	Metadata *AssetMetadata
}

// Metadata associated with the asset.
type AssetMetadata struct {
	// An international standard code for uniquely identifying sound recordings
	// and music video recordings.
	ISRC string

	// The name of the track recording for a given ISRC.
	Title string

	// The names of the recording artists for a given ISRC.
	Artists []string

	// The unique codes associated with the sale of a recording.
	UPCs []string

	// The entities that own the rights to the given UPC and are entitled to
	// license its use and collect royalties.
	Licensors map[string][]string
}

type AssetLibrary struct {
	library *C.AE_AssetLibrary
}

func (x *AssetLibrary) GetAsset(id uint64) (*Asset, error) {
	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cAsset := C.AE_Asset_New()
	if cAsset == nil {
		panic("out of memory")
	}
	defer C.AE_Asset_Delete(&cAsset)

	C.AE_AssetLibrary_GetAsset(x.library, C.uint64_t(id), cAsset, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	// We only need to allocate these when we've already successfully received the asset.
	cMetadata := C.AE_AssetMetadata_New()
	if cMetadata == nil {
		panic("out of memory")
	}
	defer C.AE_AssetMetadata_Delete(&cMetadata)

	// Extract Artists
	var cArtist *C.char
	var cArtistsPos C.size_t = 0
	var artists []string

	for C.AE_AssetMetadata_NextArtist(cMetadata, &cArtist, &cArtistsPos) {
		artists = append(artists, C.GoString(cArtist))
	}

	// Extract UPCs
	var cUpc *C.char
	var cUpcsPos C.size_t = 0
	var upcs []string

	for C.AE_AssetMetadata_NextUPC(cMetadata, &cUpc, &cUpcsPos) {
		upcs = append(upcs, C.GoString(cUpc))
	}

	// Extract Licensors
	cAssetLicensors := C.AE_AssetLicensors_New()
	if cAssetLicensors == nil {
		panic("out of memory")
	}
	defer C.AE_AssetLicensors_Delete(&cAssetLicensors)

	var cAssetLicensorsPos C.size_t = 0
	var assetLicensors map[string][]string

	for C.AE_AssetMetadata_NextLicensors(cMetadata, cAssetLicensors, &cAssetLicensorsPos) {
		var cLicensor *C.char
		var cLicensorsPos C.size_t = 0
		var licensors []string

		for C.AE_AssetLicensors_NextLicensor(cAssetLicensors, &cLicensor, &cLicensorsPos) {
			licensors = append(licensors, C.GoString(cLicensor))
		}

		territory := C.GoString(C.AE_AssetLicensors_GetTerritory(cAssetLicensors))
		assetLicensors[territory] = licensors
	}

	return &Asset{
		Metadata: &AssetMetadata{
			ISRC:      C.GoString(C.AE_AssetMetadata_GetISRC(cMetadata)),
			Title:     C.GoString(C.AE_AssetMetadata_GetTitle(cMetadata)),
			Artists:   artists,
			UPCs:      upcs,
			Licensors: assetLicensors,
		},
	}, nil
}
