// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v2_0_0 //nolint:revive,stylecheck // app version

import (
	"context"
	"time"

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

		unbondingfails, err := getInfoUnbondingFail(ctx, keepers)
		if err != nil {
			ctx.Logger().Error("Error fixUnbondingICSRemove:", "message", err.Error())
		}

		vm, err = mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		err = unbondNow(ctx, keepers, unbondingfails)
		if err != nil {
			ctx.Logger().Error("Error fixUnbondingICSRemove:", "message", err.Error())
		}

		err = initializeConsensusParamVersion(ctx, keepers.ConsensusParamsKeeper)
		if err != nil {
			// don't hard fail here, as this is not critical for the upgrade to succeed
			ctx.Logger().Error("Error initializing ConsensusParam Version:", "message", err.Error())
		}

		err = collectionsInitializeMintParam(ctx, keepers.MintKeeper)
		if err != nil {
			ctx.Logger().Error("Error initializing mint params:", "message", err.Error())
		}

		return vm, nil
	}
}

func initializeConsensusParamVersion(ctx sdk.Context, consensusKeeper consensusparamkeeper.Keeper) error {
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

func unbondNow(ctx sdk.Context, keepers *keepers.AppKeepers, ubfs []UnbondingFail) error {
	for i := 0; i < len(ubfs); i++ {
		u := ubfs[i]
		ubd, err := keepers.StakingKeeper.GetUnbondingDelegationByUnbondingID(ctx, u.UnbondingId)
		if err != nil {
			continue
		}

		// ensure accuracy
		if u.Delegator == ubd.DelegatorAddress && u.Validator == ubd.ValidatorAddress && ubd.Entries[u.Index].Balance.Equal(u.Balance) {
			ubd.Entries[u.Index].UnbondingId = u.UnbondingId
			ubd.Entries[u.Index].CompletionTime = ctx.BlockHeader().Time.Add(-1 * time.Hour)
			ubd.Entries[u.Index].UnbondingOnHoldRefCount = 1

			err := keepers.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
			if err != nil {
				return err

			}
			err = keepers.StakingKeeper.UnbondingCanComplete(ctx, u.UnbondingId)
			if err != nil {
				return err
			}
			// index ubd update
			indexUpdate(u.Delegator, u.Validator, u.Index, ubfs)
		}
	}

	return nil
}

// getInfoUnbondingFail Get info UnbondingFail
func getInfoUnbondingFail(ctx sdk.Context, keepers *keepers.AppKeepers) (ubf []UnbondingFail, err error) {
	// ubdID errors: 3500-3700 (3599-3678)
	for i := 3500; i < 3700; i++ {
		id := uint64(i)
		ubd, err := keepers.StakingKeeper.GetUnbondingDelegationByUnbondingID(ctx, id)
		if err != nil {
			continue // not found ubd for id
		}

		for idx, entry := range ubd.Entries {
			if entry.UnbondingOnHoldRefCount != 0 && entry.UnbondingId == id {
				ubf = append(ubf, UnbondingFail{
					Delegator:   ubd.DelegatorAddress,
					Validator:   ubd.ValidatorAddress,
					Balance:     entry.Balance,
					UnbondingId: id,
					Index:       idx,
				})
				break
			}

		}
	}

	return ubf, nil
}

func indexUpdate(del, val string, index int, ubf []UnbondingFail) {
	for i := 0; i < len(ubf); i++ {
		if ubf[i].Delegator == del && ubf[i].Validator == val {
			if index < ubf[i].Index {
				ubf[i].Index = ubf[i].Index - 1
			}
		}
	}
}
