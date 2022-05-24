package main

import (
	"testing"
	"time"

	"github.com/onomyprotocol/onomy/testutil/integration"
)

var bootstrappingTimeout = time.Minute // nolint:gochecknoglobals

func TestInitAndRunChain(t *testing.T) {
	// run onomy chain
	onomyChain, err := integration.NewOnomyChain()
	if err != nil {
		t.Fatal(err)
	}

	if err := onomyChain.Start(bootstrappingTimeout); err != nil {
		t.Fatal(err)
	}
}
