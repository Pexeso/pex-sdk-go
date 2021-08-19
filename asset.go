// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/asset.h>
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
	AssetTypeInvalid = AssetType(0)

	// A copywritten audio-visual work.
	AssetTypeVideo = AssetType(1)

	// An audio recording. Searched content may match a recording asset via 1-1
	// audio matches, or by matching it's melody (e.g a cover song).
	AssetTypeAUdioRecording = AssetType(2)

	// The composition representing the underlying lyrics and melody of a song.
	// Composition Assets are linked to associated Recording Assets.
	AssetTypeAudioComposition = AssetType(3)

	AssetTypeImage = AssetType(4)

	AssetTypeText = AssetType(5)
)

// Asset contains all information about a particular asset. Searches performed
// through the SDK match against assets representing copyrighted works.
type Asset struct {
	// The ID of the asset.
	ID uint64

	// The type of the asset.
	Type AssetType

	// The title of the asset.
	Title string

	// The artists who contributed to the asset.
	Artists []string
}

func newAssetFromC(cAsset *C.AE_Asset) *Asset {
	var cArtist *C.char
	var cArtistsPos C.int = 0
	var artists []string

	for C.AE_Asset_NextArtist(cAsset, &cArtist, &cArtistsPos) {
		artists = append(artists, C.GoString(cArtist))
	}

	return &Asset{
		ID:      uint64(C.AE_Asset_GetID(cAsset)),
		Type:    AssetType(C.AE_Asset_GetType(cAsset)),
		Title:   C.GoString(C.AE_Asset_GetTitle(cAsset)),
		Artists: artists,
	}
}
