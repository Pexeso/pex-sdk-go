// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/client.h>
// #include <pex/ae/sdk/asset.h>
// #include <pex/ae/sdk/stream_search.h>
// #include <stdlib.h>
import "C"
import "unsafe"

type StreamEventType int

const (
	// Sent when the search has started (before matches were found)
	StreamEventSearchStarted = StreamEventType(C.kStreamEventSearchStarted)

	// Sent when a match in the is first found.
	StreamEventMatchStarted = StreamEventType(C.kStreamEventMatchStarted)

	// Sent when a previously started match in the stream has ended.
	StreamEventMatchEnded = StreamEventType(C.kStreamEventMatchEnded)

	// Sent when input stream has ended.
	StreamEventStreamEnded = StreamEventType(C.kStreamEventStreamEnded)

	// Sent when the search has ended and no more events will be sent.
	StreamEventSearchEnded = StreamEventType(C.kStreamEventSearchEnded)

	// Sent when the search has failed with an error (e.g. downloading a stream chunk has failed)
	StreamEventSearchError = StreamEventType(C.kStreamEventSearchError)
)

func (x StreamEventType) String() string {
	switch x {
	case StreamEventSearchStarted:
		return "SearchStarted"
	case StreamEventMatchStarted:
		return "MatchStarted"
	case StreamEventMatchEnded:
		return "MatchEnded"
	case StreamEventStreamEnded:
		return "StreamEnded"
	case StreamEventSearchEnded:
		return "SearchEnded"
	case StreamEventSearchError:
		return "SearchError"
	}

	return "<unknown>"
}

type StreamEvent struct {
	Type           StreamEventType
	Err            error
	Asset          *Asset
	QueryTimestamp int64
	AssetTimestamp int64
}

type StreamSearch struct {
	c *C.AE_StreamSearch
}

func (x *Client) StartStreamSearch(url string) (*StreamSearch, error) {
	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))

	cSearch, err := C.AE_StreamSearch_New()
	if err != nil {
		panic("out of memory")
	}

	C.AE_StreamSearch_SetClient(cSearch, x.c)
	C.AE_StreamSearch_SetStreamPath(cSearch, cURL)

	C.AE_StreamSearch_StartSearch(cSearch, cStatus)
	if err := statusToError(cStatus); err != nil {
		C.AE_StreamSearch_Delete(&cSearch)
		return nil, err
	}

	return &StreamSearch{
		c: cSearch,
	}, nil
}

func (x *StreamSearch) Close() {
	C.AE_StreamSearch_EndSearch(x.c)
	C.AE_StreamSearch_Delete(&x.c)
}

func (x *StreamSearch) NextEvent() (*StreamEvent, error) {
	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cEvent, err := C.AE_StreamSearchEvent_New()
	if err != nil {
		panic("out of memory")
	}
	defer C.AE_StreamSearchEvent_Delete(&cEvent)

	C.AE_StreamSearch_GetNextEvent(x.c, cEvent, cStatus)
	if err := statusToError(cStatus); err != nil {
		return nil, err
	}

	event := &StreamEvent{
		Type: StreamEventType(C.AE_StreamSearchEvent_GetType(cEvent)),
	}

	switch event.Type {
	case StreamEventSearchError:
		getEventError(event, cEvent, cStatus)
	case StreamEventMatchStarted:
		getEventAsset(event, cEvent, cStatus)
	case StreamEventMatchEnded:
		getEventAsset(event, cEvent, cStatus)
	}
	return event, nil
}

func getEventError(event *StreamEvent, cEvent *C.AE_StreamSearchEvent, cStatus *C.AE_Status) {
	cErr := C.AE_Status_New()
	if cErr == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cErr)

	C.AE_StreamSearchEvent_GetError(cEvent, cErr, cStatus)
	if err := statusToError(cStatus); err != nil {
		panic(err)
	}
	event.Err = statusToError(cErr)
}

func getEventAsset(event *StreamEvent, cEvent *C.AE_StreamSearchEvent, cStatus *C.AE_Status) {
	cAsset := C.AE_Asset_New()
	if cAsset == nil {
		panic("out of memory")
	}
	defer C.AE_Asset_Delete(&cAsset)

	C.AE_StreamSearchEvent_GetAsset(cEvent, cAsset, cStatus)
	if err := statusToError(cStatus); err != nil {
		panic(err)
	}

	queryTimestamp := C.AE_StreamSearchEvent_GetQueryTimestamp(cEvent, cStatus)
	if err := statusToError(cStatus); err != nil {
		panic(err)
	}

	assetTimestamp := C.AE_StreamSearchEvent_GetAssetTimestamp(cEvent, cStatus)
	if err := statusToError(cStatus); err != nil {
		panic(err)
	}

	event.Asset = newAssetFromC(cAsset)
	event.QueryTimestamp = int64(queryTimestamp)
	event.AssetTimestamp = int64(assetTimestamp)
}
