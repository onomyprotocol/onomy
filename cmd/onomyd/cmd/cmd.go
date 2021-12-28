// Package cmd contains cli wrapper for the onomy cli.
package cmd

import (
	gravitycmd "github.com/althea-net/cosmos-gravity-bridge/module/cmd/gravity/cmd"
	cosmossimappcmd "github.com/cosmos/cosmos-sdk/simapp/simd/cmd"
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
	cmdsToReplace := map[string]*cobra.Command{
		"gentx [key_name] [amount]": gravitycmd.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		"collect-gentxs":            gravitycmd.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		"add-genesis-account [address_or_key_name] [coin][,[coin]]": cosmossimappcmd.AddGenesisAccountCmd(app.DefaultNodeHome),
	}

	for _, v := range rootCmd.Commands() {
		cmd, ok := cmdsToReplace[v.Use]
		if ok {
			rootCmd.RemoveCommand(v)
			rootCmd.AddCommand(cmd)
			delete(cmdsToReplace, v.Use)
		}
	}
	if len(cmdsToReplace) != 0 {
		panic("on onomy cmd replacements, not all of the commands were replaced")
	}
	// eth_keys cmd
	rootCmd.AddCommand(gravitycmd.Commands(app.DefaultNodeHome))

	return rootCmd, encodingConfig
}
