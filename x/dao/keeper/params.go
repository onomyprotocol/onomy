package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// GetParams returns the total set of dao parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return k.getParams(ctx)
}

// SetParams set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.ps.SetParamSet(ctx, &params)
}

// WithdrawRewardPeriod - the blocks period to dao staking reward.
func (k Keeper) WithdrawRewardPeriod(ctx sdk.Context) (res int64) {
	k.ps.Get(ctx, types.KeyWithdrawRewardPeriod, &res)
	return
}

// StakingMaxCommissionRate - the max validator's commission to be staked by the dao.
func (k Keeper) StakingMaxCommissionRate(ctx sdk.Context) (res sdk.Dec) {
	k.ps.Get(ctx, types.KeyStakingMaxCommissionRate, &res)
	return
}

// StakingTokenMaxProposalRate - the rate of max allowed proposal from the treasury amount.
func (k Keeper) StakingTokenMaxProposalRate(ctx sdk.Context) (res sdk.Dec) {
	k.ps.Get(ctx, types.KeyStakingTokenMaxProposalRate, &res)
	return
}

// StakingTokenPoolRate - the rate of total dao's staking coins to keep unstaked.
func (k Keeper) StakingTokenPoolRate(ctx sdk.Context) (res sdk.Dec) {
	k.ps.Get(ctx, types.KeyStakingTokenPoolRate, &res)
	return
}

func (k Keeper) getParams(ctx sdk.Context) (params types.Params) {
	k.ps.GetParamSet(ctx, &params)
	return params
}
