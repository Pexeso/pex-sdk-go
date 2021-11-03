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

// FingerprintFile is used to generate a fingerprint from a
// file stored on a disk. The parameter to the function must be a path
// to a valid file in supported format.
func (x *Client) FingerprintFile(path string) (*Fingerprint, error) {
	return newFingerprint([]byte(path), true)
}

// FingerprintBuffer is used to generate a fingerprint from a
// media file loaded in memory as a byte slice.
func (x *Client) FingerprintBuffer(buffer []byte) (*Fingerprint, error) {
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

// LoadFingerprint loads a fingerprint previously serialized by
// the Fingerprint.Dump() function.
func (x *Client) LoadFingerprint(dump []byte) (*Fingerprint, error) {
	status := C.AE_Status_New()
	if status == nil {
		panic("out of memory")
	}

	defer C.AE_Status_Delete(&status)
	ft := C.AE_Fingerprint_New()
	if ft == nil {
		panic("out of memory")
	}

	b := C.AE_Buffer_New()
	if b == nil {
		panic("out of memory")
	}
	defer C.AE_Buffer_Delete(&b)

	cDump := C.CBytes(dump)
	defer C.free(cDump)

	C.AE_Buffer_Set(b, cDump, C.size_t(len(dump)))

	C.AE_Fingerprint_Load(ft, b, status)
	if err := statusToError(status); err != nil {
		return nil, err
	}
	return &Fingerprint{ft}, nil
}
