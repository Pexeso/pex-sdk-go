// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <stdlib.h>
// #include <pex/ae/sdk/fingerprint.h>
import "C"

// Fingerprint is how the SDK identifies a piece of digital content.
// It can be generated from a media file or from a memory buffer. The
// content must be encoded in one of the supported formats and must be
// longer than 1 second.
type Fingerprint struct {
	ft *C.AE_Fingerprint
}

// Close releases allocated resources and memory.
func (f *Fingerprint) Close() error {
	C.AE_Fingerprint_Delete(&f.ft)
	return nil
}

// Dump serializes the fingerprint into a byte slice so that it can be
// stored on a disk or in a dabase. It can later be deserialized with
// the LoadDumpedFingerprint() function.
func (f *Fingerprint) Dump() ([]byte, error) {
	status := C.AE_Status_New()
	if status == nil {
		panic("out of memory")
	}

	defer C.AE_Status_Delete(&status)
	b := C.AE_Buffer_New()
	if b == nil {
		panic("out of memory")
	}
	defer C.AE_Buffer_Delete(&b)

	C.AE_Fingerprint_Dump(f.ft, b, status)
	if err := statusToError(status); err != nil {
		return nil, err
	}

	data := C.AE_Buffer_GetData(b)
	size := C.int(C.AE_Buffer_GetSize(b))

	return C.GoBytes(data, size), nil
}
