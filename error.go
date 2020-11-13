// Copyright 2020 Pexeso Inc. All rights reserved.

package ae

// #include <pex/aesdk/c/status.h>
import "C"
import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func statusToError(status *C.AE_Status) *Error {
	if !C.AE_Status_OK(status) {
		return &Error{
			Code:    int(C.AE_Status_GetCode(status)),
			Message: C.GoString(C.AE_Status_GetMessage(status)),
		}
	}
	return nil
}
