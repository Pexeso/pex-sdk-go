package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	pex "github.com/Pexeso/pex-sdk-go/v4"
	"os"
	"sort"
	"strings"
)

const (
	clientID     = "A9c30H96v6ke1iZdKGylThNUZxpH7FnC"
	clientSecret = "https://pwpush.com/p/wqiznrriw8o2-q/r"
	inputFile    = "/path/to/file.mp3"
)

type Segment struct {
	AssetStart int `json:"asset_start"`
	AssetEnd   int `json:"asset_end"`
	QueryStart int `json:"query_start"`
	QueryEnd   int `json:"query_end"`
	Confidence int `json:"confidence"`
}

type Match struct {
	Asset struct {
		ID          string  `json:"id"`
		Duration    float32 `json:"duration_seconds"`
		Title       string  `json:"title"`
		Subtitle    string  `json:"subtitle"`
		Artist      string  `json:"artist"`
		Isrc        string  `json:"isrc"`
		Label       string  `json:"label"`
		Distributor string  `json:"distributor"`
	} `json:"asset"`

	MatchDetails struct {
		Audio struct {
			Segments []Segment `json:"segments"`
		}
	} `json:"match_details"`
}

type Data struct {
	Matches []Match `json:"matches"`
}

func main() {
	// Initialize and authenticate the client.
	client, err := pex.NewPexSearchClient(clientID, clientSecret)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Optionally mock the client. If a client is mocked, it will only communicate
	// with the local mockserver instead of production servers. This is useful for
	// testing.
	if err := pex.MockClient(client); err != nil {
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

	dir := "06_12_2023_1"
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		ft, err := client.FingerprintFile(fmt.Sprintf("./%s/%s", dir, f.Name()))
		if err != nil {
			panic(err)
		}

		// Build the request.
		req := &pex.PexSearchRequest{
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

		var data *Data
		err = json.Unmarshal(j, &data)
		if err != nil {
			panic(err)
		}

		//Write the CSV data
		ff, err := os.Create(fmt.Sprintf("%s.csv", strings.TrimSuffix(f.Name(), ".flac")))
		if err != nil {
			panic(err)
		}
		defer ff.Close()

		mWriter := csv.NewWriter(ff)
		mWriter.UseCRLF = true
		defer mWriter.Flush()

		y := data.Matches
		sort.Slice(y, func(i, j int) bool {
			return y[i].MatchDetails.Audio.Segments[0].QueryStart < y[j].MatchDetails.Audio.Segments[0].QueryStart
		})

		// write
		for _, m := range y {
			for _, s := range m.MatchDetails.Audio.Segments {
				filename := "pretest.flac"
				var confidence int
				if s.Confidence == 100 {
					confidence = 99
				} else {
					confidence = int(s.Confidence)
				}

				r := []string{
					filename,
					m.Asset.Title,
					m.Asset.Subtitle,
					m.Asset.Artist,
					m.Asset.Isrc,
					fmt.Sprintf("%.1f", float32(s.QueryStart)),
					fmt.Sprintf("%.1f", float32(s.QueryEnd)),
					fmt.Sprintf("%.1f", float32(s.AssetStart)),
					fmt.Sprintf("%d", confidence),
				}
				wErr := mWriter.Write(r)
				if wErr != nil {
					panic(wErr)
				}

			}
		}

		if f.Name() == "06_12_2023_152143.combined_stream.flac" {
			break
		}
	}

}
