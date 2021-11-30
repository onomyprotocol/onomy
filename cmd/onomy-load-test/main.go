//go:build tmload
// +build tmload

package main

import (
	"github.com/onomyprotocol/tm-load-test/pkg/loadtest"
)

func main() {
	if err := loadtest.RegisterClientFactory("onomy-testing-app", &MyABCIAppClientFactory{}); err != nil {
		panic(err)
	}
	// The loadtest.Run method will handle CLI argument parsing, errors,
	// configuration, instantiating the load test and/or master/slave
	// operations, etc. All it needs is to know which client factory to use for
	// its load testing.
	loadtest.Run(&loadtest.CLIConfig{
		AppName:              "onomy-load-test",
		AppShortDesc:         "onomy testing application for My cosmos node",
		AppLongDesc:          "Some long description on how to use the tool",
		DefaultClientFactory: "onomy-testing-app",
	})
}
