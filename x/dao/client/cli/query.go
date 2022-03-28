package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(_ string) *cobra.Command {
	// Group dao queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2, // nolint:gomnd
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())

	return cmd
}
