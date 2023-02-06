// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/asset.h>
// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/client.h>
// #include <pex/ae/sdk/search.h>
// #include <stdlib.h>
import "C"
import "unsafe"

// Holds all data necessary to perform a pex search. A search can only be
// performed using a fingerprint, but additional parameters may be supported in
// the future.
type PexSearchRequest struct {
	// A fingerprint obtained by calling either NewFingerprintFromFile
	// or NewFingerprintFromBuffer. This field is required.
	Fingerprint *Fingerprint
}

// This object is returned from PexSearchFuture.Get upon successful
// completion.
type PexSearchResult struct {
	// An ID that uniquely identifies a particular search. Can be used for diagnostics.
	LookupID string

	// The assets which the query matched against.
	Matches []*PexSearchMatch
}

// PexSearchMatch contains detailed information about the match,
// including information about the matched asset, and the matching
// segments.
type PexSearchMatch struct {
	// The asset whose fingerprint matches the query.
	Asset *Asset

	// The matching time segments on the query and asset respectively.
	Segments []*Segment
}

// PexSearchFuture object is returned by the PexSearchClient.StartSearch
// function and is used to retrieve a search result.
type PexSearchFuture struct {
	client *PexSearchClient

	LookupID string
}

// Get blocks until the search result is ready and then returns it. It
// also releases all the allocated resources, so it will return an
// error when called multiple times.
func (x *PexSearchFuture) Get() (*PexSearchResult, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cRequest := C.AE_CheckSearchRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.AE_CheckSearchRequest_Delete(&cRequest)

	cResult := C.AE_CheckSearchResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_CheckSearchResult_Delete(&cResult)

	cLookupID := C.CString(x.LookupID)
	defer C.free(unsafe.Pointer(cLookupID))

	C.AE_CheckSearchRequest_SetLookupID(cRequest, cLookupID, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	C.AE_CheckSearch(x.client.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}
	return x.processResult(cResult, cStatus)
}

func (x *PexSearchFuture) processResult(cResult *C.AE_CheckSearchResult, cStatus *C.AE_Status) (*PexSearchResult, error) {
	cMatch := C.AE_SearchMatch_New()
	if cMatch == nil {
		panic("out of memory")
	}
	defer C.AE_SearchMatch_Delete(&cMatch)

	cAsset := C.AE_Asset_New()
	if cAsset == nil {
		panic("out of memory")
	}
	defer C.AE_Asset_Delete(&cAsset)

	var cMatchesPos C.int = 0
	var matches []*PexSearchMatch

	for C.AE_CheckSearchResult_NextMatch(cResult, cMatch, &cMatchesPos) {
		var cQueryStart C.int64_t
		var cQueryEnd C.int64_t
		var cAssetStart C.int64_t
		var cAssetEnd C.int64_t
		var cType C.int
		var cSegmentsPos C.int = 0
		var segments []*Segment

		for C.AE_SearchMatch_NextSegment(cMatch, &cQueryStart, &cQueryEnd, &cAssetStart, &cAssetEnd, &cType, &cSegmentsPos) {
			segments = append(segments, &Segment{
				Type:       SegmentType(cType),
				QueryStart: int64(cQueryStart),
				QueryEnd:   int64(cQueryEnd),
				AssetStart: int64(cAssetStart),
				AssetEnd:   int64(cAssetEnd),
			})
		}

		C.AE_SearchMatch_GetAsset(cMatch, cAsset, cStatus)
		if err := statusToError(cStatus); err != nil {
			return nil, err
		}

		matches = append(matches, &PexSearchMatch{
			Asset:    newAssetFromC(cAsset),
			Segments: segments,
		})
	}

	return &PexSearchResult{
		LookupID: x.LookupID,
		Matches:  matches,
	}, nil
}

// Client serves as an entry point to all operations that
// communicate with the Attribution Engine backend service. It
// automatically handles the connection and authentication with the
// service.
type PexSearchClient struct {
	c *C.AE_Client
}

func NewPexSearchClient(clientID, clientSecret string) (*PexSearchClient, error) {
	cClient, err := newClient(C.AE_PEX_SEARCH, clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	return &PexSearchClient{
		c: cClient,
	}, nil
}

// StartSearch starts a Pex search. This operation does not block until
// the search is finished, it does however perform a network operation
// to initiate the search on the backend service.
func (x *PexSearchClient) StartSearch(req *PexSearchRequest) (*PexSearchFuture, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cRequest := C.AE_StartSearchRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.AE_StartSearchRequest_Delete(&cRequest)

	cResult := C.AE_StartSearchResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_StartSearchResult_Delete(&cResult)

	cBuffer := C.AE_Buffer_New()
	if cBuffer == nil {
		panic("out of memory")
	}
	defer C.AE_Buffer_Delete(&cBuffer)

	ftData := unsafe.Pointer(&req.Fingerprint.b[0])
	ftSize := C.size_t(len(req.Fingerprint.b))

	C.AE_Buffer_Set(cBuffer, ftData, ftSize)

	C.AE_StartSearchRequest_SetFingerprint(cRequest, cBuffer, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	C.AE_StartSearch(x.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	return &PexSearchFuture{
		client:   x,
		LookupID: C.GoString(C.AE_StartSearchResult_GetLookupID(cResult)),
	}, nil
}
