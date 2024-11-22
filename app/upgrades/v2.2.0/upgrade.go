// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v2_1_0 //nolint:revive,stylecheck // app version

import (
	"context"
	"slices"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibctypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	"github.com/onomyprotocol/onomy/app/keepers"
)

// Name is migration name. Migrate denom "anom" to "ono".
const Name = "v2.2.0"

var newDenom = "ono"

// UpgradeHandler is an x/upgrade handler.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		// staking.
		paramStaking, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, err
		}
		paramStaking.BondDenom = newDenom
		err = keepers.StakingKeeper.SetParams(ctx, paramStaking)
		if err != nil {
			return vm, err
		}

		// ibc
		var escrowAddress []string
		balancesIBCTranfer := keepers.TransferKeeper.GetAllTotalEscrowed(ctx)
		keepers.IBCKeeper.ChannelKeeper.IterateChannels(ctx, func(ic ibctypes.IdentifiedChannel) bool {
			escrowAddress = append(escrowAddress, transfertypes.GetEscrowAddress(ic.PortId, ic.ChannelId).String())
			return false
		})

		// bank.
		keepers.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
			Base: newDenom,
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    newDenom,
					Exponent: 18,
				},
			},
			Display:     newDenom,
			Name:        "ONO",
			Symbol:      "ONO",
			Description: `Migrate denom "anom" to "ono"`,
		})

		err = keepers.BankKeeper.Balances.Walk(ctx, nil, func(key collections.Pair[sdk.AccAddress, string], value math.Int) (stop bool, err error) {
			addr := sdk.AccAddress(key.K1())
			if slices.Contains(escrowAddress, addr.String()) {
				return false, nil
			}
			if key.K2() == "anom" {
				err = keepers.BankKeeper.Balances.Set(ctx, collections.Join(addr, newDenom), value)
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

		err = keepers.BankKeeper.Supply.Walk(ctx, nil, func(key string, value math.Int) (stop bool, err error) {
			if key == "anom" {
				err = keepers.BankKeeper.Supply.Set(ctx, newDenom, value.Sub(balancesIBCTranfer.AmountOf("anom")))
				if err != nil {
					return true, err
				}
				err = keepers.BankKeeper.Supply.Set(ctx, "anom", balancesIBCTranfer.AmountOf("anom"))
			}
			return false, err
		})
		if err != nil {
			return vm, err
		}

		// gov min_deposit, expedited_min_deposit.

		govParams, err := keepers.GovKeeper.Params.Get(ctx)
		if err != nil {
			return vm, err
		}
		for i, coin := range govParams.MinDeposit {
			if coin.Denom == "anom" {
				govParams.MinDeposit[i] = sdk.NewCoin(newDenom, coin.Amount)
			}
		}
		for i, coin := range govParams.ExpeditedMinDeposit {
			if coin.Denom == "anom" {
				govParams.ExpeditedMinDeposit[i] = sdk.NewCoin(newDenom, coin.Amount)
			}
		}

		err = keepers.GovKeeper.Params.Set(ctx, govParams)
		if err != nil {
			return vm, err
		}
		// mint
		mintParams, err := keepers.MintKeeper.Params.Get(ctx)
		if err != nil {
			return vm, err
		}
		mintParams.MintDenom = newDenom
		err = keepers.MintKeeper.Params.Set(ctx, mintParams)
		if err != nil {
			return vm, err
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
