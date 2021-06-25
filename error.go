// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <pex/ae/sdk/c/status.h>
import "C"
import "fmt"

// StatusCode is used together with Error as a hint on why the error
// was returned.
type StatusCode int

const (
	StatusOK               = StatusCode(0)
	StatusDeadlineExceeded = StatusCode(1)
	StatusPermissionDenied = StatusCode(2)
	StatusUnauthenticated  = StatusCode(3)
	StatusNotFound         = StatusCode(4)
	StatusInvalidInput     = StatusCode(5)
	StatusOutOfMemory      = StatusCode(6)
	StatusInternalError    = StatusCode(7)
	StatusNotInitialized   = StatusCode(8)
	StatusConnectionError  = StatusCode(9)
	StatusLookupFailed     = StatusCode(10)
	StatusLookupTimedOut   = StatusCode(11)
)

// Error will be returend by most SDK functions. Besides an error
// message, it also includes a status code, which can be used to
// determine the underlying issue, e.g. AssetLibrary.GetAsset will return
// an error with StatusNotFound if the asset couldn't be found.
type Error struct {
	Code    StatusCode
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func statusToError(status *C.AE_Status) *Error {
	if !C.AE_Status_OK(status) {
		return &Error{
			Code:    StatusCode(C.AE_Status_GetCode(status)),
			Message: C.GoString(C.AE_Status_GetMessage(status)),
		}
	}
	return nil
}
