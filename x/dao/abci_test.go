package dao_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestEndBlocker(t *testing.T) {
	fiftyPercents := sdk.NewDec(1).Quo(sdk.NewDec(2))
	tenPercents := sdk.NewDec(1).Quo(sdk.NewDec(10))
	genesisCoins := sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction))
	lowBondCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000) // not enough for validator to be bonded
	normalBondCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	lowCommission := stakingtypes.NewCommissionRates(tenPercents, tenPercents, tenPercents)
	highCommission := stakingtypes.NewCommissionRates(fiftyPercents, fiftyPercents, fiftyPercents)

	type valReq struct {
		selfBondCoin sdk.Coin
		commission   stakingtypes.CommissionRates
	}

	type valAssertion struct {
		bondStatus     stakingtypes.BondStatus
		selfBondAmount sdk.Dec
		daoBondAmount  sdk.Dec
	}

	type args struct {
		vals         map[string]valReq
		treasuryBond sdk.Coin
	}

	type wantAssertion struct {
		vals         map[string]valAssertion
		treasuryBond sdk.Coin
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
						selfBondCoin: normalBondCoin,
						commission:   lowCommission,
					},
					"val2": { // bonded
						selfBondCoin: normalBondCoin,
						commission:   lowCommission,
					},
					"val3": { // won't be bonded
						selfBondCoin: lowBondCoin,
						commission:   lowCommission,
					},
					"val4": { // bonded, but high commission to be staked
						selfBondCoin: normalBondCoin,
						commission:   highCommission,
					},
				},
				treasuryBond: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)),
			},
			want: wantAssertion{
				vals: map[string]valAssertion{
					"val1": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: normalBondCoin.Amount.ToDec(),
						daoBondAmount:  sdk.TokensFromConsensusPower(95, sdk.DefaultPowerReduction).ToDec().Quo(sdk.NewDec(2)),
					},
					"val2": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: normalBondCoin.Amount.ToDec(),
						daoBondAmount:  sdk.TokensFromConsensusPower(95, sdk.DefaultPowerReduction).ToDec().Quo(sdk.NewDec(2)),
					},
					"val3": {
						bondStatus:     stakingtypes.Unbonded,
						selfBondAmount: lowBondCoin.Amount.ToDec(),
						daoBondAmount:  sdk.ZeroDec(),
					},
					"val4": {
						bondStatus:     stakingtypes.Bonded,
						selfBondAmount: normalBondCoin.Amount.ToDec(),
						daoBondAmount:  sdk.ZeroDec(),
					},
				},
				treasuryBond: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(5, sdk.DefaultPowerReduction)),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// prepare account genesis params
			balances := make([]banktypes.Balance, 0, len(tt.args.vals))
			privateKeys := make(map[string]*secp256k1.PrivKey, len(tt.args.vals))
			for moniker := range tt.args.vals {
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
					daoGenesis.TreasuryBalance = sdk.NewCoins(tt.args.treasuryBond)
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

			for moniker, val := range tt.args.vals {
				description := stakingtypes.Description{Moniker: moniker}
				simApp.CreateValidator(t, val.selfBondCoin, description, val.commission, sdk.OneInt(), privateKeys[moniker])
			}

			simApp.BeginNextBlock()
			ctx := simApp.NewContext()
			dao.EndBlocker(ctx, simApp.OnomyApp().DaoKeeper)

			// the chain is running
			// assertions

			accountKeeper := simApp.OnomyApp().AccountKeeper
			daoKeeper := simApp.OnomyApp().DaoKeeper
			stakingKeeper := simApp.OnomyApp().StakingKeeper
			daoAddress := accountKeeper.GetModuleAddress(types.ModuleName)

			updatedValidators := stakingKeeper.GetAllValidators(ctx)
			assert.Equal(t, len(tt.args.vals), len(updatedValidators))

			for _, val := range updatedValidators {
				delegations := stakingKeeper.GetValidatorDelegations(ctx, val.GetOperator())
				assert.LessOrEqual(t, 1, len(delegations))

				valAssert, ok := tt.want.vals[val.GetMoniker()]
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

				treasuryBond := daoKeeper.Treasury(ctx).AmountOf(stakingKeeper.BondDenom(ctx)).ToDec()
				assert.Equal(t, tt.want.treasuryBond.Amount.ToDec(), treasuryBond)

				// TBD:
				// check that sum of all coins in the staking and remaining are eq initial pool
				// check that all staking queues are empty
				// add tests with the unbonding -> bonding validator
				// add tests with the treasury bond pool change
				// add test with the redelegation
				// add tests with decimals rounding
			}
		})
	}
}
