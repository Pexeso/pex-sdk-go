package ae

import "testing"

func TestClientInitialization(t *testing.T) {
	c, err := NewClient("hurr", "durr")
	if err != nil {
		t.Fatalf("expected err == nil, got: %v", err)
	}

	s := c.MetadataSearch()

	s.Do(&MetadataSearchRequest{})
}
