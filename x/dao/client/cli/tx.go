// Package cli contains dao module cli utils.
package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2, // nolint:gomnd
		RunE:                       client.ValidateCmd,
	}

	return cmd
}
