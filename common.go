// Copyright 2020 Pexeso Inc. All rights reserved.

package pex

type MatchDetails struct {
	Audio  *SegmentDetails `json:"audio"`
	Melody *SegmentDetails `json:"melody"`
	Video  *SegmentDetails `json:"video"`
}

type SegmentDetails struct {
	QueryMatchDurationSeconds float32 `json:"query_match_duration_seconds"`
	QueryMatchPercentage      float32 `json:"query_match_percentage"`
	AssetMatchDurationSeconds float32 `json:"asset_match_duration_seconds"`
	AssetMatchPercentage      float32 `json:"asset_match_percentage"`

	Segments []Segment `json:"segments"`
}

// Segment is the range [start, end) in both the query and the asset of
// where the match was found within the asset.
type Segment struct {
	// The start of the matched range int the query in seconds (inclusive).
	QueryStart int64 `json:"query_start"`

	// The end of the matched range in the query in seconds (exclusive).
	QueryEnd int64 `json:"query_end"`

	// The start of the matched range in the asset in seconds (inclusive).
	AssetStart int64 `json:"asset_start"`

	// The end of the matched range in the asset in seconds (exclusive).
	AssetEnd int64 `json:"asset_end"`

	AudioPitch          *int64 `json:"audio_pitch"`
	AudioSpeed          *int64 `json:"audio_speed"`
	MelodyTransposition *int64 `json:"melody_transposition"`

	Confidence int64 `json:"confidence"`

	DebugInfo any `json:"debug_info,omitempty"`
}

type DSP struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
