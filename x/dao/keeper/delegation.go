package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// ReBalanceDelegation re-balances the DAO staking among validators bases on the current validators self bond.
func (k Keeper) ReBalanceDelegation(ctx sdk.Context) error {
	vals := k.stakingKeeper.GetAllValidators(ctx)
	targetDelegations := k.getTargetDelegationState(ctx, vals)
	return k.reBalanceDaoStaking(ctx, vals, targetDelegations)
}

// getTargetDelegationState builds a map of the validators and the stake amount they should have now.
// if the validator is not in the map, the DAO stake is zero.
func (k Keeper) getTargetDelegationState(ctx sdk.Context, vals []stakingtypes.Validator) map[string]sdk.Dec {
	maxStakingCommissionRate := k.StakingMaxCommissionRate(ctx)
	valsSelfBonds := make(map[string]sdk.Dec) // the key is OperatorAddress
	valsSelfBondsSupply := sdk.NewDec(0)
	for _, val := range vals {
		if !val.IsBonded() {
			continue
		}
		if val.GetCommission().GT(maxStakingCommissionRate) {
			continue
		}

		valOperator := val.GetOperator()
		valAddr := sdk.AccAddress(valOperator)
		selfDelegation := k.stakingKeeper.Delegation(ctx, valAddr, valOperator)
		if selfDelegation != nil && !selfDelegation.GetShares().IsZero() {
			valsSelfBonds[valOperator.String()] = selfDelegation.GetShares()
			valsSelfBondsSupply = valsSelfBondsSupply.Add(selfDelegation.GetShares())
		}
	}

	daoDelegationSupply := k.getDaoDelegationSupply(ctx)
	daoBondDenomSupply := k.treasuryBondDenomAmount(ctx).ToDec().Add(daoDelegationSupply)

	stakingTokenPoolRate := k.StakingTokenPoolRate(ctx)
	daoBondDenomToDelegate := daoBondDenomSupply.Mul(sdk.OneDec().Sub(stakingTokenPoolRate))

	targetDelegationState := make(map[string]sdk.Dec) // the key is OperatorAddress
	for valAddr, selfDelegationAmount := range valsSelfBonds {
		valDaoDelegationAmount := selfDelegationAmount.Quo(valsSelfBondsSupply).Mul(daoBondDenomToDelegate)
		if !valDaoDelegationAmount.IsZero() {
			targetDelegationState[valAddr] = valDaoDelegationAmount
		}
	}

	return targetDelegationState
}

func (k Keeper) getDaoDelegationSupply(ctx sdk.Context) sdk.Dec {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	delegations := k.stakingKeeper.GetAllDelegatorDelegations(ctx, daoAddr)
	totalStakingSupply := sdk.NewDec(0)
	for _, delegation := range delegations {
		totalStakingSupply = totalStakingSupply.Add(delegation.GetShares())
	}

	return totalStakingSupply
}

func (k Keeper) reBalanceDaoStaking(ctx sdk.Context, vals []stakingtypes.Validator, targetDaoStaking map[string]sdk.Dec) error {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	delegations := make(map[string]sdk.Int)
	undelegations := make(map[string]sdk.Dec)
	for _, val := range vals {
		valAddr := val.GetOperator()
		targetDaoDelegation, ok := targetDaoStaking[valAddr.String()]
		delegation := k.stakingKeeper.Delegation(ctx, daoAddr, valAddr)
		delegatedByDao := sdk.NewDec(0)
		if delegation != nil {
			delegatedByDao = delegation.GetShares()
		}
		// for the validators not in the target list the target amount is zero
		if !ok {
			targetDaoDelegation = sdk.NewDec(0)
		}

		delegationDelta := targetDaoDelegation.RoundInt().Sub(delegatedByDao.RoundInt())
		if delegationDelta.IsZero() {
			continue
		}

		if delegationDelta.IsNegative() {
			undelegations[valAddr.String()] = delegatedByDao.RoundInt().Sub(targetDaoDelegation.RoundInt()).ToDec()
			continue
		}

		delegations[valAddr.String()] = delegationDelta
	}

	err := undelegateValidators(ctx, vals, undelegations, k, daoAddr)
	if err != nil {
		return err
	}

	err2 := k.delegateValidators(ctx, vals, delegations, daoAddr)
	if err2 != nil {
		return err2
	}

	return nil
}

func undelegateValidators(ctx sdk.Context, vals []stakingtypes.Validator, undelegations map[string]sdk.Dec, k Keeper, daoAddr sdk.AccAddress) error {
	for _, val := range vals {
		undelegation, ok := undelegations[val.GetOperator().String()]
		if !ok {
			continue
		}
		if _, err := k.stakingKeeper.Undelegate(ctx, daoAddr, val.GetOperator(), undelegation); err != nil {
			return err
		}
		if err := k.completeUnbonding(ctx, daoAddr, val.GetOperator()); err != nil {
			return err
		}
	}
	return nil
}

// completeUnbinding completes the unbonding of all entries in the
// retrieved unbonding delegation object or an error upon failure.
func (k Keeper) completeUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	ubd, found := k.stakingKeeper.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return stakingtypes.ErrNoUnbondingDelegation
	}

	bondDenom := k.stakingKeeper.BondDenom(ctx)
	delegatorAddress, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		return err
	}

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if !entry.Balance.IsZero() {
			amt := sdk.NewCoin(bondDenom, entry.Balance)
			if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
				ctx, stakingtypes.NotBondedPoolName, delegatorAddress, sdk.NewCoins(amt),
			); err != nil {
				return err
			}
		}
	}

	k.stakingKeeper.RemoveUnbondingDelegation(ctx, ubd)
	return nil
}

func (k Keeper) delegateValidators(ctx sdk.Context, vals []stakingtypes.Validator, delegations map[string]sdk.Int, daoAddr sdk.AccAddress) error {
	for _, val := range vals {
		delegation, ok := delegations[val.GetOperator().String()]
		if !ok {
			continue
		}
		if _, err := k.stakingKeeper.Delegate(ctx, daoAddr, delegation, stakingtypes.Unbonded, val, true); err != nil {
			return err
		}
	}
	return nil
}
