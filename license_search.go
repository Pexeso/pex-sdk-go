// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/c/license_search.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"sync"
)

// BasicPolicy is an enumeration of possible license policies for queried
// content.
type BasicPolicy int

const (
	// The content should be allowed to be uploaded to the platform.
	BasicPolicyAllow = BasicPolicy(0)

	// The content should not be allowed to be uploaded to the platform,
	// because it includes copyrighted content.
	BasicPolicyBlock = BasicPolicy(1)
)

// Holds all data necessary to perform a license search. A search can only be
// performed using a fingerprint, but additional parameters may be supported in
// the future.
type LicenseSearchRequest struct {
	// A fingerprint obtained by calling either NewFingerprintFromFile or
	// NewFingerprintFromBuffer. This field is required.
	Fingerprint *Fingerprint
}

// This object is returned from LicenseSearchFuture.Get upon successful
// completion.
type LicenseSearchResult struct {
	// An ID that uniquely identifies a particular search. Can be used for
	// diagnostics.
	LookupID uint64

	// An ID that uniquely identifies the UGC. It is used to provide UGC metadata back to Pex.
	UGCID uint64

	// A map where the key is a territory and the value is BasicPolicy (either
	// allow or block). The territory codes conform to
	// the ISO 3166-1 alpha-2 standard. For more information visit
	// https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2.
	Policies map[string]BasicPolicy
}

// This class encapsulates all operations necessary to perform a license
// search. Instead of instantiating the class directly,
// Client.LicenseSearch should be used.
type LicenseSearch struct {
	c *C.AE_LicenseSearch
}

// Starts a license search. This operation does not block until the
// search is finished, it does however perform a network operation to
// initiate the search on the backend service.
func (x *LicenseSearch) Start(req *LicenseSearchRequest) (*LicenseSearchFuture, error) {
	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cRequest := C.AE_LicenseSearchRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseSearchRequest_Delete(&cRequest)

	cFuture := C.AE_LicenseSearchFuture_New()
	if cFuture == nil {
		panic("out of memory")
	}

	C.AE_LicenseSearchRequest_SetFingerprint(cRequest, req.Fingerprint.ft)

	C.AE_LicenseSearch_Start(x.c, cRequest, cFuture, cStatus)
	if err := statusToError(cStatus); err != nil {
		// Delete the resource here to prevent leaking.
		C.AE_LicenseSearchFuture_Delete(&cFuture)
		return nil, err
	}

	return &LicenseSearchFuture{
		c: cFuture,
	}, nil
}

// LicenseSearchFuture is returned by the LicenseSearch.Start method
// and is used to retrieve a search result.
type LicenseSearchFuture struct {
	c *C.AE_LicenseSearchFuture
	m sync.Mutex
}

func (x *LicenseSearchFuture) LookupID() uint64 {
	return C.AE_LicenseSearchFuture_GetLookupID(x.c)
}

// Get blocks until the search result is ready and then returns it. It
// also releases all the allocated resources, so it will return an
// error when called multiple times.
func (x *LicenseSearchFuture) Get() (*LicenseSearchResult, error) {
	x.m.Lock()
	defer x.m.Unlock()

	if x.c == nil {
		return nil, errors.New("already called")
	}
	defer x.close()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cResult := C.AE_LicenseSearchResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseSearchResult_Delete(&cResult)

	C.AE_LicenseSearchFuture_Get(x.c, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}
	return x.processResult(cResult), nil
}

func (x *LicenseSearchFuture) close() {
	C.AE_LicenseSearchFuture_Delete(&x.c)
	x.c = nil
}

func (x *LicenseSearchFuture) processResult(cResult *C.AE_LicenseSearchResult) *LicenseSearchResult {
	var cTerritory *C.char
	var cPolicy C.int
	var cPoliciesPos C.size_t = 0
	policies := make(map[string]BasicPolicy)

	for C.AE_LicenseSearchResult_NextPolicy(cResult, &cTerritory, &cPolicy, &cPoliciesPos) {
		policies[C.GoString(cTerritory)] = BasicPolicy(cPolicy)
	}

	return &LicenseSearchResult{
		LookupID: uint64(C.AE_LicenseSearchResult_GetLookupID(cResult)),
		UGCID:    uint64(C.AE_LicenseSearchResult_GetUGCID(cResult)),
		Policies: policies,
	}
}
