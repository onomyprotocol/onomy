package dao_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

var (
	fiftyPercents  = sdk.NewDec(1).QuoInt64(2)                                                                       //nolint:gochecknoglobals
	tenPercents    = sdk.NewDec(1).Quo(sdk.NewDec(10))                                                               //nolint:gochecknoglobals
	genesisCoins   = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)) //nolint:gochecknoglobals
	nanoBondCoins  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000000)                                              // not enough for validator to be bonded
	tenBondCoins   = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))  //nolint:gochecknoglobals
	lowCommission  = stakingtypes.NewCommissionRates(tenPercents, tenPercents, tenPercents)                          //nolint:gochecknoglobals
	highCommission = stakingtypes.NewCommissionRates(fiftyPercents, fiftyPercents, fiftyPercents)                    //nolint:gochecknoglobals
)

type valReq struct {
	selfBondCoin sdk.Coin
	commission   stakingtypes.CommissionRates
	reward       sdk.Coin
}

type valAssertion struct {
	bondStatus     stakingtypes.BondStatus
	selfBondAmount sdk.Dec
	daoBondAmount  sdk.Dec
}

func TestEndBlocker_ReBalance(t *testing.T) {
	type args struct {
		vals            map[string]valReq
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
				vals: map[string]valReq{
					"val1": { // bonded
						selfBondCoin: tenBondCoins,
						commission:   lowCommission,
					},
					"val2": { // bonded
						selfBondCoin: tenBondCoins,
						commission:   lowCommission,
					},
					"val3": { // won't be bonded
						selfBondCoin: nanoBondCoins,
						commission:   lowCommission,
					},
					"val4": { // bonded, but high commission to be staked
						selfBondCoin: tenBondCoins,
						commission:   highCommission,
					},
				},
				treasuryBalance: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)),
			},
			want: wantAssertion{
				vals: map[string]valAssertion{
					"val1": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						daoBondAmount:  sdk.TokensFromConsensusPower(95, sdk.DefaultPowerReduction).ToDec().QuoInt64(2),
					},
					"val2": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						daoBondAmount:  sdk.TokensFromConsensusPower(95, sdk.DefaultPowerReduction).ToDec().QuoInt64(2),
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
				treasuryBalance: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(5, sdk.DefaultPowerReduction)),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			simApp := createSimAppWithValidators(t, tt.args.vals, tt.args.treasuryBalance)

			simApp.BeginNextBlock()
			ctx := simApp.NewContext()
			dao.EndBlocker(ctx, simApp.OnomyApp().DaoKeeper)

			// assertions
			assertValidators(t, simApp, ctx, tt.want.vals)

			daoKeeper := simApp.OnomyApp().DaoKeeper
			gotTreasuryBalance := daoKeeper.Treasury(ctx)
			assert.Equal(t, sdk.NewCoins(tt.want.treasuryBalance), gotTreasuryBalance)

			// the remaining pool is expected
			// (staked + current) * pool rate = current
			assert.Equal(t, daoKeeper.GetDaoDelegationSupply(ctx).Add(gotTreasuryBalance[0].Amount.ToDec()).Mul(types.DefaultStakingTokenPoolRate).RoundInt().ToDec(),
				gotTreasuryBalance[0].Amount.ToDec())

			// TBD:
			// add tests with the unbonding -> bonding validator
			// add test with the redelegation
			// add tests with decimals rounding
		})
	}
}

func TestEndBlocker_WithdrawReward(t *testing.T) {
	validatorReward := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000000)
	expectedDaoFullReward := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1486956434)

	type args struct {
		vals            map[string]valReq
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
				vals: map[string]valReq{
					"val1": { // bonded
						selfBondCoin: tenBondCoins,
						commission:   lowCommission,
						reward:       validatorReward,
					},
					"val2": { // bonded
						selfBondCoin: tenBondCoins,
						commission:   lowCommission,
						reward:       validatorReward,
					},
					"val3": { // won't be bonded
						selfBondCoin: nanoBondCoins,
						commission:   lowCommission,
						reward:       validatorReward,
					},
				},
				treasuryBalance: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)),
			},
			want: wantAssertion{
				vals: map[string]valAssertion{
					"val1": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						daoBondAmount:
						// initial dao staking
						sdk.TokensFromConsensusPower(95, sdk.DefaultPowerReduction).ToDec().QuoInt64(2).
							// the reward
							Add(expectedDaoFullReward.Amount.ToDec().QuoInt64(2).Mul(sdk.OneDec().Sub(types.DefaultStakingTokenPoolRate))).TruncateDec(),
					},
					"val2": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: tenBondCoins.Amount.ToDec(),
						daoBondAmount:
						// initial dao staking
						sdk.TokensFromConsensusPower(95, sdk.DefaultPowerReduction).ToDec().QuoInt64(2).
							// the reward
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
			simApp := createSimAppWithValidators(t, tt.args.vals, tt.args.treasuryBalance)
			simApp.BeginNextBlock()
			ctx := simApp.NewNextContext()
			// update dao params to withdraw reward
			daoKeeper := simApp.OnomyApp().DaoKeeper
			daoParams := daoKeeper.GetParams(ctx)
			daoParams.WithdrawRewardPeriod = withdrawRewardPeriod // withdraw reward each 10 block
			daoKeeper.SetParams(ctx, daoParams)
			// allocate validator rewards
			for moniker := range tt.args.vals {
				moniker := moniker
				simApp.OnomyApp().StakingKeeper.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
					if moniker == validator.GetMoniker() {
						// workaround to fix the distribution invariant
						require.NoError(t, simApp.OnomyApp().BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(tt.args.vals[moniker].reward)))
						require.NoError(t, simApp.OnomyApp().BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distrtypes.ModuleName, sdk.NewCoins(tt.args.vals[moniker].reward)))

						simApp.OnomyApp().DistrKeeper.AllocateTokensToValidator(ctx, validator, sdk.NewDecCoinsFromCoins(tt.args.vals[moniker].reward))
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
			assert.Equal(t, sdk.NewCoins(tt.want.treasuryBalance), gotTreasuryBalance)

			// the remaining pool is expected
			// (staked + current) * pool rate = current
			assert.Equal(t, daoKeeper.GetDaoDelegationSupply(ctx).Add(gotTreasuryBalance[0].Amount.ToDec()).Mul(types.DefaultStakingTokenPoolRate).RoundInt().ToDec(),
				gotTreasuryBalance[0].Amount.ToDec())
		})
	}
}

func assertValidators(t *testing.T, simApp *simapp.SimApp, ctx sdk.Context, vals map[string]valAssertion) {
	t.Helper()

	accountKeeper := simApp.OnomyApp().AccountKeeper
	stakingKeeper := simApp.OnomyApp().StakingKeeper
	daoAddress := accountKeeper.GetModuleAddress(types.ModuleName)
	updatedValidators := stakingKeeper.GetAllValidators(ctx)
	assert.Equal(t, len(vals), len(updatedValidators))

	for _, val := range updatedValidators {
		delegations := stakingKeeper.GetValidatorDelegations(ctx, val.GetOperator())
		assert.LessOrEqual(t, 1, len(delegations))

		valAssert, ok := vals[val.GetMoniker()]
		assert.True(t, ok)
		assert.Equal(t, valAssert.bondStatus, val.Status, val.GetMoniker())

		for _, delegation := range delegations {
			switch delegation.DelegatorAddress {
			case daoAddress.String():
				{
					assert.Equal(t, valAssert.daoBondAmount, delegation.Shares)
				}
			case sdk.AccAddress(val.GetOperator()).String():
				{
					assert.Equal(t, valAssert.selfBondAmount, delegation.Shares)
				}
			default:
				{
					t.Errorf("unexpected delegation %+v, of val %q, address: %q", delegation, val.GetMoniker(), val.GetOperator().String())
				}
			}
		}
	}
}

func createSimAppWithValidators(t *testing.T, vals map[string]valReq, treasuryBalance sdk.Coin) *simapp.SimApp {
	t.Helper()

	// prepare account genesis params
	balances := make([]banktypes.Balance, 0, len(vals))
	privateKeys := make(map[string]*secp256k1.PrivKey, len(vals))
	for moniker := range vals {
		privateKey := secp256k1.GenPrivKey()
		privateKeys[moniker] = privateKey
		address := sdk.AccAddress(privateKey.PubKey().Address())
		balances = append(balances, banktypes.Balance{
			Address: address.String(),
			Coins:   sdk.Coins{genesisCoins},
		})
	}
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

	simApp := simapp.Setup(
		simapp.WithGenesisAccountsAndBalances(balances...),
		treasuryOverrideOpt,
		simapp.WithAppCommit(),
	)

	for moniker, val := range vals {
		description := stakingtypes.Description{Moniker: moniker}
		simApp.CreateValidator(t, val.selfBondCoin, description, val.commission, sdk.OneInt(), privateKeys[moniker])
	}
	return simApp
}
