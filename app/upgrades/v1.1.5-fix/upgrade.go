// Package v1_1_5 is contains chain upgrade of the corresponding version.
package v1_1_5_fix //nolint:revive,stylecheck // app version

import (
	"fmt"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibcproviderkeeper "github.com/cosmos/interchain-security/x/ccv/provider/keeper"
	ibcprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	ccv "github.com/cosmos/interchain-security/x/ccv/types"
	"github.com/onomyprotocol/onomy/app/upgrades"
)

var targetIds = []uint64{uint64(3372), uint64(3374), uint64(3373)}

func CreateFork(
	sk *stakingkeeper.Keeper,
	pk *ibcproviderkeeper.Keeper,
	providerStoreKey sdk.StoreKey,
) upgrades.Fork {
	forkLogic := func(ctx sdk.Context) {
		for _, id := range targetIds {
			var consumerChainIDS []string

			for _, chain := range pk.GetAllConsumerChains(ctx) {
				consumerChainIDS = append(consumerChainIDS, chain.ChainId)
			}

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
				fmt.Println("old", ubdIds)

				newIds := []uint64{}

				for _, ubdId := range ubdIds {
					fmt.Println("should", ubdId, id)
					if ubdId == id {
						ctx.Logger().Info(fmt.Sprintf("filter out id %d", ubdId))
						continue
					}

					newIds = append(newIds, ubdId)
				}

				fmt.Println("new", newIds)

				// filter out invalid ID
				pk.SetUnbondingOpIndex(ctx, consumerChainID, valsetUpdateID, newIds)
			}

			// remove ubd entries
			_, found := pk.GetUnbondingOp(ctx, id)
			if found {
				pk.DeleteUnbondingOp(ctx, id)
				ctx.Logger().Info(fmt.Sprintf("delete id %d", id))
			}
		}

		// clear invalid mature ubd entries
		ids := []uint64{}
		for _, id := range pk.GetMaturedUnbondingOps(ctx) {
			fmt.Println("should true", slices.Contains(targetIds, id))

			if slices.Contains(targetIds, id) {
				ctx.Logger().Info(fmt.Sprintf("filter matured id %d", id))
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
		ctx.Logger().Info(fmt.Sprint("updated matured ids"))
	}
	return upgrades.Fork{
		UpgradeName:    Name,
		UpgradeHeight:  UpgradeHeight,
		BeginForkLogic: forkLogic,
	}
}
