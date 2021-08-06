package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/onomyprotocol/onomy/x/onomy/types"
)

const (
	querySuggestionsMinimumDistance = 2
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(_ string) *cobra.Command {
	// Group onomy queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: querySuggestionsMinimumDistance,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	return cmd
}
