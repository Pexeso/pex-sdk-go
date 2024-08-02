// Copyright 2020 Pexeso Inc. All rights reserved.

package pex

// #include <pex/sdk/lock.h>
// #include <pex/sdk/client.h>
// #include <pex/sdk/ingestion.h>
// #include <pex/sdk/search.h>
// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"fmt"
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
	// IDs that uniquely identify a particular search. Can be used for diagnostics.
	LookupIDs []string `json:"lookup_ids"`

	// The assets which the query matched against.
	Matches []*PrivateSearchMatch `json:"matches"`

	QueryFileDurationSeconds float32 `json:"query_file_duration_seconds"`
}

// PrivateSearchMatch contains detailed information about the match,
// including information about the matched asset, and the matching
// segments.
type PrivateSearchMatch struct {
	// The ID provided during ingestion.
	ProvidedID string `json:"provided_id"`

	// The matching time segments on the query and asset respectively.
	MatchDetails *MatchDetails `json:"match_details"`
}

// PrivateSearchFuture object is returned by the Client.StartPrivateSearch
// function and is used to retrieve a search result.
type PrivateSearchFuture struct {
	client *PrivateSearchClient

	LookupIDs []string
}

// Get blocks until the search result is ready and then returns it. It
// also releases all the allocated resources, so it will return an
// error when called multiple times.
func (x *PrivateSearchFuture) Get() (*PrivateSearchResult, error) {
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

func (x *PrivateSearchFuture) processResult(cResult *C.Pex_CheckSearchResult, cStatus *C.Pex_Status) (*PrivateSearchResult, error) {
	cJSON := C.Pex_CheckSearchResult_GetJSON(cResult)
	j := C.GoString(cJSON)

	res := new(PrivateSearchResult)
	if err := json.Unmarshal([]byte(j), res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	res.LookupIDs = x.LookupIDs
	return res, nil
}

// PrivateSearchClient serves as an entry point to all operations that
// communicate with Pex backend services. It
// automatically handles the connection and authentication with the
// service.
type PrivateSearchClient struct {
	fingerprinter

	c *C.Pex_Client
}

func NewPrivateSearchClient(clientID, clientSecret string) (*PrivateSearchClient, error) {
	cClient, err := newClient(C.Pex_PRIVATE_SEARCH, clientID, clientSecret)
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

func (x *PrivateSearchClient) getCClient() *C.Pex_Client {
	return x.c
}

// StartSearch starts a private search. This operation does not block until
// the search is finished, it does however perform a network operation
// to initiate the search on the backend service.
func (x *PrivateSearchClient) StartSearch(req *PrivateSearchRequest) (*PrivateSearchFuture, error) {
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

	return &PrivateSearchFuture{
		client:    x,
		LookupIDs: lookupIDs,
	}, nil
}

// Ingest ingests a fingerprint into the private search
// catalog. The catalog is determined from the authentication credentials used
// when initializing the client. If you want to ingest into multiple catalogs
// within one application, you need to use multiple clients. The id parameter
// identifies the fingerprint and will be returned during search to identify
// the matched asset.
func (x *PrivateSearchClient) Ingest(id string, ft *Fingerprint) error {
	C.Pex_Lock()
	defer C.Pex_Unlock()

	cStatus := C.Pex_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.Pex_Status_Delete(&cStatus)

	cID := C.CString(id)
	defer C.free(unsafe.Pointer(cID))

	cBuffer := C.Pex_Buffer_New()
	if cBuffer == nil {
		panic("out of memory")
	}
	defer C.Pex_Buffer_Delete(&cBuffer)

	ftData := unsafe.Pointer(&ft.b[0])
	ftSize := C.size_t(len(ft.b))

	C.Pex_Buffer_Set(cBuffer, ftData, ftSize)

	C.Pex_Ingest(x.c, cID, cBuffer, cStatus)
	return statusToError(cStatus)
}

// Archive archives a fingerprint that was previously ingested into the private search
// catalog. The catalog is determined from the authentication credentials used
// when initializing the client.
func (x *PrivateSearchClient) Archive(id string, types ...FingerprintType) error {
	C.Pex_Lock()
	defer C.Pex_Unlock()

	cStatus := C.Pex_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.Pex_Status_Delete(&cStatus)

	cID := C.CString(id)
	defer C.free(unsafe.Pointer(cID))

	C.Pex_Archive(x.c, cID, C.int(reduceTypes(types)), cStatus)
	return statusToError(cStatus)
}

type Entry struct {
	ProvidedID       string            `json:"provided_id"`
	FingerprintTypes []FingerprintType `json:"fingerprint_types"`
}

type ListEntriesRequest struct {
	Limit int
	After string
}

type ListEntriesResult struct {
	Entries     []Entry `json:"entries"`
	EndCursor   string  `json:"end_cursor"`
	HasNextPage bool    `json:"has_next_page"`
}

type Lister struct {
	c *C.Pex_Client

	Limit       int
	EndCursor   string
	HasNextPage bool
}

func (x *Lister) List() ([]Entry, error) {
	C.Pex_Lock()
	defer C.Pex_Unlock()

	cReq := C.Pex_ListRequest_New()
	if cReq == nil {
		panic("out of memory")
	}
	defer C.Pex_ListRequest_Delete(&cReq)

	cRes := C.Pex_ListResult_New()
	if cRes == nil {
		panic("out of memory")
	}
	defer C.Pex_ListResult_Delete(&cRes)

	cStatus := C.Pex_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.Pex_Status_Delete(&cStatus)

	cAfter := C.CString(x.EndCursor)
	defer C.free(unsafe.Pointer(cAfter))

	C.Pex_ListRequest_SetLimit(cReq, C.int(x.Limit))
	C.Pex_ListRequest_SetAfter(cReq, cAfter)

	C.Pex_List(x.c, cReq, cRes, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	j := C.GoString(C.Pex_ListResult_GetJSON(cRes))

	var res ListEntriesResult
	if err := json.Unmarshal([]byte(j), &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	x.EndCursor = res.EndCursor
	x.HasNextPage = res.HasNextPage

	return res.Entries, nil
}

func (x *PrivateSearchClient) ListEntries(req *ListEntriesRequest) *Lister {
	return &Lister{
		c:           x.c,
		Limit:       req.Limit,
		EndCursor:   req.After,
		HasNextPage: true,
	}
}
