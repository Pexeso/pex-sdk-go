package main

import (
	"encoding/json"
	"fmt"

	pexae "github.com/Pexeso/ae-sdk-go/v3"
)

const (
	clientID     = "#YOUR_CLIENT_ID_HERE"
	clientSecret = "#YOUR_CLIENT_SECRET_HERE"
	inputFile    = "/path/to/file.mp3"
)

func main() {
	// Initialize and authenticate the client.
	client, err := pexae.NewPexSearchClient(clientID, clientSecret)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Optionally mock the client. If a client is mocked, it will only communicate
	// with the local mockserver instead of production servers. This is useful for
	// testing.
	if err := pexae.MockClient(client); err != nil {
		panic(err)
	}

	// Fingerprint a file. You can also fingerprint a buffer with
	//
	//   client.FingerprintBuffer([]byte).
	//
	// Both the files and the memory buffers
	// must be valid media content in following formats:
	//
	//   Audio: aac
	//   Video: h264, h265
	//
	// Keep in mind that generating a fingerprint is CPU bound operation and
	// might consume a significant amount of your CPU time.
	ft, err := client.FingerprintFile(inputFile)
	if err != nil {
		panic(err)
	}

	// Build the request.
	req := &pexae.PexSearchRequest{
		Fingerprint: ft,
	}

	// Start the search.
	fut, err := client.StartSearch(req)
	if err != nil {
		panic(err)
	}

	// Retrieve the result.
	res, err := fut.Get()
	if err != nil {
		panic(err)
	}

	// Print the result.
	j, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
