// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <stdlib.h>
// #include <pex/ae/sdk/fingerprint.h>
import "C"
import "unsafe"

// Fingerprint is how the SDK identifies a piece of digital content.
// It can be generated from a media file or from a memory buffer. The
// content must be encoded in one of the supported formats and must be
// longer than 1 second.
type Fingerprint struct {
	ft *C.AE_Fingerprint
}

// NewFingerprintFromFile is used to generate a fingerprint from a
// file stored on a disk. The parameter to the function must be a path
// to a valid file in supported format.
func NewFingerprintFromFile(path string) (*Fingerprint, error) {
	return newFingerprint([]byte(path), true)
}

// NewFingerprintFromBuffer is used to generate a fingerprint from a
// media file loaded in memory as a byte slice.
func NewFingerprintFromBuffer(buffer []byte) (*Fingerprint, error) {
	return newFingerprint(buffer, false)
}

func newFingerprint(input []byte, isFile bool) (*Fingerprint, error) {
	status := C.AE_Status_New()
	if status == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&status)

	ft := C.AE_Fingerprint_New()
	if ft == nil {
		panic("out of memory")
	}

	if isFile {
		cFile := C.CString(string(input))
		defer C.free(unsafe.Pointer(cFile))

		C.AE_Fingerprint_FromFile(ft, cFile, status)
	} else {
		buffer := C.AE_Buffer_New()
		if buffer == nil {
			panic("out of memory")
		}
		defer C.AE_Buffer_Delete(&buffer)

		cInput := C.CBytes(input)
		defer C.free(cInput)

		C.AE_Buffer_Set(buffer, cInput, C.size_t(len(input)))
		C.AE_Fingerprint_FromBuffer(ft, buffer, status)
	}

	if err := statusToError(status); err != nil {
		C.AE_Fingerprint_Delete(&ft)
		return nil, err
	}

	return &Fingerprint{ft}, nil
}

// Close releases allocated resources and memory.
func (f *Fingerprint) Close() error {
	C.AE_Fingerprint_Delete(&f.ft)
	return nil
}
