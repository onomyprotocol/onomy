package main

import (
	"os"

	gravitycmd "github.com/althea-net/cosmos-gravity-bridge/module/cmd/gravity/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/spm/cosmoscmd"

	"github.com/onomyprotocol/onomy/app"
)

func main() {
	encodingConfig := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)
	rootCmd, _ := cosmoscmd.NewRootCmd(
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

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
