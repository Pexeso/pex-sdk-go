// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// Segment is the range [start, end) in both the query and the asset of
// where the match was found within the asset.
type Segment struct {
	// The start of the matched range int the query in seconds (inclusive).
	QueryStart int64

	// The end of the matched range in the query in seconds (exclusive).
	QueryEnd int64

	// The start of the matched range in the asset in seconds (inclusive).
	AssetStart int64

	// The end of the matched range in the asset in seconds (exclusive).
	AssetEnd int64
}
