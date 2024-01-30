// Copyright 2020 Pexeso Inc. All rights reserved.

package pex

// #include <pex/sdk/lock.h>
// #include <pex/sdk/client.h>
// #include <pex/sdk/search.h>
// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// PexSearchType can optionally be specified in the PexSearchRequest and will
// allow to retrieve results that are more relevant to the given use-case.
type PexSearchType int

const (
	// IdentifyMusic is a type of PexSearch that will return results that will
	// help identify the music in the provided media file.
	IdentifyMusic = C.Pex_CheckSearchType_IdentifyMusic

	// FindMatches is a type of PexSearch that will return all assets that
	// matched against the given media file.
	FindMatches = C.Pex_CheckSearchType_FindMatches
)

// Holds all data necessary to perform a pex search. A search can only be
// performed using a fingerprint, but additional parameters may be supported in
// the future.
type PexSearchRequest struct {
	// A fingerprint obtained by calling either NewFingerprintFromFile
	// or NewFingerprintFromBuffer. This field is required.
	Fingerprint *Fingerprint

	// Type is optional and when specified will allow to retrieve results that
	// are more relevant to the given use-case.
	Type PexSearchType
}

// This object is returned from PexSearchFuture.Get upon successful
// completion.
type PexSearchResult struct {
	// IDs that uniquely identify a particular search. Can be used for diagnostics.
	LookupIDs []string `json:"lookup_ids"`

	// The assets which the query matched against.
	Matches []*PexSearchMatch `json:"matches"`

	QueryFileDurationSeconds float32 `json:"query_file_duration_seconds"`
}

type PexSearchAsset struct {
	ID string `json:"id"`

	// The title of the asset.
	Title string `json:"title"`

	// The artist who contributed to the asset.
	Artist string `json:"artist"`

	// International Standard Recording Code.
	ISRC string `json:"isrc"`

	// The label that owns the asset (e.g. Sony Music Entertainment).
	Label string `json:"label"`

	// The total duration of the asset in seconds.
	DurationSeconds float32 `json:"duration_seconds"`

	Barcode     string `json:"barcode"`
	Distributor string `json:"distributor"`
	Subtitle    string `json:"subtitle"`
	AlbumName   string `json:"album_name"`
	ReleaseDate struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"release_date"`

	DSP []*DSP `json:"dsp"`
}

// PexSearchMatch contains detailed information about the match,
// including information about the matched asset, and the matching
// segments.
type PexSearchMatch struct {
	// The asset whose fingerprint matches the query.
	Asset *PexSearchAsset `json:"asset"`

	// The matching time segments on the query and asset respectively.
	MatchDetails MatchDetails `json:"match_details"`
}

// PexSearchFuture object is returned by the PexSearchClient.StartSearch
// function and is used to retrieve a search result.
type PexSearchFuture struct {
	client *PexSearchClient

	LookupIDs []string
	Type      PexSearchType
}

// Get blocks until the search result is ready and then returns it. It
// also releases all the allocated resources, so it will return an
// error when called multiple times.
func (x *PexSearchFuture) Get() (*PexSearchResult, error) {
	return x.client.CheckSearch(x.LookupIDs, x.Type)
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
		Type:      req.Type,
	}, nil
}

func (x *PexSearchClient) CheckSearch(lookupIDs []string, searchType PexSearchType) (*PexSearchResult, error) {
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

	for _, lookupID := range lookupIDs {
		cLookupID := C.CString(lookupID)
		defer C.free(unsafe.Pointer(cLookupID))
		C.Pex_CheckSearchRequest_AddLookupID(cRequest, cLookupID)
	}

	C.Pex_CheckSearchRequest_SetType(cRequest, C.Pex_CheckSearchType(searchType))

	C.Pex_CheckSearch(x.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}
	return x.processResult(cResult, lookupIDs)
}

func (x *PexSearchClient) processResult(cResult *C.Pex_CheckSearchResult, lookupIDs []string) (*PexSearchResult, error) {
	cJSON := C.Pex_CheckSearchResult_GetJSON(cResult)
	j := C.GoString(cJSON)

	res := new(PexSearchResult)
	if err := json.Unmarshal([]byte(j), res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	res.LookupIDs = lookupIDs
	return res, nil
}
