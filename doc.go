// Copyright 2020 Pexeso Inc. All rights reserved.

// Package ae
//
// Fingerprinting
//
// A fingerprint is how the SDK identifies a piece of digital content.
// It can be generated from a media file or from a memory buffer. The
// content must be encoded in one of the supported formats and must be
// longer than 1 second.
//
// You can generate a fingerprint from a media file:
//
//     ft, err := ae.NewFingerprintFromFile("/path/to/file.mp4")
//     if err != nil {
//       panic(err)
//     }
//     defer ft.Close()
//
//     // ...
//
// Or you can generate a fingerprint from a memory buffer:
//
//     ft, err := ae.NewFingerprintFromBuffer(buf)
//     if err != nil {
//       panic(err)
//     }
//     defer ft.Close()
//
//     // ...
//
// Both the files and the memory buffers must be valid media content in
// following formats:
//
//     Audio: aac
//     Video: h264, h265
//
// Keep in mind that generating a fingerprint is CPU bound operation and
// might consume a significant amount of your CPU time.
//
//
// Lookup
//
// A lookup is a search that takes a generated fingerprint and compares it to the
// Pex Asset Registry with the goal of identifying matches.
//
// To perform a lookup:
//
//     client := ae.NewType2Client("__client__", "__secret__")
//
//     req := &Type2LookupRequest{
//       Fingerprint: ft,
//       Metadata: &Metadata{
//         UPCs: []string{123},
//       },
//     }
//
//     res, err := client.Lookup(req)
//     if err != nil {
//       panic(err)
//     }
//
//     fmt.Println("lookup_id:", res.LookupID)
//
// You can also retrieve lookup detail using the GetDetail() function like this:
//
//     req := &Type2DetailRequest{
//       LookupID: lookupID,
//     }
//
//     res, err := client.GetDetail(req)
//     if err != nil {
//       panic(err)
//     }
//
//     fmt.Println("detail": res)
//
//
package ae
