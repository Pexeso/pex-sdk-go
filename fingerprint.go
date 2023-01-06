// Copyright 2020 Pexeso Inc. All rights reserved.

package pexae

// #include <stdlib.h>
// #include <pex/ae/sdk/lock.h>
// #include <pex/ae/sdk/fingerprint.h>
import "C"
import "unsafe"

// FingerprintType is a bit flag specifying one or more fingerprint types.
type FingerprintType int

const (
	FingerprintTypeVideo  FingerprintType = 1
	FingerprintTypeAudio  FingerprintType = 2
	FingerprintTypeMelody FingerprintType = 4
	FingerprintTypeAll                    = FingerprintTypeVideo | FingerprintTypeAudio | FingerprintTypeMelody
)

// Fingerprint is how the SDK identifies a piece of digital content.
// It can be generated from a media file or from a memory buffer. The
// content must be encoded in one of the supported formats and must be
// longer than 1 second.
type Fingerprint struct {
	b []byte
}

func NewFingerprint(b []byte) *Fingerprint {
	return &Fingerprint{
		b: b,
	}
}

// FingerprintFile is used to generate a fingerprint from a
// file stored on a disk. The path parameter must be a path
// to a valid file in supported format. The types parameter
// specifies which types of fingerprints to create. If not
// types are provided, FingerprintTypeAll is assumed.
func (x *Client) FingerprintFile(path string, types ...FingerprintType) (*Fingerprint, error) {
	return newFingerprint([]byte(path), true, reduceTypes(types))
}

// FingerprintBuffer is used to generate a fingerprint from a
// media file loaded in memory as a byte slice. The types parameter
// specifies which types of fingerprints to create. If not
// types are provided, FingerprintTypeAll is assumed.
func (x *Client) FingerprintBuffer(buffer []byte, types ...FingerprintType) (*Fingerprint, error) {
	return newFingerprint(buffer, false, reduceTypes(types))
}

func reduceTypes(in []FingerprintType) (out FingerprintType) {
	if len(in) == 0 {
		return FingerprintTypeAll
	}

	for _, t := range in {
		out = out | t
	}
	return out
}

func newFingerprint(input []byte, isFile bool, typ FingerprintType) (*Fingerprint, error) {
	C.AE_Lock()
	defer C.AE_Unlock()

	status := C.AE_Status_New()
	if status == nil {
		panic("out of memory")
	}
	defer C.AE_Status_Delete(&status)

	ft := C.AE_Buffer_New()
	if ft == nil {
		panic("out of memory")
	}

	if isFile {
		cFile := C.CString(string(input))
		defer C.free(unsafe.Pointer(cFile))

		C.AE_Fingerprint_File_For_Types(cFile, ft, status, C.int(typ))
	} else {
		buf := C.AE_Buffer_New()
		if buf == nil {
			panic("out of memory")
		}
		defer C.AE_Buffer_Delete(&buf)

		data := unsafe.Pointer(&input[0])
		size := C.size_t(len(input))

		C.AE_Buffer_Set(buf, data, size)
		C.AE_Fingerprint_Buffer_For_Types(buf, ft, status, C.int(typ))
	}

	if err := statusToError(status); err != nil {
		C.AE_Buffer_Delete(&ft)
		return nil, err
	}

	data := C.AE_Buffer_GetData(ft)
	size := C.int(C.AE_Buffer_GetSize(ft))

	return &Fingerprint{
		b: C.GoBytes(data, size),
	}, nil
}
