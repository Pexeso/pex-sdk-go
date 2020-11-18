package ae

import (
	"fmt"
	"testing"
)

func TestAssetLibrary(t *testing.T) {
	c, err := NewClient("client01", "secret01")
	if err != nil {
		t.Fatalf("expected err to be nil, got: %v", err)
	}

	a := c.AssetLibrary()

	asset, err := a.GetAsset(501)
	if err != nil {
		t.Fatalf("expected err to be nil, got: %v", err)
	}

	fmt.Printf("%+v\n", asset)
	fmt.Printf("%+v\n", asset.Metadata)
}
