// Package v1_1_5 is contains chain upgrade of the corresponding version.
package v1_1_5_fix //nolint:revive,stylecheck // app version

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcproviderkeeper "github.com/cosmos/interchain-security/x/ccv/provider/keeper"
	ibcprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	ccv "github.com/cosmos/interchain-security/x/ccv/types"
)

// Name is migration name.
const Name = "v1.1.5-fix"

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak *authkeeper.AccountKeeper,
	sk *stakingkeeper.Keeper,
	pk *ibcproviderkeeper.Keeper,
	providerStoreKey sdk.StoreKey,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ids := []uint64{}
		for _, id := range pk.GetMaturedUnbondingOps(ctx) {
			// Attempt to complete unbonding in staking module
			_, found := sk.GetUnbondingType(ctx, id)
			if !found {
				continue
			}
			ids = append(ids, id)
		}

		maturedOps := ccv.MaturedUnbondingOps{
			Ids: ids,
		}
		bz, err := maturedOps.Marshal()
		if err != nil {
			// An error here would indicate something is very wrong,
			// maturedOps is instantiated in this method and should be able to be marshaled.
			panic(fmt.Sprintf("failed to marshal matured unbonding operations: %s", err))
		}

		// update mature ubd ids
		store := ctx.KVStore(providerStoreKey)
		store.Set(ibcprovidertypes.MaturedUnbondingOpsKey(), bz)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
