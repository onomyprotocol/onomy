// Package cli is onomy cosmos sdk cli.
package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	// "github.com/cosmos/cosmos-sdk/client/flags".
	"github.com/onomyprotocol/onomy/x/onomy/types"
)

const (
	txSuggestionsMinimumDistance = 2
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: txSuggestionsMinimumDistance,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	return cmd
}
