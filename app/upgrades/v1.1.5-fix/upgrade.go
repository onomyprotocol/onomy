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
		toBeRemovedUbdIDs := map[uint64]bool{}
		var consumerChainIDS []string
		for _, chain := range pk.GetAllConsumerChains(ctx) {
			consumerChainIDS = append(consumerChainIDS, chain.ChainId)
		}

		for _, chainID := range consumerChainIDS {
			for _, ubdOpIndex := range pk.GetAllUnbondingOpIndexes(ctx, chainID) {
				for _, id := range ubdOpIndex.UnbondingOpIds {
					if _, found := sk.GetUnbondingType(ctx, id); found {
						toBeRemovedUbdIDs[id] = true
					}
				}
			}
		}

		for id := range toBeRemovedUbdIDs {

			if len(consumerChainIDS) == 0 {
				break
			}
			valsetUpdateID := pk.GetValidatorSetUpdateId(ctx)

			// Add to indexes
			for _, consumerChainID := range consumerChainIDS {
				ubdIds, ok := pk.GetUnbondingOpIndex(ctx, consumerChainID, valsetUpdateID)
				if !ok {
					continue
				}

				newIds := []uint64{}

				for _, ubdId := range ubdIds {
					if ubdId == id {
						continue
					}

					newIds = append(newIds, ubdId)
				}

				// filter out invalid ID
				pk.SetUnbondingOpIndex(ctx, consumerChainID, valsetUpdateID, newIds)
			}

			// remove ubd entries
			_, found := pk.GetUnbondingOp(ctx, id)
			if found {
				pk.DeleteUnbondingOp(ctx, id)
			}
		}

		// clear invalid mature ubd entries
		ids := []uint64{}
		for _, id := range pk.GetMaturedUnbondingOps(ctx) {
			if toBeRemovedUbdIDs[id] {
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