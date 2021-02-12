[![docs](https://img.shields.io/badge/docs-reference-blue.svg)](https://docs.ae.pex.com/go/)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)

# Attribution Engine SDK for Golang

Go bindings for the [Attribution Engine SDK](https://docs.ae.pex.com).

### Installation

You can install the Go language bindings like this:

    go get github.com/Pexeso/ae-sdk-go


### Fingerprinting

A fingerprint is how the SDK identifies a piece of digital content.
It can be generated from a media file or from a memory buffer. The
content must be encoded in one of the supported formats and must be
longer than 1 second.

You can generate a fingerprint from a media file:

```go
ft, err := ae.NewFingerprintFromFile("/path/to/file.mp4")
if err != nil {
    panic(err)
}
defer ft.Close()

// ...
```

Or you can generate a fingerprint from a memory buffer:

```go
b, _ := ioutil.ReadFile("/path/to/file.mp4")

ft, err := ae.NewFingerprintFromBuffer(b)
if err != nil {
    panic(err)
}
defer ft.Close()

// ...
```

Both the files and the memory buffers must be valid media content in
following formats:

```
Audio: aac
Video: h264, h265
```

Keep in mind that generating a fingerprint is CPU bound operation and
might consume a significant amount of your CPU time.


### Metadata search

After the fingerprint is generated, you can use it to perform a metadata search.

```go
// First, you need to initialize a client:
client, err := ae.NewClient(clientID, clientSecret)
if err != nil {
    panic(err)
}
defer client.Close()

// Build the request.
req := &ae.MetadataSearchRequest{
    Fingerprint: ft,
}

// Start the search.
fut, err := client.MetadataSearch.Start(req)
if err != nil {
    panic(err)
}

// Do other stuff.

// Retrieve the result.
res, err := fut.Get()
if err != nil {
    panic(err)
}

// Print the result.
fmt.Printf("%+v\n", res)
```


### License search

Performing a license search is very similar to metadata search.

```go
// ...

// Build the request.
req := &ae.LicenseSearchRequest{
    Fingerprint: ft,
}

// Start the search.
fut, err := client.LicenseSearch.Start(req)
if err != nil {
    panic(err)
}

// ...
```

The most significant difference between the searches currently is in the
results they return. See MetadataSearchResult and LicenseSearchResult for
more information.


### Asset library

You can use AssetLibrary to retrieve information about matched assets.

```go
// After successful metadata search.
for _, match := range res.Matches {
    asset, err := client.AssetLibrary.GetAsset(match.AssetID)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", asset)
}
```
