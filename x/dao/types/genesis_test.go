package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "positive_default",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "positive_custom_params",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod:        1,
					StakingTokenPoolRate:        sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingTokenMaxProposalRate: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingMaxCommissionRate:    sdk.NewDec(1).Quo(sdk.NewDec(10)),
				},
			},
			valid: true,
		},
		{
			desc:     "negative_empty",
			genState: &types.GenesisState{},
			valid:    false,
		},
		{
			desc: "negative_negative_withdraw_reward_period",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod:        -1,
					StakingTokenPoolRate:        sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingTokenMaxProposalRate: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingMaxCommissionRate:    sdk.NewDec(1).Quo(sdk.NewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_staking_token_pool_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod:        1,
					StakingTokenPoolRate:        sdk.NewDec(-1).Quo(sdk.NewDec(10)),
					StakingTokenMaxProposalRate: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingMaxCommissionRate:    sdk.NewDec(1).Quo(sdk.NewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_more_than_one_staking_token_pool_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod:        1,
					StakingTokenPoolRate:        sdk.NewDec(11).Quo(sdk.NewDec(10)),
					StakingTokenMaxProposalRate: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingMaxCommissionRate:    sdk.NewDec(1).Quo(sdk.NewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_staking_token_max_proposal_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod:        1,
					StakingTokenPoolRate:        sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingTokenMaxProposalRate: sdk.NewDec(-1).Quo(sdk.NewDec(10)),
					StakingMaxCommissionRate:    sdk.NewDec(1).Quo(sdk.NewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_more_than_one_staking_token_max_proposal_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod:        1,
					StakingTokenPoolRate:        sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingTokenMaxProposalRate: sdk.NewDec(11).Quo(sdk.NewDec(10)),
					StakingMaxCommissionRate:    sdk.NewDec(1).Quo(sdk.NewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_staking_max_commission_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod:        1,
					StakingTokenPoolRate:        sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingTokenMaxProposalRate: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingMaxCommissionRate:    sdk.NewDec(-1).Quo(sdk.NewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_more_than_one_staking_max_commission_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod:        1,
					StakingTokenPoolRate:        sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingTokenMaxProposalRate: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					StakingMaxCommissionRate:    sdk.NewDec(11).Quo(sdk.NewDec(10)),
				},
			},
			valid: false,
		},
	} {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
