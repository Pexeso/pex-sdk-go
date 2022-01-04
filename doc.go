// Copyright 2020 Pexeso Inc. All rights reserved.

// Welcome to the Go bindings reference manual for the Attribution Engine's SDK.
//
// Important! Please make sure to install the core library, as described in the
// following link: https://docs.ae.pex.com/installation/, before trying to use
// the Go bindings.
//
//
// Installation
//
// You can install the Go language bindings like this:
//
//     go get github.com/Pexeso/ae-sdk-go/v2
//
//
// Client
//
// Before you can do any operation with the SDK you need to initialize a client.
//
//     client, err := pexae.NewClient(clientID, clientSecret)
//     if err != nil {
//         panic(err)
//     }
//     defer client.Close()
//
// If you want to test the SDK using the mockserver you need to mock the client:
//
//     if err := pexae.MockClient(client); err != nil {
//         panic(err)
//     }
//
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
//     ft, err := client.FingerprintFile("/path/to/file.mp4")
//     if err != nil {
//         panic(err)
//     }
//
// Or you can generate a fingerprint from a memory buffer:
//
//     b, _ := ioutil.ReadFile("/path/to/file.mp4")
//
//     ft, err := pexae.FingerprintBuffer(b)
//     if err != nil {
//         panic(err)
//     }
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
// Metadata search
//
// After the fingerprint is generated, you can use it to perform a metadata search.
//
//     // Build the request.
//     req := &pexae.MetadataSearchRequest{
//         Fingerprint: ft,
//     }
//
//     // Start the search.
//     fut, err := client.StartMetadataSearch(req)
//     if err != nil {
//         panic(err)
//     }
//
//     // Do other stuff.
//
//     // Retrieve the result.
//     res, err := fut.Get()
//     if err != nil {
//         panic(err)
//     }
//
//     // Print the result.
//     fmt.Printf("%+v\n", res)
//
//
// License search
//
// Performing a license search is very similar to metadata search.
//
//     // ...
//
//     // Build the request.
//     req := &pexae.LicenseSearchRequest{
//         Fingerprint: ft,
//     }
//
//     // Start the search.
//     fut, err := client.StartLicenseSearch(req)
//     if err != nil {
//         panic(err)
//     }
//
//     // ...
//
// The most significant difference between the searches currently is in the
// results they return. See MetadataSearchResult and LicenseSearchResult for
// more information.
package pexae
