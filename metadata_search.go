// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/asset.h>
// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/client.h>
// #include <pex/ae/sdk/metadata_search.h>
// #include <stdlib.h>
import "C"
import "unsafe"

// Holds all data necessary to perform a metadata search. A search can only be
// performed using a fingerprint, but additional parameters may be supported in
// the future.
type MetadataSearchRequest struct {
	// A fingerprint obtained by calling either NewFingerprintFromFile
	// or NewFingerprintFromBuffer. This field is required.
	Fingerprint *Fingerprint
}

// This object is returned from MetadataSearchFuture.Get upon successful
// completion.
type MetadataSearchResult struct {
	// An ID that uniquely identifies a particular search. Can be used for diagnostics.
	LookupID uint64

	// The assets which the query matched against.
	Matches []*MetadataSearchMatch
}

// MetadataSearchMatch contains detailed information about the match,
// including information about the matched asset, and the matching
// segments.
type MetadataSearchMatch struct {
	// The asset whose fingerprint matches the query.
	Asset *Asset

	// The matching time segments on the query and asset respectively.
	Segments []*Segment
}

// MetadataSearchFuture object is returned by the Client.StartMetadataSearch
// function and is used to retrieve a search result.
type MetadataSearchFuture struct {
	client *Client

	LookupID uint64
}

// Get blocks until the search result is ready and then returns it. It
// also releases all the allocated resources, so it will return an
// error when called multiple times.
func (x *MetadataSearchFuture) Get() (*MetadataSearchResult, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cResult := C.AE_MetadataSearchResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_MetadataSearchResult_Delete(&cResult)

	C.AE_MetadataSearch_Check(x.client.c, C.uint64_t(x.LookupID), cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}
	return x.processResult(cResult), nil
}

func (x *MetadataSearchFuture) processResult(cResult *C.AE_MetadataSearchResult) *MetadataSearchResult {
	cMatch := C.AE_MetadataSearchMatch_New()
	if cMatch == nil {
		panic("out of memory")
	}
	defer C.AE_MetadataSearchMatch_Delete(&cMatch)

	cAsset := C.AE_Asset_New()
	if cAsset == nil {
		panic("out of memory")
	}
	defer C.AE_Asset_Delete(&cAsset)

	var cMatchesPos C.int = 0
	var matches []*MetadataSearchMatch

	for C.AE_MetadataSearchResult_NextMatch(cResult, cMatch, &cMatchesPos) {
		var cQueryStart C.int64_t
		var cQueryEnd C.int64_t
		var cAssetStart C.int64_t
		var cAssetEnd C.int64_t
		var cSegmentsPos C.int = 0
		var segments []*Segment

		for C.AE_MetadataSearchMatch_NextSegment(cMatch, &cQueryStart, &cQueryEnd, &cAssetStart, &cAssetEnd, &cSegmentsPos) {
			segments = append(segments, &Segment{
				QueryStart: int64(cQueryStart),
				QueryEnd:   int64(cQueryEnd),
				AssetStart: int64(cAssetStart),
				AssetEnd:   int64(cAssetEnd),
			})
		}

		C.AE_MetadataSearchMatch_GetAsset(cMatch, cAsset)

		matches = append(matches, &MetadataSearchMatch{
			Asset:    newAssetFromC(cAsset),
			Segments: segments,
		})
	}

	return &MetadataSearchResult{
		LookupID: uint64(C.AE_MetadataSearchResult_GetLookupID(cResult)),
		Matches:  matches,
	}
}

// StartMetadataSearch starts a metadata search. This operation does not block until
// the search is finished, it does however perform a network operation
// to initiate the search on the backend service.
func (x *Client) StartMetadataSearch(req *MetadataSearchRequest) (*MetadataSearchFuture, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cRequest := C.AE_MetadataSearchRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.AE_MetadataSearchRequest_Delete(&cRequest)

	cBuffer := C.AE_Buffer_New()
	if cBuffer == nil {
		panic("out of memory")
	}
	defer C.AE_Buffer_Delete(&cBuffer)

	ftData := unsafe.Pointer(&req.Fingerprint.b[0])
	ftSize := C.size_t(len(req.Fingerprint.b))

	C.AE_Buffer_Set(cBuffer, ftData, ftSize)

	C.AE_MetadataSearchRequest_SetFingerprint(cRequest, cBuffer, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	var lookupID C.uint64_t
	C.AE_MetadataSearch_Start(x.c, cRequest, &lookupID, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	return &MetadataSearchFuture{
		client:   x,
		LookupID: uint64(lookupID),
	}, nil
}
