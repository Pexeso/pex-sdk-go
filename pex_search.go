// Copyright 2020 Pexeso Inc. All rights reserved.

package pex

// #include <pex/sdk/asset.h>
// #include <pex/sdk/lock.h>
// #include <pex/sdk/client.h>
// #include <pex/sdk/search.h>
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
	// IDs that uniquely identify a particular search. Can be used for diagnostics.
	LookupIDs []string

	// The assets which the query matched against.
	Matches []*PexSearchMatch
}

type PexSearchAsset struct {
	// The title of the asset.
	Title string

	// The artist who contributed to the asset.
	Artist string

	// International Standard Recording Code.
	ISRC string

	// The label that owns the asset (e.g. Sony Music Entertainment).
	Label string

	// The total duration of the asset in seconds.
	Duration float32
}

func newPexSearchAssetFromC(cAsset *C.Pex_Asset) *PexSearchAsset {
	return &PexSearchAsset{
		Title:    C.GoString(C.Pex_Asset_GetTitle(cAsset)),
		Artist:   C.GoString(C.Pex_Asset_GetArtist(cAsset)),
		ISRC:     C.GoString(C.Pex_Asset_GetISRC(cAsset)),
		Label:    C.GoString(C.Pex_Asset_GetLabel(cAsset)),
		Duration: float32(C.Pex_Asset_GetDuration(cAsset)),
	}
}

// PexSearchMatch contains detailed information about the match,
// including information about the matched asset, and the matching
// segments.
type PexSearchMatch struct {
	// The asset whose fingerprint matches the query.
	Asset *PexSearchAsset

	// The matching time segments on the query and asset respectively.
	Segments []*Segment
}

// PexSearchFuture object is returned by the PexSearchClient.StartSearch
// function and is used to retrieve a search result.
type PexSearchFuture struct {
	client *PexSearchClient

	LookupIDs []string
}

// Get blocks until the search result is ready and then returns it. It
// also releases all the allocated resources, so it will return an
// error when called multiple times.
func (x *PexSearchFuture) Get() (*PexSearchResult, error) {
	C.Pex_Lock()
	defer C.Pex_Unlock()

	cStatus := C.Pex_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.Pex_Status_Delete(&cStatus)

	cRequest := C.Pex_CheckSearchRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.Pex_CheckSearchRequest_Delete(&cRequest)

	cResult := C.Pex_CheckSearchResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.Pex_CheckSearchResult_Delete(&cResult)

	for _, lookupID := range x.LookupIDs {
		cLookupID := C.CString(lookupID)
		defer C.free(unsafe.Pointer(cLookupID))
		C.Pex_CheckSearchRequest_AddLookupID(cRequest, cLookupID)
	}

	C.Pex_CheckSearch(x.client.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}
	return x.processResult(cResult, cStatus)
}

func (x *PexSearchFuture) processResult(cResult *C.Pex_CheckSearchResult, cStatus *C.Pex_Status) (*PexSearchResult, error) {
	cMatch := C.Pex_SearchMatch_New()
	if cMatch == nil {
		panic("out of memory")
	}
	defer C.Pex_SearchMatch_Delete(&cMatch)

	cAsset := C.Pex_Asset_New()
	if cAsset == nil {
		panic("out of memory")
	}
	defer C.Pex_Asset_Delete(&cAsset)

	var cMatchesPos C.int = 0
	var matches []*PexSearchMatch

	for C.Pex_CheckSearchResult_NextMatch(cResult, cMatch, &cMatchesPos) {
		var cQueryStart C.int64_t
		var cQueryEnd C.int64_t
		var cAssetStart C.int64_t
		var cAssetEnd C.int64_t
		var cType C.int
		var cSegmentsPos C.int = 0
		var segments []*Segment

		for C.Pex_SearchMatch_NextSegment(cMatch, &cQueryStart, &cQueryEnd, &cAssetStart, &cAssetEnd, &cType, &cSegmentsPos) {
			segments = append(segments, &Segment{
				Type:       SegmentType(cType),
				QueryStart: int64(cQueryStart),
				QueryEnd:   int64(cQueryEnd),
				AssetStart: int64(cAssetStart),
				AssetEnd:   int64(cAssetEnd),
			})
		}

		C.Pex_SearchMatch_GetAsset(cMatch, cAsset, cStatus)
		if err := statusToError(cStatus); err != nil {
			return nil, err
		}

		matches = append(matches, &PexSearchMatch{
			Asset:    newPexSearchAssetFromC(cAsset),
			Segments: segments,
		})
	}

	return &PexSearchResult{
		LookupIDs: x.LookupIDs,
		Matches:   matches,
	}, nil
}

// PexSearchClient serves as an entry point to all operations that
// communicate with Pex backend services. It
// automatically handles the connection and authentication with the
// service.
type PexSearchClient struct {
	fingerprinter

	c *C.Pex_Client
}

func NewPexSearchClient(clientID, clientSecret string) (*PexSearchClient, error) {
	cClient, err := newClient(C.Pex_PEX_SEARCH, clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	return &PexSearchClient{
		c: cClient,
	}, nil
}

// Close closes all connections to the backend service and releases
// the memory manually allocated by the core library.
func (x *PexSearchClient) Close() error {
	return closeClient(&x.c)
}

func (x *PexSearchClient) getCClient() *C.Pex_Client {
	return x.c
}

// StartSearch starts a Pex search. This operation does not block until
// the search is finished, it does however perform a network operation
// to initiate the search on the backend service.
func (x *PexSearchClient) StartSearch(req *PexSearchRequest) (*PexSearchFuture, error) {
	C.Pex_Lock()
	defer C.Pex_Unlock()

	cStatus := C.Pex_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.Pex_Status_Delete(&cStatus)

	cRequest := C.Pex_StartSearchRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.Pex_StartSearchRequest_Delete(&cRequest)

	cResult := C.Pex_StartSearchResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.Pex_StartSearchResult_Delete(&cResult)

	cBuffer := C.Pex_Buffer_New()
	if cBuffer == nil {
		panic("out of memory")
	}
	defer C.Pex_Buffer_Delete(&cBuffer)

	ftData := unsafe.Pointer(&req.Fingerprint.b[0])
	ftSize := C.size_t(len(req.Fingerprint.b))

	C.Pex_Buffer_Set(cBuffer, ftData, ftSize)

	C.Pex_StartSearchRequest_SetFingerprint(cRequest, cBuffer, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	C.Pex_StartSearch(x.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	var cLookupIDPos C.size_t = 0
	var lookupIDs []string
	var cLookupID *C.char

	for C.Pex_StartSearchResult_NextLookupID(cResult, &cLookupIDPos, &cLookupID) {
		lookupIDs = append(lookupIDs, C.GoString(cLookupID))
	}

	return &PexSearchFuture{
		client:    x,
		LookupIDs: lookupIDs,
	}, nil
}
