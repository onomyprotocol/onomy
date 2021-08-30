package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/onomyprotocol/onomy/cmd/onomyd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		switch e := err.(type) { // nolint:errorlint
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
