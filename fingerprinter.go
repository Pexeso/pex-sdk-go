// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <stdlib.h>
// #include <pex/ae/sdk/fingerprint.h>
import "C"
import (
	"errors"
	"unsafe"
)

type Fingerprinter struct {
	embedded bool
}

// NewFingerprintFromFile is used to generate a fingerprint from a
// file stored on a disk. The parameter to the function must be a path
// to a valid file in supported format.
func (x *Fingerprinter) FromFile(path string) (*Fingerprint, error) {
	if !x.embedded {
		return nil, errors.New("use Client.Fingerprinter instead of creating a new one")
	}
	return newFingerprint([]byte(path), true)
}

// NewFingerprintFromBuffer is used to generate a fingerprint from a
// media file loaded in memory as a byte slice.
func (x *Fingerprinter) FromBuffer(buffer []byte) (*Fingerprint, error) {
	if !x.embedded {
		return nil, errors.New("use Client.Fingerprinter instead of creating a new one")
	}
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

// LoadDumpedfingerprint loads a fingerprint previously serialized by
// the Fingerprint.Dump() function.
func (x *Fingerprinter) Load(dump []byte) (*Fingerprint, error) {
	if !x.embedded {
		return nil, errors.New("use Client.Fingerprinter instead of creating a new one")
	}

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
