package ae

// #include <pex/ae/metadata_search.h>
// #include <stdlib.h>
import "C"
import "time"

// MetadataSearch implements all functions necessary to perform MetadataSearch search
// (public search - public db).
type MetadataSearch struct {
	search *C.AE_MetadataSearch
}

type MetadataSearchRequest struct {
	// A fingerprint obtained by calling either NewFingerprintFromFile or
	// NewFingerprintFromBuffer. This field is required.
	Fingerprint *Fingerprint
}

type MetadataSearchResult struct {
	// ID that uniquely identifies this lookup. It can later be used for
	// retrieving the lookup detail.
	LookupID uint64

	CompletedAt time.Time

	Matches []*MetadataSearchMatch
}

type MetadataSearchMatch struct {
	AssetID   uint64
	AssetType AssetType
	Segments  []*Segment
}

// Lookup performs the actual search.
func (x *MetadataSearch) Do(req *MetadataSearchRequest) (*MetadataSearchResult, error) {
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

	cResult := C.AE_MetadataSearchResult_New()
	if cResult == nil {
		panic("out of memory")
	}
	defer C.AE_MetadataSearchResult_Delete(&cResult)

	C.AE_MetadataSearchRequest_SetFingerprint(cRequest, req.Fingerprint.ft)

	C.AE_MetadataSearch_Do(x.search, cRequest, cResult, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	cMatch := C.AE_MetadataSearchMatch_New()
	if cMatch == nil {
		panic("out of memory")
	}
	defer C.AE_MetadataSearchMatch_Delete(&cMatch)

	var cMatchesPos C.size_t = 0
	var matches []*MetadataSearchMatch

	for C.AE_MetadataSearchResult_NextMatch(cResult, cMatch, &cMatchesPos) {
		var cQueryStart C.int64_t
		var cQueryEnd C.int64_t
		var cAssetStart C.int64_t
		var cAssetEnd C.int64_t
		var cSegmentsPos C.size_t = 0
		var segments []*Segment

		for C.AE_MetadataSearchMatch_NextSegment(cMatch, &cQueryStart, &cQueryEnd, &cAssetStart, &cAssetEnd, &cSegmentsPos) {
			segments = append(segments, &Segment{
				QueryStart: int64(cQueryStart),
				QueryEnd:   int64(cQueryEnd),
				AssetStart: int64(cAssetStart),
				AssetEnd:   int64(cAssetEnd),
			})
		}

		matches = append(matches, &MetadataSearchMatch{
			AssetID:   uint64(C.AE_MetadataSearchMatch_GetAssetID(cMatch)),
			AssetType: AssetType(C.AE_MetadataSearchMatch_GetAssetType(cMatch)),
			Segments:  segments,
		})
	}

	completedAtUnix := int64(C.AE_MetadataSearchResult_GetCompletedAt(cResult))

	return &MetadataSearchResult{
		LookupID:    uint64(C.AE_MetadataSearchResult_GetLookupID(cResult)),
		CompletedAt: time.Unix(completedAtUnix, 0),
		Matches:     matches,
	}, nil
}
