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
		daoDelegation := k.stakingKeeper.Delegation(ctx, daoAddr, valOperator)
		if daoDelegation == nil || daoDelegation.GetShares().IsZero() {
			continue
		}

		if _, err := k.distributionKeeper.WithdrawDelegationRewards(ctx, daoAddr, valOperator); err != nil {
			return err
		}
	}
	return nil
}
