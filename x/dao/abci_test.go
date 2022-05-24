package dao_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

var (
	fiftyPercents                     = sdk.NewDec(1).QuoInt64(2)                                                                       //nolint:gochecknoglobals
	tenPercents                       = sdk.NewDec(1).Quo(sdk.NewDec(10))                                                               //nolint:gochecknoglobals
	nanoBondCoins                     = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000000)                                              // not enough for validator to be bonded
	twoBondCoins                      = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(2, sdk.DefaultPowerReduction))   //nolint:gochecknoglobals
	tenBondCoins                      = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))  //nolint:gochecknoglobals
	hundredBondCoins                  = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)) //nolint:gochecknoglobals
	lowCommission                     = stakingtypes.NewCommissionRates(tenPercents, tenPercents, tenPercents)                          //nolint:gochecknoglobals
	highCommission                    = stakingtypes.NewCommissionRates(fiftyPercents, fiftyPercents, fiftyPercents)                    //nolint:gochecknoglobals
	hundredBondWithoutStakingPoolRate = hundredBondCoins.Amount.ToDec().Mul(sdk.OneDec().Sub(types.DefaultStakingTokenPoolRate))        //nolint:gochecknoglobals
)

func TestEndBlocker_ReBalance(t *testing.T) {
	type args struct {
		vals            map[string]simapp.ValReq
		treasuryBalance sdk.Coin
	}

	type wantAssertion struct {
		vals            map[string]valAssertion
		treasuryBalance sdk.Coin
	}

	tests := []struct {
		name string
		args args
		want wantAssertion
	}{
		{
			name: "positive",
			args: args{
				vals: map[string]simapp.ValReq{
					"val1": { // bonded
						SelfBondCoin: twoBondCoins,
						Commission:   lowCommission,
					},
					"val2": { // bonded
						SelfBondCoin: tenBondCoins,
						Commission:   lowCommission,
					},
					"val3": { // won't be bonded
						SelfBondCoin: nanoBondCoins,
						Commission:   lowCommission,
					},
					"val4": { // bonded, but high Commission to be staked
						SelfBondCoin: tenBondCoins,
						Commission:   highCommission,
					},
				},
				treasuryBalance: hundredBondCoins,
			},
			want: wantAssertion{
				vals: map[string]valAssertion{
					"val1": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: twoBondCoins.Amount.ToDec(),
						// full * self bond / total bond
						daoBondAmount: twoBondCoins.Amount.ToDec().
							Mul(hundredBondWithoutStakingPoolRate).
							Quo(twoBondCoins.Amount.Add(tenBondCoins.Amount).ToDec()).TruncateDec(),
					},
					"val2": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						daoBondAmount: tenBondCoins.Amount.ToDec().
							Mul(hundredBondWithoutStakingPoolRate).
							Quo(twoBondCoins.Amount.Add(tenBondCoins.Amount).ToDec()).TruncateDec(),
					},
					"val3": {
						bondStatus:     stakingtypes.Unbonded,
						selfBondAmount: nanoBondCoins.Amount.ToDec(),
						daoBondAmount:  sdk.ZeroDec(),
					},
					"val4": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						daoBondAmount:  sdk.ZeroDec(),
					},
				},
				treasuryBalance: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(5, sdk.DefaultPowerReduction).AddRaw(1)),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			simApp, _ := createSimAppWithValidatorsAndTreasury(t, tt.args.vals, tt.args.treasuryBalance)

			// iterate couple times to check that the state is the same
			for i := 0; i < 10; i++ {
				simApp.BeginNextBlock()
				ctx := simApp.NewNextContext()
				simApp.EndBlockAndCommit(ctx)
			}

			simApp.BeginNextBlock()
			ctx := simApp.NewNextContext()

			// assertions
			assertValidators(t, simApp, ctx, tt.want.vals)

			daoKeeper := simApp.OnomyApp().DaoKeeper
			gotTreasuryBalance := daoKeeper.Treasury(ctx)
			require.Equal(t, sdk.NewCoins(tt.want.treasuryBalance), gotTreasuryBalance)

			// pool rate = current pool / total
			require.Equal(t, gotTreasuryBalance[0].Amount.ToDec().Quo(daoKeeper.GetDaoDelegationSupply(ctx).Add(gotTreasuryBalance[0].Amount.ToDec())), types.DefaultStakingTokenPoolRate)

			// the check the overall balance remains the same
			require.Equal(t, daoKeeper.GetDaoDelegationSupply(ctx).Add(gotTreasuryBalance[0].Amount.ToDec()), tt.args.treasuryBalance.Amount.ToDec())
		})
	}
}

func TestEndBlocker_WithdrawReward(t *testing.T) {
	validatorReward := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000000)
	expectedDaoFullReward := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1486956434)

	type args struct {
		vals            map[string]simapp.ValReq
		treasuryBalance sdk.Coin
	}

	type wantAssertion struct {
		vals            map[string]valAssertion
		treasuryBalance sdk.Coin
	}

	tests := []struct {
		name string
		args args
		want wantAssertion
	}{
		{
			name: "positive",
			args: args{
				vals: map[string]simapp.ValReq{
					"val1": { // bonded
						SelfBondCoin: tenBondCoins,
						Commission:   lowCommission,
						Reward:       validatorReward,
					},
					"val2": { // bonded
						SelfBondCoin: tenBondCoins,
						Commission:   lowCommission,
						Reward:       validatorReward,
					},
					"val3": { // won't be bonded
						SelfBondCoin: nanoBondCoins,
						Commission:   lowCommission,
						Reward:       validatorReward,
					},
				},
				treasuryBalance: hundredBondCoins,
			},
			want: wantAssertion{
				vals: map[string]valAssertion{
					"val1": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						daoBondAmount:
						// initial dao staking
						hundredBondWithoutStakingPoolRate.QuoInt64(2).
							// the Reward
							Add(expectedDaoFullReward.Amount.ToDec().QuoInt64(2).Mul(sdk.OneDec().Sub(types.DefaultStakingTokenPoolRate))).TruncateDec(),
					},
					"val2": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						daoBondAmount:
						// initial dao staking
						hundredBondWithoutStakingPoolRate.QuoInt64(2).
							// the Reward
							Add(expectedDaoFullReward.Amount.ToDec().QuoInt64(2).Mul(sdk.OneDec().Sub(types.DefaultStakingTokenPoolRate))).TruncateDec(),
					},
					"val3": {
						bondStatus:     stakingtypes.Unbonded,
						selfBondAmount: nanoBondCoins.Amount.ToDec(),
						daoBondAmount:  sdk.ZeroDec(),
					},
				},
				treasuryBalance: sdk.NewCoin(sdk.DefaultBondDenom,
					sdk.TokensFromConsensusPower(5, sdk.DefaultPowerReduction).
						Add(expectedDaoFullReward.Amount.ToDec().Mul(types.DefaultStakingTokenPoolRate).RoundInt())),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const withdrawRewardPeriod = 6 // the simApp.BeginNextBlock() in assertion will be executed with that block number
			simApp, _ := createSimAppWithValidatorsAndTreasury(t, tt.args.vals, tt.args.treasuryBalance)
			simApp.BeginNextBlock()
			ctx := simApp.NewNextContext()
			// update dao params to withdraw Reward
			daoKeeper := simApp.OnomyApp().DaoKeeper
			daoParams := daoKeeper.GetParams(ctx)
			daoParams.WithdrawRewardPeriod = withdrawRewardPeriod
			daoKeeper.SetParams(ctx, daoParams)
			// allocate validator rewards
			for moniker := range tt.args.vals {
				moniker := moniker
				simApp.OnomyApp().StakingKeeper.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
					if moniker == validator.GetMoniker() {
						// mind and send coins as a validator Reward
						require.NoError(t, simApp.OnomyApp().BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(tt.args.vals[moniker].Reward)))
						require.NoError(t, simApp.OnomyApp().BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distrtypes.ModuleName, sdk.NewCoins(tt.args.vals[moniker].Reward)))

						simApp.OnomyApp().DistrKeeper.AllocateTokensToValidator(ctx, validator, sdk.NewDecCoinsFromCoins(tt.args.vals[moniker].Reward))
						return true
					}
					return false
				})
			}
			simApp.EndBlockAndCommit(ctx)

			// assertions
			simApp.BeginNextBlock()
			ctx = simApp.NewNextContext()
			dao.EndBlocker(ctx, simApp.OnomyApp().DaoKeeper)
			assertValidators(t, simApp, ctx, tt.want.vals)

			gotTreasuryBalance := daoKeeper.Treasury(ctx)
			require.Equal(t, sdk.NewCoins(tt.want.treasuryBalance), gotTreasuryBalance)

			// pool rate = current pool / total
			require.Equal(t, gotTreasuryBalance[0].Amount.ToDec().Quo(daoKeeper.GetDaoDelegationSupply(ctx).Add(gotTreasuryBalance[0].Amount.ToDec())), types.DefaultStakingTokenPoolRate)

			// the check the overall balance is increased
			require.Equal(t, daoKeeper.GetDaoDelegationSupply(ctx).Add(gotTreasuryBalance[0].Amount.ToDec()).
				// substitute the Reward from the total dao
				Sub(expectedDaoFullReward.Amount.ToDec()), tt.args.treasuryBalance.Amount.ToDec())
		})
	}
}

func TestEndBlocker_Vote(t *testing.T) {
	type valWithProposalsReq struct {
		simapp.ValReq
		deposit sdk.Coin
	}

	type args struct {
		vals map[string]valWithProposalsReq
	}

	type wantAssertion struct {
		wantDaoProposal map[string]bool // [moniker]should dao vote
	}

	tests := []struct {
		name string
		args args
		want wantAssertion
	}{
		{
			name: "positive_two_active_proposals",
			args: args{
				vals: map[string]valWithProposalsReq{
					"val1": {
						ValReq: simapp.ValReq{
							Balance:      sdk.NewCoins(tenBondCoins.Add(tenBondCoins)),
							SelfBondCoin: tenBondCoins,
							Commission:   lowCommission,
						},
						deposit: tenBondCoins,
					},
					"val2": {
						ValReq: simapp.ValReq{
							Balance:      sdk.NewCoins(tenBondCoins.Add(tenBondCoins)),
							SelfBondCoin: tenBondCoins,
							Commission:   lowCommission,
						},
						deposit: tenBondCoins,
					},
				},
			},
			want: wantAssertion{
				wantDaoProposal: map[string]bool{
					"val1": true,
					"val2": true,
				},
			},
		},
		{
			name: "positive_one_active_proposal",
			args: args{
				vals: map[string]valWithProposalsReq{
					"val1": {
						ValReq: simapp.ValReq{
							Balance:      sdk.NewCoins(tenBondCoins.Add(tenBondCoins)),
							SelfBondCoin: tenBondCoins,
							Commission:   lowCommission,
						},
						deposit: tenBondCoins,
					},
					"val2": {
						ValReq: simapp.ValReq{
							Balance:      sdk.NewCoins(tenBondCoins.Add(tenBondCoins)),
							SelfBondCoin: tenBondCoins,
							Commission:   lowCommission,
						},
						deposit: nanoBondCoins, // low deposit so the dao shouldn't vote
					},
				},
			},
			want: wantAssertion{
				wantDaoProposal: map[string]bool{
					"val1": true,
					"val2": false,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const proposalNamePattern = "proposal-%s"

			vals := make(map[string]simapp.ValReq, len(tt.args.vals))
			for moniker := range tt.args.vals {
				vals[moniker] = tt.args.vals[moniker].ValReq
			}
			simApp, privs := createSimAppWithValidatorsAndTreasury(t, vals, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0))

			// create the text proposals
			for moniker := range tt.args.vals {
				priv := privs[moniker]
				textProposalContent := govtypes.ContentFromProposalType(fmt.Sprintf(proposalNamePattern, moniker), "description", govtypes.ProposalTypeText)
				simApp.CreateProposal(t, textProposalContent, tt.args.vals[moniker].deposit, priv)
			}

			simApp.BeginNextBlock()
			ctx := simApp.NewContext()
			dao.EndBlocker(ctx, simApp.OnomyApp().DaoKeeper)

			// assertions
			govKeeper := simApp.OnomyApp().GovKeeper
			accountKeeper := simApp.OnomyApp().AccountKeeper
			daoAddress := accountKeeper.GetModuleAddress(types.ModuleName)

			votes := govKeeper.GetAllVotes(ctx)
			for moniker, want := range tt.want.wantDaoProposal {
				found := false
				for _, vote := range votes {
					// all votes should from the dao only
					require.Equal(t, daoAddress.String(), vote.Voter)
					proposal, _ := govKeeper.GetProposal(ctx, vote.ProposalId)
					if fmt.Sprintf(proposalNamePattern, moniker) == proposal.GetContent().GetTitle() {
						found = true
						break
					}
				}
				require.Equal(t, want, found)
			}
		})
	}
}

func TestEndBlocker_Slashing_Protection(t *testing.T) {
	// 50% slashing fraction
	fraction := sdk.NewDecWithPrec(5, 1)

	type valWithSlashingReq struct {
		simapp.ValReq
		shouldSlash bool
	}

	type args struct {
		vals            map[string]valWithSlashingReq
		treasuryBalance sdk.Coin
	}

	type wantAssertion struct {
		vals            map[string]valAssertion
		treasuryBalance sdk.Coin
	}

	tests := []struct {
		name string
		args args
		want wantAssertion
	}{
		{
			name: "positive",
			args: args{
				vals: map[string]valWithSlashingReq{
					"val1": {
						ValReq: simapp.ValReq{
							SelfBondCoin: tenBondCoins,
							Commission:   lowCommission,
						},
						shouldSlash: false,
					},
					"val2": { // bonded
						ValReq: simapp.ValReq{
							SelfBondCoin: tenBondCoins,
							Commission:   lowCommission,
						},
						shouldSlash: true,
					},
				},
				treasuryBalance: hundredBondCoins,
			},
			want: wantAssertion{
				vals: map[string]valAssertion{
					"val1": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						// the val2 was slashed so the final amount will higher than val2
						// also the slashing of the validator is based on the voting power, hence the initial
						// amount to slash will be rounded
						// full * self bond / total bond
						daoBondAmount: tenBondCoins.Amount.ToDec().
							// 25*10^16 here is the rounded part
							Mul(hundredBondWithoutStakingPoolRate).
							Quo(tenBondCoins.Amount.ToDec().Mul(fraction).Add(sdk.NewIntWithDecimal(25, 16).ToDec()).Add(tenBondCoins.Amount.ToDec())).TruncateDec(),
					},
					"val2": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec().Mul(fraction).Add(sdk.NewIntWithDecimal(25, 16).ToDec()),
						// the val2 was slashed so the final amount will higher lower val1
						// full * self bond / total bond
						daoBondAmount: tenBondCoins.Amount.ToDec().Mul(fraction).Add(sdk.NewIntWithDecimal(25, 16).ToDec()).
							// 25*10^16 here is the rounded part
							Mul(hundredBondWithoutStakingPoolRate).
							Quo(tenBondCoins.Amount.ToDec().Mul(fraction).Add(sdk.NewIntWithDecimal(25, 16).ToDec()).Add(tenBondCoins.Amount.ToDec())).TruncateDec(),
					},
				},
				treasuryBalance: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(5, sdk.DefaultPowerReduction).AddRaw(1)),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			vals := make(map[string]simapp.ValReq, len(tt.args.vals))
			for moniker := range tt.args.vals {
				vals[moniker] = tt.args.vals[moniker].ValReq
			}
			simApp, _ := createSimAppWithValidatorsAndTreasury(t, vals, tt.args.treasuryBalance)
			// initial rebalance
			simApp.BeginNextBlock()
			ctx := simApp.NewNextContext()
			simApp.EndBlockAndCommit(ctx)

			// slashing
			simApp.BeginNextBlock()
			ctx = simApp.NewNextContext()
			for moniker := range tt.args.vals {
				if !tt.args.vals[moniker].shouldSlash {
					continue
				}
				for _, val := range simApp.OnomyApp().StakingKeeper.GetAllValidators(ctx) {
					if val.GetMoniker() == moniker {
						power := simApp.OnomyApp().StakingKeeper.GetLastValidatorPower(ctx, val.GetOperator())
						consAddr, err := val.GetConsAddr()
						require.NoError(t, err)
						simApp.OnomyApp().StakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), power, fraction)
					}
				}
			}
			simApp.EndBlockAndCommit(ctx)

			// finalize rebalance
			simApp.BeginNextBlock()
			ctx = simApp.NewNextContext()
			simApp.EndBlockAndCommit(ctx)

			// assertions
			assertValidators(t, simApp, ctx, tt.want.vals)

			daoKeeper := simApp.OnomyApp().DaoKeeper
			gotTreasuryBalance := daoKeeper.Treasury(ctx)
			require.Equal(t, sdk.NewCoins(tt.want.treasuryBalance), gotTreasuryBalance)

			// pool rate = current pool / total
			require.Equal(t, gotTreasuryBalance[0].Amount.ToDec().Quo(daoKeeper.GetDaoDelegationSupply(ctx).Add(gotTreasuryBalance[0].Amount.ToDec())), types.DefaultStakingTokenPoolRate)

			// the check the overall balance remains the same
			require.Equal(t, daoKeeper.GetDaoDelegationSupply(ctx).Add(gotTreasuryBalance[0].Amount.ToDec()), tt.args.treasuryBalance.Amount.ToDec())
		})
	}
}

func createSimAppWithValidatorsAndTreasury(t *testing.T, vals map[string]simapp.ValReq, treasuryBalance sdk.Coin) (*simapp.SimApp, map[string]*secp256k1.PrivKey) {
	t.Helper()

	// treasury genesis
	treasuryOverrideOpt := simapp.WithGenesisOverride(
		func(m map[string]json.RawMessage) map[string]json.RawMessage {
			daoGenesis := types.DefaultGenesis()
			daoGenesis.TreasuryBalance = sdk.NewCoins(treasuryBalance)
			daoGenesisString, err := json.Marshal(daoGenesis)
			require.NoError(t, err)
			m[types.ModuleName] = daoGenesisString
			return m
		})

	return simapp.SetupWithValidators(t, vals, treasuryOverrideOpt)
}

type valAssertion struct {
	bondStatus     stakingtypes.BondStatus
	selfBondAmount sdk.Dec
	daoBondAmount  sdk.Dec
}

func assertValidators(t *testing.T, simApp *simapp.SimApp, ctx sdk.Context, vals map[string]valAssertion) {
	t.Helper()

	accountKeeper := simApp.OnomyApp().AccountKeeper
	stakingKeeper := simApp.OnomyApp().StakingKeeper
	daoAddress := accountKeeper.GetModuleAddress(types.ModuleName)
	updatedValidators := stakingKeeper.GetAllValidators(ctx)
	require.Equal(t, len(vals), len(updatedValidators))

	for _, val := range updatedValidators {
		delegations := stakingKeeper.GetValidatorDelegations(ctx, val.GetOperator())
		require.LessOrEqual(t, 1, len(delegations))

		valAssert, ok := vals[val.GetMoniker()]
		require.True(t, ok)
		require.Equal(t, valAssert.bondStatus, val.Status, val.GetMoniker())

		for _, delegation := range delegations {
			switch delegation.DelegatorAddress {
			case daoAddress.String():
				{
					require.Equal(t, valAssert.daoBondAmount, val.TokensFromShares(delegation.Shares), val.GetMoniker())
				}
			case sdk.AccAddress(val.GetOperator()).String():
				{
					require.Equal(t, valAssert.selfBondAmount, val.TokensFromShares(delegation.Shares), val.GetMoniker())
				}
			default:
				{
					t.Errorf("unexpected delegation %+v, of val %q, address: %q", delegation, val.GetMoniker(), val.GetOperator().String())
				}
			}
		}
	}
}
