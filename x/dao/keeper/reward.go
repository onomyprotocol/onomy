package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// WithdrawReward withdraw dao delegation reward.
func (k Keeper) WithdrawReward(ctx sdk.Context) error {
	vals := k.stakingKeeper.GetAllValidators(ctx)
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	for _, val := range vals {
		valOperator := val.GetOperator()
		_, found := k.stakingKeeper.GetDelegation(ctx, daoAddr, valOperator)
		if !found {
			continue
		}
		// check existence of delegator starting info
		if !k.distributionKeeper.HasDelegatorStartingInfo(ctx, valOperator, daoAddr) {
			continue
		}
		if _, err := k.distributionKeeper.WithdrawDelegationRewards(ctx, daoAddr, valOperator); err != nil {
			return err
		}
	}
	return nil
}
