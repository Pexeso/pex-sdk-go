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
// a CMO representing the song writer.
type AssetType int

func (x AssetType) String() string {
	switch x {
	case 1:
		return "video"
	case 2:
		return "audio_recording"
	case 3:
		return "audio_composition"
	default:
		return "invalid"
	}
}

const (
	AssetTypeInvalid = AssetType(0)

	// A copywritten audio-visual work.
	AssetTypeVideo = AssetType(1)

	// An audio recording. Searched content may match a recording asset via 1-1
	// audio matches, or by matching it's melody (e.g a cover song).
	AssetTypeAudioRecording = AssetType(2)

	// The composition representing the underlying lyrics and melody of a song.
	// Composition Assets are linked to associated Recording Assets.
	AssetTypeAudioComposition = AssetType(3)
)

// Asset contains all information about a particular asset. Searches performed
// through the SDK match against assets representing copyrighted works.
type Asset struct {
	// The ID of the asset.
	ID string

	// The type of the asset.
	Type AssetType

	// The title of the asset.
	Title string

	// The artist who contributed to the asset.
	Artist string
}

func newAssetFromC(cAsset *C.AE_Asset) *Asset {
	return &Asset{
		ID:     C.GoString(C.AE_Asset_GetID(cAsset)),
		Type:   AssetType(C.AE_Asset_GetType(cAsset)),
		Title:  C.GoString(C.AE_Asset_GetTitle(cAsset)),
		Artist: C.GoString(C.AE_Asset_GetArtist(cAsset)),
	}
}
