package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// ReBalanceDelegation re-balances the DAO staking among validators bases on the current validators self bond.
func (k Keeper) ReBalanceDelegation(ctx sdk.Context) error {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	vals := k.stakingKeeper.GetAllValidators(ctx)
	targetDaoStaking := k.getTargetDelegationState(ctx, vals)
	delegations, undelegations := k.computeDelegationsAndUndelegation(ctx, daoAddr, vals, targetDaoStaking)

	if len(delegations) == 0 && len(undelegations) == 0 {
		return nil
	}

	// If we have updates in the (un)delegations we should withdraw the rewards and recompute the (un)delegations
	//	with the received reward. Otherwise, it will be withdrawn during the (un)delegations execution (via staking hook),
	//	and that will cause the re-balancing for each following block.
	if err := k.WithdrawReward(ctx); err != nil {
		return err
	}

	vals = k.stakingKeeper.GetAllValidators(ctx)
	targetDaoStaking = k.getTargetDelegationState(ctx, vals)
	delegations, undelegations = k.computeDelegationsAndUndelegation(ctx, daoAddr, vals, targetDaoStaking)

	if err := undelegateValidators(ctx, vals, undelegations, k, daoAddr); err != nil {
		return err
	}

	return k.delegateValidators(ctx, vals, delegations, daoAddr)
}

// GetDaoDelegationSupply returns total amount of the treasury bonded coins.
func (k Keeper) GetDaoDelegationSupply(ctx context.Context) math.LegacyDec {
	return k.getDaoDelegationSupply(ctx, k.stakingKeeper.GetAllValidators(ctx))
}

// getTargetDelegationState builds a map of the validators and the stake amount they should have now.
// if the validator is not in the map, the DAO stake is zero.
func (k Keeper) getTargetDelegationState(ctx sdk.Context, vals []stakingtypes.Validator) map[string]math.Int {
	maxValCommission := k.MaxValCommission(ctx)
	valsSelfBonds := make(map[string]math.LegacyDec) // the key is OperatorAddress
	valsSelfBondsSupply := math.LegacyZeroDec()
	for _, val := range vals {
		if !val.IsBonded() {
			continue
		}
		if val.GetCommission().GT(maxValCommission) {
			continue
		}

		valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
		if err != nil {
			panic(err)
		}
		accVal := sdk.AccAddress(valAddr)
		selfDelegation, found := k.stakingKeeper.GetDelegation(ctx, accVal, valAddr)
		if !found || selfDelegation.GetShares().IsZero() {
			continue
		}
		selfDelegationAmount := val.TokensFromShares(selfDelegation.GetShares())
		valsSelfBonds[valAddr.String()] = selfDelegationAmount
		valsSelfBondsSupply = valsSelfBondsSupply.Add(selfDelegationAmount)
	}

	daoDelegationSupply := k.getDaoDelegationSupply(ctx, vals)
	daoBondDenomSupply := k.treasuryBondDenomAmount(ctx).ToLegacyDec().Add(daoDelegationSupply)

	daoBondDenomToDelegate := daoBondDenomSupply.Sub(daoBondDenomSupply.Mul(k.PoolRate(ctx)))

	targetDelegationState := make(map[string]math.Int) // the key is OperatorAddress
	for valAddr, selfDelegationAmt := range valsSelfBonds {
		valDaoDelegationAmt := selfDelegationAmt.Mul(daoBondDenomToDelegate).Quo(valsSelfBondsSupply).TruncateInt()
		if !valDaoDelegationAmt.IsZero() {
			targetDelegationState[valAddr] = valDaoDelegationAmt
		}
	}

	return targetDelegationState
}

// getDaoDelegationSupply returns total amount of the treasury bonded coins.
func (k Keeper) getDaoDelegationSupply(ctx context.Context, vals []stakingtypes.Validator) math.LegacyDec {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	totalStakingSupply := math.LegacyZeroDec()
	for _, val := range vals {
		valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
		if err != nil {
			panic(err)
		}
		delegation, found := k.stakingKeeper.GetDelegation(ctx, daoAddr, valAddr)
		if !found || delegation.GetShares().IsZero() {
			continue
		}
		delegationAmount := val.TokensFromShares(delegation.GetShares())
		totalStakingSupply = totalStakingSupply.Add(delegationAmount)
	}

	return totalStakingSupply
}

// computeDelegationsAndUndelegation computes the target (un)delegations.
func (k Keeper) computeDelegationsAndUndelegation(
	ctx sdk.Context,
	daoAddr sdk.AccAddress,
	vals []stakingtypes.Validator,
	targetDaoStaking map[string]math.Int,
) (map[string]math.Int, map[string]math.Int) {
	delegations := make(map[string]math.Int)
	undelegations := make(map[string]math.Int)
	for _, val := range vals {
		valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
		if err != nil {
			panic(err)
		}
		targetDaoDelegation, ok := targetDaoStaking[valAddr.String()]
		delegatedByDao := math.ZeroInt()
		delegation, found := k.stakingKeeper.GetDelegation(ctx, daoAddr, valAddr)
		if found {
			delegatedByDao = val.TokensFromShares(delegation.GetShares()).TruncateInt()
		}
		// for the validators not in the target list the target amount is zero
		if !ok {
			targetDaoDelegation = math.ZeroInt()
		}

		delegationDelta := targetDaoDelegation.Sub(delegatedByDao)
		if delegationDelta.IsZero() {
			continue
		}

		if delegationDelta.IsNegative() {
			undelegations[valAddr.String()] = delegatedByDao.Sub(targetDaoDelegation)
			continue
		}

		delegations[valAddr.String()] = delegationDelta
	}
	return delegations, undelegations
}

// undelegateValidators undelegates the requested amount from the validators in the undelegations.
func undelegateValidators(ctx sdk.Context, vals []stakingtypes.Validator, undelegations map[string]math.Int, k Keeper, daoAddr sdk.AccAddress) error {
	for _, val := range vals {
		valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
		if err != nil {
			return err
		}
		undelegationAmt, ok := undelegations[valAddr.String()]
		if !ok {
			continue
		}
		if undelegationAmt.IsZero() {
			continue
		}
		undelegationShares, err := val.SharesFromTokens(undelegationAmt)
		if err != nil {
			return err
		}
		if undelegationShares.IsZero() {
			continue
		}

		if _, err := k.stakingKeeper.UnbondAndUndelegateCoins(ctx, daoAddr, valAddr, undelegationShares); err != nil {
			return err
		}
	}
	if len(undelegations) > 0 {
		k.Logger(ctx).Info(fmt.Sprintf("rebalanced, undelegated: %+v", undelegations))
	}

	return nil
}

// UndelegateAllValidators undelegate all validators.
func (k Keeper) UndelegateAllValidators(ctx context.Context) (err error) {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	vals := k.stakingKeeper.GetAllValidators(ctx)

	for _, val := range vals {
		valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
		if err != nil {
			return err
		}

		delegation, found := k.stakingKeeper.GetDelegation(ctx, daoAddr, valAddr)
		if !found || delegation.GetShares().IsZero() {
			continue
		}

		if _, err = k.stakingKeeper.UnbondAndUndelegateCoins(ctx, daoAddr, valAddr, delegation.GetShares()); err != nil {
			return err
		}
	}

	return err
}

// delegateValidators delegates the requested amount from the validators in the delegations.
func (k Keeper) delegateValidators(ctx sdk.Context, vals []stakingtypes.Validator, delegations map[string]math.Int, daoAddr sdk.AccAddress) error {
	for _, val := range vals {
		delegation, ok := delegations[val.GetOperator()]
		if !ok {
			continue
		}
		if delegation.IsZero() {
			continue
		}
		if _, err := k.stakingKeeper.Delegate(ctx, daoAddr, delegation, stakingtypes.Unbonded, val, true); err != nil {
			return err
		}
	}
	if len(delegations) > 0 {
		k.Logger(ctx).Info(fmt.Sprintf("rebalanced, delegated: %+v", delegations))
	}

	return nil
}
