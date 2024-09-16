package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// WithdrawReward withdraw dao delegation reward.
func (k Keeper) WithdrawReward(ctx context.Context) error {
	vals, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	for _, val := range vals {
		valOperator := val.GetOperator()
		valAddr, err := sdk.ValAddressFromBech32(valOperator)
		if err != nil {
			return err
		}
		_, err = k.stakingKeeper.GetDelegation(ctx, daoAddr, valAddr)
		if err != nil {
			continue
		}
		// check existence of delegator starting info
		if has, err := k.distributionKeeper.HasDelegatorStartingInfo(ctx, valAddr, daoAddr); err != nil || !has {
			continue
		}

		reward, err := k.distributionKeeper.WithdrawDelegationRewards(ctx, daoAddr, valAddr)
		if err != nil {
			return err
		}
		k.Logger(ctx).Info(fmt.Sprintf("withdrawn reward: %s", reward.String()))
	}
	return nil
}
