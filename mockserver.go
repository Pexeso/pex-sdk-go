package pexae

// #cgo pkg-config: pexae
// #include <pex/ae/sdk/c/client.h>
// #include <pex/ae/sdk/c/mockserver.h>
// #include <stdlib.h>
import "C"
import "unsafe"

// NewMockserverClient creates a new instance of the client that will  using
// provided credentials for authentication.
func NewMockserverClient(clientID, clientSecret string) (*Client, error) {
	cStatus := C.AE_Status_New()
	if cStatus == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&cStatus)

	cClientID := C.CString(clientID)
	defer C.free(unsafe.Pointer(cClientID))

	cClientSecret := C.CString(clientSecret)
	defer C.free(unsafe.Pointer(cClientSecret))

	cClient := C.AE_Client_New()
	if cClient == nil {
		panic("out of memory")
	}

	C.AE_Mockserver_InitClient(cClient, cClientID, cClientSecret, cStatus)
	if err := statusToError(cStatus); err != nil {
		C.free(unsafe.Pointer(cClient))
		return nil, err
	}
	return buildClient(cClient), nil
}
