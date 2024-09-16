package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// WithdrawReward withdraw dao delegation reward.
func (k Keeper) WithdrawReward(ctx context.Context) error {
	vals := k.stakingKeeper.GetAllValidators(ctx)
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	for _, val := range vals {
		valOperator := val.GetOperator()
		valAddr, err := sdk.ValAddressFromBech32(valOperator)
		if err != nil {
			return err
		}
		_, found := k.stakingKeeper.GetDelegation(ctx, daoAddr, valAddr)
		if !found {
			continue
		}
		// check existence of delegator starting info
		if !k.distributionKeeper.HasDelegatorStartingInfo(ctx, valAddr, daoAddr) {
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
