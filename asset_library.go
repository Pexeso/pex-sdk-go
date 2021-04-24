// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/c/asset_library.h>
// #include <stdlib.h>
import "C"

// AssetType is how Asset are categorized. It can be either Recording,
// Composition or Video. A single piece of content may match against
// all three asset types. E.g, a music video uploaded could match to a
// Video Asset controlled by a record label, a Recording Asset
// controlled by the same label, and a Composition Asset controlled by
//a CMO representing the song writer.
type AssetType int

const (
	// An audio recording. Searched content may match a recording asset via 1-1
	// audio matches, or by matching it's melody (e.g a cover song).
	AssetTypeRecording = AssetType(0)

	// The composition representing the underlying lyrics and melody of a song.
	// Composition Assets are linked to associated Recording Assets.
	AssetTypeComposition = AssetType(1)

	// A copywritten audio-visual work.
	AssetTypeVideo = AssetType(2)
)

// Asset contains all information about a particular asset. Searches performed
// through the SDK match against assets representing copyrighted works. It can
// be retrieved using the AssetLibrary.GetAsset function.
type Asset struct {
	Metadata *AssetMetadata
}

// AssetMetadata contains metadata associated with the asset. Usually
// it comes as a part of the Asset struct.
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

	// It is a map where the key is a territory code that conforms to
	// the ISO 3166-1 alpha-2 standard. For more information visit
	// https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2.
	Licensors map[string][]string
}

// AssetLibrary encapsulates all operations on assets. Instead of
// initializating this struct directly, Client.AssetLibrary should be
// used.
type AssetLibrary struct {
	c *C.AE_AssetLibrary
}

// GetAsset retrieves information about an asset based on an asset ID.
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

	C.AE_AssetLibrary_GetAsset(x.c, C.uint64_t(id), cAsset, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	// We only need to allocate these when we've already successfully received the asset.
	cMetadata := C.AE_AssetMetadata_New()
	if cMetadata == nil {
		panic("out of memory")
	}
	defer C.AE_AssetMetadata_Delete(&cMetadata)

	C.AE_Asset_GetMetadata(cAsset, cMetadata)

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
	var assetLicensors = make(map[string][]string)

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
