package main

import (
	"testing"
	"time"

	"github.com/onomyprotocol/onomy/testutil"
)

var bootstrappingTimeout = time.Minute // nolint:gochecknoglobals

func TestInitAndRunChain(t *testing.T) {
	// run onomy chain
	onomyChain, err := testutil.NewOnomyChain()
	if err != nil {
		t.Fatal(err)
	}

	if err := onomyChain.Start(bootstrappingTimeout); err != nil {
		t.Fatal(err)
	}
}
