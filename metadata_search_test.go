package ae

import (
	"fmt"
	"testing"
)

func TestMetadataSearch(t *testing.T) {
	c, err := NewClient("hurr", "durr")
	if err != nil {
		t.Fatalf("expected err to be nil, got: %v", err)
	}

	ft, err := NewFingerprintFromFile("/home/stepan/Downloads/Black_Hole_Sun.m4a")
	if err != nil {
		t.Fatalf("expected err to be nil, got: %v", err)
	}

	s := c.MetadataSearch()

	res, err := s.Do(&MetadataSearchRequest{
		Fingerprint: ft,
	})
	if err != nil {
		t.Fatalf("expected err to be nil, got: %v", err)
	}

	fmt.Printf("%+v\n", res)
}
