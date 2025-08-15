// Copyright 2020 Pexeso Inc. All rights reserved.

// Package pex contains the Go bindings for the Pex SDK.
//
// Important! Please make sure to install the core library, as described in the
// following link: https://docs.search.pex.com/installation/, before trying to use
// the Go bindings.
//
// # Installation
//
// You can install the Go language bindings like this:
//
//	go get github.com/Pexeso/pex-sdk-go/v4
//
// # Pex search example
//
//	package main
//
//	import (
//		"encoding/json"
//		"fmt"
//
//		pex "github.com/Pexeso/pex-sdk-go/v4"
//	)
//
//	const (
//		clientID     = "#YOUR_CLIENT_ID_HERE"
//		clientSecret = "#YOUR_CLIENT_SECRET_HERE"
//		inputFile    = "/path/to/file.mp3"
//	)
//
//	func main() {
//		// Initialize and authenticate the client.
//		client, err := pex.NewPexSearchClient(clientID, clientSecret)
//		if err != nil {
//			panic(err)
//		}
//		defer client.Close()
//
//		// Fingerprint a file. You can also fingerprint a buffer with
//		//
//		//   client.FingerprintBuffer([]byte).
//		//
//		// Both the files and the memory buffers
//		// must be valid media content in following formats:
//		//
//		//   Audio: aac
//		//   Video: h264, h265
//		//
//		// Keep in mind that generating a fingerprint is CPU bound operation and
//		// might consume a significant amount of your CPU time.
//		ft, err := client.FingerprintFile(inputFile)
//		if err != nil {
//			panic(err)
//		}
//
//		// Build the request.
//		req := &pex.PexSearchRequest{
//			Fingerprint: ft,
//		}
//
//		// Start the search.
//		fut, err := client.StartSearch(req)
//		if err != nil {
//			panic(err)
//		}
//
//		// Retrieve the result.
//		res, err := fut.Get()
//		if err != nil {
//			panic(err)
//		}
//
//		// Print the result.
//		j, err := json.MarshalIndent(res, "", "  ")
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(string(j))
//	}
//
// # Private search
//
//	package main
//
//	import (
//		"encoding/json"
//		"fmt"
//
//		pex "github.com/Pexeso/pex-sdk-go/v4"
//	)
//
//	const (
//		clientID     = "#YOUR_CLIENT_ID_HERE"
//		clientSecret = "#YOUR_CLIENT_SECRET_HERE"
//		inputFile    = "/path/to/file.mp3"
//	)
//
//	func main() {
//		// Initialize and authenticate the client.
//		client, err := pex.NewPrivateSearchClient(clientID, clientSecret)
//		if err != nil {
//			panic(err)
//		}
//		defer client.Close()
//
//		// Fingerprint a file. You can also fingerprint a buffer with
//		//
//		//   client.FingerprintBuffer([]byte).
//		//
//		// Both the files and the memory buffers
//		// must be valid media content in following formats:
//		//
//		//   Audio: aac
//		//   Video: h264, h265
//		//
//		// Keep in mind that generating a fingerprint is CPU bound operation and
//		// might consume a significant amount of your CPU time.
//		ft, err := client.FingerprintFile(inputFile)
//		if err != nil {
//			panic(err)
//		}
//
//		// Ingest it into your private catalog.
//		if err := client.Ingest("my_id_1", ft); err != nil {
//			panic(err)
//		}
//
//		// Build the request.
//		req := &pex.PrivateSearchRequest{
//			Fingerprint: ft,
//		}
//
//		// Start the search.
//		fut, err := client.StartSearch(req)
//		if err != nil {
//			panic(err)
//		}
//
//		// Retrieve the result.
//		res, err := fut.Get()
//		if err != nil {
//			panic(err)
//		}
//
//		// Print the result.
//		j, err := json.MarshalIndent(res, "", "  ")
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(string(j))
//	}
package pex
