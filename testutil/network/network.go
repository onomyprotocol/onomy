// Package network contains utils to setup the chain network.
package network

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"

	"github.com/onomyprotocol/onomy/app"
	"github.com/onomyprotocol/onomy/testutil/simapp"
)

type (
	Network = network.Network // nolint:revive // simapp test var
	Config  = network.Config  // nolint:revive // simapp test var
)

// TestNetwork defines the test network wrapper.
type TestNetwork struct {
	*network.Network
}

// Option is an option pattern function used fot the test network customisations.
type Option func(*network.Config)

// WithGenesisOverride returns genesis override Option.
func WithGenesisOverride(override func(map[string]json.RawMessage) map[string]json.RawMessage) Option {
	return func(c *network.Config) {
		c.GenesisState = override(c.GenesisState)
	}
}

// New setups the test network.
func New(t *testing.T, opts ...Option) *TestNetwork {
	t.Helper()

	cfg := network.DefaultConfig()

	cfg.NumValidators = 1
	encCfg := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)
	cfg.GenesisState = app.ModuleBasics.DefaultGenesis(encCfg.Marshaler)
	cfg.AppConstructor = func(val network.Validator) servertypes.Application {
		onomyApp := simapp.Setup().OnomyApp()
		// the override is required in order not to face the issue with the gravity end blocker validation
		// because it requires the eth address to be linked with the validator account
		onomyApp.SetOrderEndBlockers(crisistypes.ModuleName, govtypes.ModuleName, stakingtypes.ModuleName)
		return onomyApp
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	onomyNetwork := network.New(t, cfg)

	_, err := onomyNetwork.WaitForHeight(1)
	require.NoError(t, err)

	return &TestNetwork{onomyNetwork}
}

// TxValidator1Args returns the tx params for the 1s network validator.
func (testNetwork *TestNetwork) TxValidator1Args() []string {
	return []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, testNetwork.Validators[0].Address.String()),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(testNetwork.Config.BondDenom, sdk.NewInt(10))).String()), // nolint:gomnd //test constant
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}
}

// Validator1Ctx returns the context of the 1st validator in the network.
func (testNetwork *TestNetwork) Validator1Ctx() client.Context {
	return testNetwork.Validators[0].ClientCtx
}

// Validator1Address returns the address of the 1st validator in the network.
func (testNetwork *TestNetwork) Validator1Address() sdk.AccAddress {
	return testNetwork.Validators[0].Address
}
