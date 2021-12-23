// Package cmd contains cli wrapper for the onomy cli.
package cmd

import (
	gravitycmd "github.com/althea-net/cosmos-gravity-bridge/module/cmd/gravity/cmd"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/spm/cosmoscmd"

	"github.com/onomyprotocol/onomy/app"
)

// NewRootCmd initiates the cli for onomy chain.
func NewRootCmd() (*cobra.Command, cosmoscmd.EncodingConfig) {
	rootCmd, encodingConfig := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.Name,
		app.ModuleBasics,
		app.New,
	)

	// customize gentx and collect-gentxs commands
	for _, v := range rootCmd.Commands() {
		if v.Use == "gentx [key_name] [amount]" {
			rootCmd.RemoveCommand(v)
			rootCmd.AddCommand(gravitycmd.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome))
		}
		if v.Use == "collect-gentxs" {
			rootCmd.RemoveCommand(v)
			rootCmd.AddCommand(gravitycmd.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome))
		}
	}
	// eth_keys cmd
	rootCmd.AddCommand(gravitycmd.Commands(app.DefaultNodeHome))

	return rootCmd, encodingConfig
}
