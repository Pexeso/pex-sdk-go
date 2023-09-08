// Copyright 2020 Pexeso Inc. All rights reserved.

package pex

import (
	"encoding/json"
	"errors"
)

// SegmentType shows whether the segment matched on audio, video or melody.
type SegmentType int

func (x SegmentType) String() string {
	switch x {
	case SegmentTypeUnspecified:
		return "unspecified"
	case SegmentTypeAudio:
		return "audio"
	case SegmentTypeVideo:
		return "video"
	case SegmentTypeMelody:
		return "melody"
	default:
		return "unknown"
	}
}

func ParseSegmentType(s string) (SegmentType, error) {
	switch s {
	case "unspecified":
		return SegmentTypeUnspecified, nil
	case "audio":
		return SegmentTypeAudio, nil
	case "video":
		return SegmentTypeVideo, nil
	case "melody":
		return SegmentTypeMelody, nil
	default:
		return SegmentTypeUnspecified, errors.New("unknown segment type")
	}
}

func (x SegmentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}

func (x *SegmentType) UnmarshalJSON(data []byte) (err error) {
	var segment string
	if err := json.Unmarshal(data, &segment); err != nil {
		return err
	}
	if *x, err = ParseSegmentType(segment); err != nil {
		return err
	}
	return nil
}

const (
	SegmentTypeUnspecified = SegmentType(0)
	SegmentTypeAudio       = SegmentType(1)
	SegmentTypeVideo       = SegmentType(2)
	SegmentTypeMelody      = SegmentType(3)
)

// Segment is the range [start, end) in both the query and the asset of
// where the match was found within the asset.
type Segment struct {
	// Type of the segment (audio, video, melody).
	Type SegmentType

	// The start of the matched range int the query in seconds (inclusive).
	QueryStart int64

	// The end of the matched range in the query in seconds (exclusive).
	QueryEnd int64

	// The start of the matched range in the asset in seconds (inclusive).
	AssetStart int64

	// The end of the matched range in the asset in seconds (exclusive).
	AssetEnd int64
}
