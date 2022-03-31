// Package network contains utils to setup the chain network.
package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmdb "github.com/tendermint/tm-db"

	"github.com/onomyprotocol/onomy/app"
)

type (
	Network = network.Network // nolint:revive // simapp test var
	Config  = network.Config  // nolint:revive // simapp test var
	// TestNetwork defines the test network wrapper.
	TestNetwork struct {
		*network.Network
	}
)

// SetupNetwork setups the test network.
func SetupNetwork(t *testing.T) *TestNetwork {
	t.Helper()

	onomyNetwork := New(t)

	_, err := onomyNetwork.WaitForHeight(1)
	if err != nil {
		panic(err)
	}

	return &TestNetwork{onomyNetwork}
}

// Validator1Ctx returns the context of the 1st validator in the network.
func (testNetwork *TestNetwork) Validator1Ctx() client.Context {
	return testNetwork.Validators[0].ClientCtx
}

// Validator1Address returns the address of the 1st validator in the network.
func (testNetwork *TestNetwork) Validator1Address() sdk.AccAddress {
	return testNetwork.Validators[0].Address
}

// New creates instance with fully configured cosmos network.
// Accepts optional config, that will be used in place of the DefaultConfig() if provided.
func New(t *testing.T, configs ...network.Config) *network.Network {
	t.Helper()

	if len(configs) > 1 {
		panic("at most one config should be provided")
	}
	var cfg network.Config
	if len(configs) == 0 {
		cfg = DefaultConfig()
	} else {
		cfg = configs[0]
	}
	net := network.New(t, cfg)
	t.Cleanup(net.Cleanup)
	return net
}

// DefaultConfig will initialize config for the network with custom application,
// genesis and single validator. All other parameters are inherited from cosmos-sdk/testutil/network.DefaultConfig.
func DefaultConfig() network.Config {
	encoding := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)
	return network.Config{
		Codec:             encoding.Marshaler,
		TxConfig:          encoding.TxConfig,
		LegacyAmino:       encoding.Amino,
		InterfaceRegistry: encoding.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.Validator) servertypes.Application {
			return app.New(
				val.Ctx.Logger, tmdb.NewMemDB(), nil, true, map[int64]bool{}, val.Ctx.Config.RootDir, 0,
				encoding,
				simapp.EmptyAppOptions{},
				baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
				baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
			)
		},
		GenesisState:    app.ModuleBasics.DefaultGenesis(encoding.Marshaler),
		TimeoutCommit:   2 * time.Second,                    // nolint:gomnd // simapp param
		ChainID:         "chain-" + tmrand.NewRand().Str(6), // nolint:gomnd // simapp param
		NumValidators:   3,                                  // nolint:gomnd // default validators number
		BondDenom:       sdk.DefaultBondDenom,
		MinGasPrices:    fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000000000000000, sdk.DefaultPowerReduction), // nolint:gomnd // simapp param
		StakingTokens:   sdk.TokensFromConsensusPower(500000000000000, sdk.DefaultPowerReduction),  // nolint:gomnd // simapp param
		BondedTokens:    sdk.TokensFromConsensusPower(100000000000000, sdk.DefaultPowerReduction),  // nolint:gomnd // simapp param
		PruningStrategy: storetypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions:  []keyring.Option{},
	}
}
