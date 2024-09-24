// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v1_1_4 //nolint:revive,stylecheck // app version

import (
	"context"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"

	"github.com/onomyprotocol/onomy/app/keepers"
)

// Name is migration name.
const Name = "v2.0.0"

// UpgradeHandler is an x/upgrade handler.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {

	return func(c context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		err = InitializeConsensusParamVersion(ctx, keepers.ConsensusParamsKeeper)
		if err != nil {
			// don't hard fail here, as this is not critical for the upgrade to succeed
			ctx.Logger().Error("Error initializing ConsensusParam Version:", "message", err.Error())
		}

		err = collectionsInitializeMintParam(ctx, keepers.MintKeeper)
		if err != nil {
			// don't hard fail here, as this is not critical for the upgrade to succeed
			ctx.Logger().Error("Error initializing mint params:", "message", err.Error())
		}
		return vm, nil
	}
}

func InitializeConsensusParamVersion(ctx sdk.Context, consensusKeeper consensusparamkeeper.Keeper) error {
	params, err := consensusKeeper.ParamsStore.Get(ctx)
	if err != nil {
		return err
	}
	params.Version = &cmtproto.VersionParams{}
	return consensusKeeper.ParamsStore.Set(ctx, params)
}

// ######### Here are all the params mint from mainnet
// # blocks_per_year: "6311520"
// # goal_bonded: "0.670000000000000000"
// # inflation_max: "0.200000000000000000"
// # inflation_min: "0.070000000000000000"
// # inflation_rate_change: "0.130000000000000000"
// # mint_denom: anom
func collectionsInitializeMintParam(ctx sdk.Context, mintKeeper mintkeeper.Keeper) error {
	params, err := mintKeeper.Params.Get(ctx)
	if err != nil {
		return err
	}
	params.BlocksPerYear = 6311520
	params.GoalBonded = math.LegacyMustNewDecFromStr("0.670000000000000000")
	params.InflationMax = math.LegacyMustNewDecFromStr("0.20000000000000000")
	params.InflationMin = math.LegacyMustNewDecFromStr("0.070000000000000000")
	params.InflationRateChange = math.LegacyMustNewDecFromStr("0.130000000000000000")
	params.MintDenom = "anom"
	return mintKeeper.Params.Set(ctx, params)
}
