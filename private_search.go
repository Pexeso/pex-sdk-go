// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/client.h>
// #include <pex/ae/sdk/private_search.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

// Holds all data necessary to perform a private search. A search can only be
// performed using a fingerprint, but additional parameters may be supported in
// the future.
type PrivateSearchRequest struct {
	// A fingerprint obtained by calling either NewFingerprintFromFile
	// or NewFingerprintFromBuffer. This field is required.
	Fingerprint *Fingerprint
}

// This object is returned from PrivateSearchFuture.Get upon successful
// completion.
type PrivateSearchResult struct {
	// An ID that uniquely identifies a particular search. Can be used for diagnostics.
	LookupID string

	// The assets which the query matched against.
	Matches []*PrivateSearchMatch
}

// PrivateSearchMatch contains detailed information about the match,
// including information about the matched asset, and the matching
// segments.
type PrivateSearchMatch struct {
	// The ID provided during ingestion.
	ProvidedID string

	// The matching time segments on the query and asset respectively.
	Segments []*Segment
}

// PrivateSearchFuture object is returned by the Client.StartPrivateSearch
// function and is used to retrieve a search result.
type PrivateSearchFuture struct {
	client *Client

	LookupID string
}

// Get blocks until the search result is ready and then returns it. It
// also releases all the allocated resources, so it will return an
// error when called multiple times.
func (x *PrivateSearchFuture) Get() (*PrivateSearchResult, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cRequest := C.AE_PrivateSearchCheckRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.AE_PrivateSearchCheckRequest_Delete(&cRequest)

	cResult := C.AE_PrivateSearchCheckResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_PrivateSearchCheckResult_Delete(&cResult)

	cLookupID := C.CString(x.LookupID)
	defer C.free(unsafe.Pointer(cLookupID))

	C.AE_PrivateSearchCheckRequest_SetLookupID(cRequest, cLookupID, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	C.AE_PrivateSearch_Check(x.client.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}
	return x.processResult(cResult), nil
}

func (x *PrivateSearchFuture) processResult(cResult *C.AE_PrivateSearchCheckResult) *PrivateSearchResult {
	cMatch := C.AE_PrivateSearchMatch_New()
	if cMatch == nil {
		panic("out of memory")
	}
	defer C.AE_PrivateSearchMatch_Delete(&cMatch)

	var cMatchesPos C.int = 0
	var matches []*PrivateSearchMatch

	for C.AE_PrivateSearchCheckResult_NextMatch(cResult, cMatch, &cMatchesPos) {
		var cQueryStart C.int64_t
		var cQueryEnd C.int64_t
		var cAssetStart C.int64_t
		var cAssetEnd C.int64_t
		var cType C.int
		var cSegmentsPos C.int = 0
		var segments []*Segment

		for C.AE_PrivateSearchMatch_NextSegment(cMatch, &cQueryStart, &cQueryEnd, &cAssetStart, &cAssetEnd, &cType, &cSegmentsPos) {
			segments = append(segments, &Segment{
				Type:       SegmentType(cType),
				QueryStart: int64(cQueryStart),
				QueryEnd:   int64(cQueryEnd),
				AssetStart: int64(cAssetStart),
				AssetEnd:   int64(cAssetEnd),
			})
		}

		matches = append(matches, &PrivateSearchMatch{
			ProvidedID: C.GoString(C.AE_PrivateSearchMatch_GetProvidedID(cMatch)),
			Segments:   segments,
		})
	}

	return &PrivateSearchResult{
		LookupID: C.GoString(C.AE_PrivateSearchCheckResult_GetLookupID(cResult)),
		Matches:  matches,
	}
}

// StartPrivateSearch starts a private search. This operation does not block until
// the search is finished, it does however perform a network operation
// to initiate the search on the backend service.
func (x *Client) StartPrivateSearch(req *PrivateSearchRequest) (*PrivateSearchFuture, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cRequest := C.AE_PrivateSearchStartRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.AE_PrivateSearchStartRequest_Delete(&cRequest)

	cResult := C.AE_PrivateSearchStartResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_PrivateSearchStartResult_Delete(&cResult)

	cBuffer := C.AE_Buffer_New()
	if cBuffer == nil {
		panic("out of memory")
	}
	defer C.AE_Buffer_Delete(&cBuffer)

	ftData := unsafe.Pointer(&req.Fingerprint.b[0])
	ftSize := C.size_t(len(req.Fingerprint.b))

	C.AE_Buffer_Set(cBuffer, ftData, ftSize)

	C.AE_PrivateSearchStartRequest_SetFingerprint(cRequest, cBuffer, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	C.AE_PrivateSearch_Start(x.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	return &PrivateSearchFuture{
		client:   x,
		LookupID: C.GoString(C.AE_PrivateSearchStartResult_GetLookupID(cResult)),
	}, nil
}
