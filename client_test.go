package ae

import (
	"testing"
)

func TestClient(t *testing.T) {
	// invalid credentials
	_, err := NewClient("aaa", "bbb")
	if err == nil {
		t.Fatalf("expected err to not be nil, got nil")
	}

	// valid credentials
	_, err = NewClient("client01", "secret01")
	if err != nil {
		t.Fatalf("expected err to be nil, got: %v", err)
	}
}
