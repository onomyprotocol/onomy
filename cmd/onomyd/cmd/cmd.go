// Package cmd contains cli wrapper for the onomy cli.
package cmd

import (
	"fmt"

	gravitycmd "github.com/althea-net/cosmos-gravity-bridge/module/cmd/gravity/cmd"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cosmossimappcmd "github.com/cosmos/cosmos-sdk/simapp/simd/cmd"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/spm/cosmoscmd"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/onomyprotocol/onomy/app"
)

const (
	gravityName = "gravity"
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

	cmdsToReplace := map[string]*cobra.Command{
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

	// gravity cmd wrapper
	rootCmd.AddCommand(WrapBridgeCommands(app.DefaultNodeHome, gravityName, []*cobra.Command{
		gravitycmd.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		gravitycmd.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
	}))

	return rootCmd, encodingConfig
}

// WrapBridgeCommands registers a sub-tree of gravity commands.
func WrapBridgeCommands(defaultNodeHome, rootCmd string, cmds []*cobra.Command) *cobra.Command {
	//nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:   rootCmd,
		Short: fmt.Sprintf("Manage %s bridge.", rootCmd),
		Long:  fmt.Sprintf("Manage %s bridge.", rootCmd),
	}

	for _, childCmd := range cmds {
		cmd.AddCommand(childCmd)
	}

	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().String(cli.OutputFlag, "text", "Output format (text|json)")

	return cmd
}
