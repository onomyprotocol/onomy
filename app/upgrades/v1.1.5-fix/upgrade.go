// Package v1_1_5 is contains chain upgrade of the corresponding version.
package v1_1_5_fix //nolint:revive,stylecheck // app version

import (
	"fmt"
	"slices"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibcproviderkeeper "github.com/cosmos/interchain-security/v5/x/ccv/provider/keeper"
	ibcprovidertypes "github.com/cosmos/interchain-security/v5/x/ccv/provider/types"

	"github.com/onomyprotocol/onomy/app/upgrades"
)

var targetIds = []uint64{uint64(3372), uint64(3374), uint64(3373)}

func CreateFork(
	sk *stakingkeeper.Keeper,
	pk *ibcproviderkeeper.Keeper,
	providerStoreKey storetypes.StoreKey,
) upgrades.Fork {
	forkLogic := func(ctx sdk.Context) {
		for _, id := range targetIds {
			var consumerChainIDS []string

			consumerChainIDS = append(consumerChainIDS, pk.GetAllRegisteredConsumerChainIDs(ctx)...)

			if len(consumerChainIDS) == 0 {
				break
			}
			// Add to indexes.
			for _, consumerChainID := range consumerChainIDS {
				ubdOpIds := pk.GetAllUnbondingOpIndexes(ctx, consumerChainID)
				for _, ubdIds := range ubdOpIds {
					newIds := []uint64{}

					for _, ubdId := range ubdIds.UnbondingOpIds {
						if ubdId == id {
							continue
						}

						newIds = append(newIds, ubdId)
					}

					// filter out invalid ID.
					pk.SetUnbondingOpIndex(ctx, consumerChainID, ubdIds.VscId, newIds)
				}
			}

			// remove ubd entries.
			_, found := pk.GetUnbondingOp(ctx, id)
			if found {
				pk.DeleteUnbondingOp(ctx, id)
			}
		}

		// clear invalid mature ubd entries.
		ids := []uint64{}
		for _, id := range pk.GetMaturedUnbondingOps(ctx) {
			if slices.Contains(targetIds, id) {
				continue
			}

			ids = append(ids, id)
		}

		maturedOps := ibcprovidertypes.MaturedUnbondingOps{
			Ids: ids,
		}
		bz, err := maturedOps.Marshal()
		if err != nil {
			// An error here would indicate something is very wrong,
			// maturedOps is instantiated in this method and should be able to be marshaled.
			panic(fmt.Sprintf("failed to marshal matured unbonding operations: %s", err))
		}

		// update mature ubd ids.
		store := ctx.KVStore(providerStoreKey)
		store.Set(ibcprovidertypes.MaturedUnbondingOpsKey(), bz)
	}
	return upgrades.Fork{
		UpgradeName:    Name,
		UpgradeHeight:  UpgradeHeight,
		BeginForkLogic: forkLogic,
	}
}
