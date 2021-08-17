// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/asset.h>
// #include <pex/ae/sdk/license_search.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"sync"
)

// Holds all data necessary to perform a license search. A search can only be
// performed using a fingerprint, but additional parameters may be supported in
// the future.
type LicenseSearchRequest struct {
	// A fingerprint obtained by calling either NewFingerprintFromFile or
	// NewFingerprintFromBuffer. This field is required.
	Fingerprint *Fingerprint
}

type LicenseRightsholder struct {
	ID    uint64
	Title string
}

type LicensePolicy struct {
	ID           uint64
	CategoryID   uint64
	CategoryName string
}

type LicenseRightsholderPolicy struct {
	Rightsholder *LicenseRightsholder
	Policy       *LicensePolicy
}

type LicenseSearchMatch struct {
	// An asset that matched with the query fingerprint.
	Asset *Asset

	// A list of matching segments.
	Segments []*Segment

	// A map where the key is a territory and the value is
	// RightsholderPolicy. The territory codes conform to the ISO 3166-1
	// alpha-2 standard. For more information visit
	// https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2.
	Policies map[string][]*LicenseRightsholderPolicy
}

// This object is returned from LicenseSearchFuture.Get upon successful
// completion.
type LicenseSearchResult struct {
	// An ID that uniquely identifies a particular search. Can be used for
	// diagnostics.
	LookupID uint64

	// An ID that uniquely identifies the UGC. It is used to provide UGC metadata back to Pex.
	UGCID uint64

	// A list of matches.
	Matches []*LicenseSearchMatch
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
		LookupID: uint64(C.AE_LicenseSearchFuture_GetLookupID(cFuture)),
		c:        cFuture,
	}, nil
}

// LicenseSearchFuture is returned by the LicenseSearch.Start method
// and is used to retrieve a search result.
type LicenseSearchFuture struct {
	LookupID uint64

	c *C.AE_LicenseSearchFuture
	m sync.Mutex
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
	cMatch := C.AE_LicenseSearchMatch_New()
	if cMatch == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseSearchMatch_Delete(&cMatch)

	cAsset := C.AE_Asset_New()
	if cAsset == nil {
		panic("out of memory")
	}
	defer C.AE_Asset_Delete(&cAsset)

	cRightsholderPolicies := C.AE_LicenseRightsholderPolicies_New()
	if cRightsholderPolicies == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseRightsholderPolicies_Delete(&cRightsholderPolicies)

	cRightsholder := C.AE_LicenseRightsholder_New()
	if cRightsholder == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseRightsholder_Delete(&cRightsholder)

	cPolicy := C.AE_LicensePolicy_New()
	if cPolicy == nil {
		panic("out of memory")
	}
	defer C.AE_LicensePolicy_Delete(&cPolicy)

	var cMatchesPos C.int = 0
	var matches []*LicenseSearchMatch

	for C.AE_LicenseSearchResult_NextMatch(cResult, cMatch, &cMatchesPos) {
		// Process segments.
		var cQueryStart C.int64_t
		var cQueryEnd C.int64_t
		var cAssetStart C.int64_t
		var cAssetEnd C.int64_t
		var cSegmentsPos C.int = 0
		var segments []*Segment

		for C.AE_LicenseSearchMatch_NextSegment(cMatch, &cQueryStart, &cQueryEnd, &cAssetStart, &cAssetEnd, &cSegmentsPos) {
			segments = append(segments, &Segment{
				QueryStart: int64(cQueryStart),
				QueryEnd:   int64(cQueryEnd),
				AssetStart: int64(cAssetStart),
				AssetEnd:   int64(cAssetEnd),
			})
		}

		// Process rightsholder policies.
		var cTerritory *C.char
		var cTerritoryPoliciesPos C.int = 0
		territoryPolicies := map[string][]*LicenseRightsholderPolicy{}

		for C.AE_LicenseSearchMatch_NextTerritoryPolicies(cMatch, &cTerritory, cRightsholderPolicies, &cTerritoryPoliciesPos) {
			var cRightsholderPoliciesPos C.int = 0
			policies := make([]*LicenseRightsholderPolicy, 0)
			for C.AE_LicenseRightsholderPolicies_Next(cRightsholderPolicies, cRightsholder, cPolicy, &cRightsholderPoliciesPos) {
				policies = append(policies, &LicenseRightsholderPolicy{
					Rightsholder: &LicenseRightsholder{
						ID:    uint64(C.AE_LicenseRightsholder_GetID(cRightsholder)),
						Title: C.GoString(C.AE_LicenseRightsholder_GetTitle(cRightsholder)),
					},
					Policy: &LicensePolicy{
						ID:           uint64(C.AE_LicensePolicy_GetID(cPolicy)),
						CategoryID:   uint64(C.AE_LicensePolicy_GetCategoryID(cPolicy)),
						CategoryName: C.GoString(C.AE_LicensePolicy_GetCategoryName(cPolicy)),
					},
				})
			}
			territoryPolicies[C.GoString(cTerritory)] = policies
		}

		C.AE_LicenseSearchMatch_GetAsset(cMatch, cAsset)

		matches = append(matches, &LicenseSearchMatch{
			Asset:    newAssetFromC(cAsset),
			Segments: segments,
			Policies: territoryPolicies,
		})
	}

	return &LicenseSearchResult{
		LookupID: uint64(C.AE_LicenseSearchResult_GetLookupID(cResult)),
		UGCID:    uint64(C.AE_LicenseSearchResult_GetUGCID(cResult)),
		Matches:  matches,
	}
}
