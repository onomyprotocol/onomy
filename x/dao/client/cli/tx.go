// Package cli contains dao module cli utils.
package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s dao subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2, // nolint:gomnd
		RunE:                       client.ValidateCmd,
	}

	return cmd
}

// GetCmdFundTreasuryProposalHandler implements the command to submit a fund-treasury proposal.
func GetCmdFundTreasuryProposalHandler() *cobra.Command {
	return &cobra.Command{
		Use:   "fund-treasury",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a fund treasury proposal.",
		Long:  strings.TrimSpace(""),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}

// GetCmdExchangeWithTreasuryProposalHandler implements the command to submit a exchange-with-treasury proposal.
func GetCmdExchangeWithTreasuryProposalHandler() *cobra.Command {
	return &cobra.Command{
		Use:   "exchange-with-treasury",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an exchange with treasury proposal.",
		Long:  strings.TrimSpace(""),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
