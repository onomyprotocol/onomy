package types_test

import (
	"testing"

	"cosmossdk.io/math"
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
					WithdrawRewardPeriod: 1,
					PoolRate:             math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxProposalRate:      math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxValCommission:     math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
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
					WithdrawRewardPeriod: -1,
					PoolRate:             math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxProposalRate:      math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxValCommission:     math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_pool_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod: 1,
					PoolRate:             math.LegacyNewDec(-1).Quo(math.LegacyNewDec(10)),
					MaxProposalRate:      math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxValCommission:     math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_more_than_one_pool_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod: 1,
					PoolRate:             math.LegacyNewDec(11).Quo(math.LegacyNewDec(10)),
					MaxProposalRate:      math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxValCommission:     math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_max_proposal_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod: 1,
					PoolRate:             math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxProposalRate:      math.LegacyNewDec(-1).Quo(math.LegacyNewDec(10)),
					MaxValCommission:     math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_more_than_one_max_proposal_rate",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod: 1,
					PoolRate:             math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxProposalRate:      math.LegacyNewDec(11).Quo(math.LegacyNewDec(10)),
					MaxValCommission:     math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_max_val_commission",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod: 1,
					PoolRate:             math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxProposalRate:      math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxValCommission:     math.LegacyNewDec(-1).Quo(math.LegacyNewDec(10)),
				},
			},
			valid: false,
		},
		{
			desc: "negative_negative_more_than_one_max_val_commission",
			genState: &types.GenesisState{
				Params: types.Params{
					WithdrawRewardPeriod: 1,
					PoolRate:             math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxProposalRate:      math.LegacyNewDec(1).Quo(math.LegacyNewDec(10)),
					MaxValCommission:     math.LegacyNewDec(11).Quo(math.LegacyNewDec(10)),
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
