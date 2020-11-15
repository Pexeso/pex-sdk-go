package ae

import "testing"

func TestMetadataSearch(t *testing.T) {
	c, err := NewClient("hurr", "durr")
	if err != nil {
		t.Fatalf("expected err to be nil, got: %v", err)
	}

	ft, err := NewFingerprintFromFile("")
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

	// TODO: check the rest
}
