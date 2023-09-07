// Copyright 2020 Pexeso Inc. All rights reserved.

package pex

// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/client.h>
// #include <pex/ae/sdk/ingestion.h>
// #include <pex/ae/sdk/search.h>
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
	client *PrivateSearchClient

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

	C.AE_CheckSearchRequest_SetLookupID(cRequest, cLookupID)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	C.AE_CheckSearch(x.client.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}
	return x.processResult(cResult, cStatus)
}

func (x *PrivateSearchFuture) processResult(cResult *C.AE_CheckSearchResult, cStatus *C.AE_Status) (*PrivateSearchResult, error) {
	cMatch := C.AE_SearchMatch_New()
	if cMatch == nil {
		panic("out of memory")
	}
	defer C.AE_SearchMatch_Delete(&cMatch)

	var cMatchesPos C.int = 0
	var matches []*PrivateSearchMatch

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

		cProvidedID := C.AE_SearchMatch_GetProvidedID(cMatch, cStatus)
		if err := statusToError(cStatus); err != nil {
			return nil, err
		}

		matches = append(matches, &PrivateSearchMatch{
			ProvidedID: C.GoString(cProvidedID),
			Segments:   segments,
		})
	}

	return &PrivateSearchResult{
		LookupID: x.LookupID,
		Matches:  matches,
	}, nil
}

// PrivateSearchClient serves as an entry point to all operations that
// communicate with the Attribution Engine backend service. It
// automatically handles the connection and authentication with the
// service.
type PrivateSearchClient struct {
	fingerprinter

	c *C.AE_Client
}

func NewPrivateSearchClient(clientID, clientSecret string) (*PrivateSearchClient, error) {
	cClient, err := newClient(C.AE_PRIVATE_SEARCH, clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	return &PrivateSearchClient{
		c: cClient,
	}, nil
}

// Close closes all connections to the backend service and releases
// the memory manually allocated by the core library.
func (x *PrivateSearchClient) Close() error {
	return closeClient(&x.c)
}

func (x *PrivateSearchClient) getCClient() *C.AE_Client {
	return x.c
}

// StartSearch starts a private search. This operation does not block until
// the search is finished, it does however perform a network operation
// to initiate the search on the backend service.
func (x *PrivateSearchClient) StartSearch(req *PrivateSearchRequest) (*PrivateSearchFuture, error) {
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

	return &PrivateSearchFuture{
		client:   x,
		LookupID: C.GoString(C.AE_StartSearchResult_GetLookupID(cResult)),
	}, nil
}

// Ingest ingests a fingerprint into the private search
// catalog. The catalog is determined from the authentication credentials used
// when initializing the client. If you want to ingest into multiple catalogs
// within one application, you need to use multiple clients. The id parameter
// identifies the fingerprint and will be returned during search to identify
// the matched asset.
func (x *PrivateSearchClient) Ingest(id string, ft *Fingerprint) error {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cID := C.CString(id)
	defer C.free(unsafe.Pointer(cID))

	cBuffer := C.AE_Buffer_New()
	if cBuffer == nil {
		panic("out of memory")
	}
	defer C.AE_Buffer_Delete(&cBuffer)

	ftData := unsafe.Pointer(&ft.b[0])
	ftSize := C.size_t(len(ft.b))

	C.AE_Buffer_Set(cBuffer, ftData, ftSize)

	C.AE_Ingest(x.c, cID, cBuffer, cStatus)
	return statusToError(cStatus)
}
