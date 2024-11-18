// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v2_1_0 //nolint:revive,stylecheck // app version

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/onomyprotocol/onomy/app/keepers"
)

// Name is migration name. Migrate denom "anom" to "ono"
const Name = "v2.2.0"

// UpgradeHandler is an x/upgrade handler.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {

	return func(c context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		// staking
		paramStaking, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, err
		}
		paramStaking.BondDenom = "ono"
		err = keepers.StakingKeeper.SetParams(ctx, paramStaking)
		if err != nil {
			return vm, err
		}

		// distribution

		// bank
		keepers.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
			Base: "ono",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "ono",
					Exponent: 18,
				},
			},
			Display:     "ono",
			Name:        "ONO",
			Symbol:      "ONO",
			Description: `Migrate denom "anom" to "ono"`,
		})

		err = keepers.BankKeeper.Balances.Walk(ctx, nil, func(key collections.Pair[sdk.AccAddress, string], value math.Int) (stop bool, err error) {
			if key.K2() == "anom" {
				err = keepers.BankKeeper.Balances.Set(ctx, collections.Join(sdk.AccAddress(key.K1()), "ono"), value)
				if err != nil {
					return true, err
				}

				err = keepers.BankKeeper.Balances.Remove(ctx, key)
				if err != nil {
					return true, err
				}
			}
			return false, nil
		})
		if err != nil {
			return vm, err
		}

		// gov
		govParams, err := keepers.GovKeeper.Params.Get(ctx)
		if err != nil {
			return vm, err
		}
		for i, coin := range govParams.MinDeposit {
			if coin.Denom == "anom" {
				govParams.MinDeposit[i] = sdk.NewCoin("ono", coin.Amount)
			}
		}
		err = keepers.GovKeeper.Params.Set(ctx, govParams)
		if err != nil {
			return vm, err
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
