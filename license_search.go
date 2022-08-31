// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/asset.h>
// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/client.h>
// #include <pex/ae/sdk/license_search.h>
// #include <stdlib.h>
import "C"
import "unsafe"

// Holds all data necessary to perform a license search. A search can only be
// performed using a fingerprint, but additional parameters may be supported in
// the future.
type LicenseSearchRequest struct {
	// A fingerprint obtained by calling either NewFingerprintFromFile or
	// NewFingerprintFromBuffer. This field is required.
	Fingerprint *Fingerprint
}

type RightsholderPolicy struct {
	// The ID of the rightsholder.
	RightsholderID string

	// The title of the rightsholder.
	RightsholderTitle string

	// The ID of the policy.
	PolicyID string

	// The ID of the category this policy belongs to.
	PolicyCategoryID string

	// The name of the category this policy belongs to.
	PolicyCategoryName string
}

type LicenseSearchMatch struct {
	// The asset whose fingerprint matches the query.
	Asset *Asset

	// The matching time segments on the query and asset respectively.
	Segments []*Segment

	// A map where the key is a territory and the value is
	// RightsholderPolicy. The territory codes conform to the ISO 3166-1
	// alpha-2 standard. For more information visit
	// https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2.
	Policies map[string][]*RightsholderPolicy
}

// This object is returned from LicenseSearchFuture.Get upon successful
// completion.
type LicenseSearchResult struct {
	// An ID that uniquely identifies a particular search. Can be used for diagnostics.
	LookupID string

	// An ID that uniquely identifies the UGC. It is used to provide UGC metadata back to Pex.
	UGCID string

	// The assets which the query matched against.
	Matches []*LicenseSearchMatch
}

// LicenseSearchFuture is returned by the Client.StartLicenseSearch method
// and is used to retrieve a search result.
type LicenseSearchFuture struct {
	client *Client

	UGCID    string
	LookupID string
}

// Get blocks until the search result is ready and then returns it. It
// also releases all the allocated resources, so it will return an
// error when called multiple times.
func (x *LicenseSearchFuture) Get() (*LicenseSearchResult, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cRequest := C.AE_LicenseSearchCheckRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseSearchCheckRequest_Delete(&cRequest)

	cResult := C.AE_LicenseSearchCheckResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseSearchCheckResult_Delete(&cResult)

	cLookupID := C.CString(x.LookupID)
	defer C.free(unsafe.Pointer(cLookupID))

	C.AE_LicenseSearchCheckRequest_SetLookupID(cRequest, cLookupID, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	C.AE_LicenseSearch_Check(x.client.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}
	return x.processResult(cResult), nil
}

func (x *LicenseSearchFuture) processResult(cResult *C.AE_LicenseSearchCheckResult) *LicenseSearchResult {
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

	cRightsholderPolicies := C.AE_RightsholderPolicies_New()
	if cRightsholderPolicies == nil {
		panic("out of memory")
	}
	defer C.AE_RightsholderPolicies_Delete(&cRightsholderPolicies)

	cRightsholderPolicy := C.AE_RightsholderPolicy_New()
	if cRightsholderPolicy == nil {
		panic("out of memory")
	}
	defer C.AE_RightsholderPolicy_Delete(&cRightsholderPolicy)

	var cMatchesPos C.int = 0
	var matches []*LicenseSearchMatch

	for C.AE_LicenseSearchCheckResult_NextMatch(cResult, cMatch, &cMatchesPos) {
		// Process segments.
		var cQueryStart C.int64_t
		var cQueryEnd C.int64_t
		var cAssetStart C.int64_t
		var cAssetEnd C.int64_t
		var cType C.int
		var cSegmentsPos C.int = 0
		var segments []*Segment

		for C.AE_LicenseSearchMatch_NextSegment(cMatch, &cQueryStart, &cQueryEnd, &cAssetStart, &cAssetEnd, &cType, &cSegmentsPos) {
			segments = append(segments, &Segment{
				Type:       SegmentType(cType),
				QueryStart: int64(cQueryStart),
				QueryEnd:   int64(cQueryEnd),
				AssetStart: int64(cAssetStart),
				AssetEnd:   int64(cAssetEnd),
			})
		}

		// Process rightsholder policies.
		var cTerritory *C.char
		var cTerritoryPoliciesPos C.int = 0
		territoryPolicies := map[string][]*RightsholderPolicy{}

		for C.AE_LicenseSearchMatch_NextTerritoryPolicies(cMatch, &cTerritory, cRightsholderPolicies, &cTerritoryPoliciesPos) {
			var cRightsholderPoliciesPos C.int = 0
			policies := make([]*RightsholderPolicy, 0)
			for C.AE_RightsholderPolicies_Next(cRightsholderPolicies, cRightsholderPolicy, &cRightsholderPoliciesPos) {
				policies = append(policies, &RightsholderPolicy{
					RightsholderID:     C.GoString(C.AE_RightsholderPolicy_GetRightsholderID(cRightsholderPolicy)),
					RightsholderTitle:  C.GoString(C.AE_RightsholderPolicy_GetRightsholderTitle(cRightsholderPolicy)),
					PolicyID:           C.GoString(C.AE_RightsholderPolicy_GetPolicyID(cRightsholderPolicy)),
					PolicyCategoryID:   C.GoString(C.AE_RightsholderPolicy_GetPolicyCategoryID(cRightsholderPolicy)),
					PolicyCategoryName: C.GoString(C.AE_RightsholderPolicy_GetPolicyCategoryName(cRightsholderPolicy)),
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
		LookupID: C.GoString(C.AE_LicenseSearchCheckResult_GetLookupID(cResult)),
		UGCID:    C.GoString(C.AE_LicenseSearchCheckResult_GetUGCID(cResult)),
		Matches:  matches,
	}
}

// StartLicenseSearch starts a license search. This operation does not block until the
// search is finished, it does however perform a network operation to
// initiate the search on the backend service.
func (x *Client) StartLicenseSearch(req *LicenseSearchRequest) (*LicenseSearchFuture, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cRequest := C.AE_LicenseSearchStartRequest_New()
	if cRequest == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseSearchStartRequest_Delete(&cRequest)

	cResult := C.AE_LicenseSearchStartResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_LicenseSearchStartResult_Delete(&cResult)

	cBuffer := C.AE_Buffer_New()
	if cBuffer == nil {
		panic("out of memory")
	}
	defer C.AE_Buffer_Delete(&cBuffer)

	ftData := unsafe.Pointer(&req.Fingerprint.b[0])
	ftSize := C.size_t(len(req.Fingerprint.b))

	C.AE_Buffer_Set(cBuffer, ftData, ftSize)

	C.AE_LicenseSearchStartRequest_SetFingerprint(cRequest, cBuffer, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	C.AE_LicenseSearch_Start(x.c, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	return &LicenseSearchFuture{
		client:   x,
		LookupID: C.GoString(C.AE_LicenseSearchStartResult_GetLookupID(cResult)),
		UGCID:    C.GoString(C.AE_LicenseSearchStartResult_GetUGCID(cResult)),
	}, nil
}
