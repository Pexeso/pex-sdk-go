package pexae

import "testing"

func TestClientWithValidCredentials(t *testing.T) {
	client, err := NewMockserverClient("client01", "secret01")
	if err != nil {
		t.Fatalf("expected no error, got %+v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Fatalf("closing the client returned error: %+v", err)
		}
	}()
	if client == nil {
		t.Fatal("expected a client, got nil")
	}
	if client.MetadataSearch == nil {
		t.Fatal("improperly initialized client, missing MetadataSearch")
	}
	if client.LicenseSearch == nil {
		t.Fatal("improperly initialized client, missing LicenseSearch")
	}
	if client.AssetLibrary == nil {
		t.Fatal("improperly initialized client, missing AssetLibrary")
	}
}

func TestClientWithInvalidCredentials(t *testing.T) {
	client, err := NewMockserverClient("invalid_client", "invalid_secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if client != nil {
		t.Fatalf("expected client to be nil, got %+v", client)
	}

	e, ok := err.(*Error)
	if !ok {
		t.Fatal("got error of unexpected type")
	}

	if e.Code != StatusUnauthenticated {
		t.Fatalf("got invalid error code, expected %d, got %d", StatusUnauthenticated, e.Code)
	}
}
