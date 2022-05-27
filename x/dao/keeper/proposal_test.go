package keeper_test

import (
	"encoding/json"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestKeeper_FundTreasuryProposal(t *testing.T) {
	const (
		denom1 = "denom1"
		denom2 = "denom2"
	)

	account := simapp.GenAccount()

	type args struct {
		accountBalance sdk.Coins
		sender         string
		amount         sdk.Coins
	}

	tests := []struct {
		name                string
		args                args
		wantTreasuryBalance sdk.Coins
		wantErr             error
	}{
		{
			name: "positive_one_coin_full",
			args: args{
				accountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
				sender:         account.String(),
				amount:         sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
			},
			wantTreasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
		},
		{
			name: "positive_one_coin_partial",
			args: args{
				accountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
				sender:         account.String(),
				amount:         sdk.NewCoins(sdk.NewInt64Coin(denom1, 5)),
			},
			wantTreasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 5)),
		},
		{
			name: "positive_two_coins_partial",
			args: args{
				accountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10), sdk.NewInt64Coin(denom2, 8)),
				sender:         account.String(),
				amount:         sdk.NewCoins(sdk.NewInt64Coin(denom1, 5), sdk.NewInt64Coin(denom2, 8)),
			},
			wantTreasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 5), sdk.NewInt64Coin(denom2, 8)),
		},
		{
			name: "negative_insufficient_balance",
			args: args{
				accountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
				sender:         account.String(),
				amount:         sdk.NewCoins(sdk.NewInt64Coin(denom1, 11)),
			},
			wantErr: sdkerrors.Wrapf(types.ErrInsufficientBalance, "sender balance is less than amount to send"),
		},
		{
			name: "negative_not_existing_token",
			args: args{
				accountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
				sender:         account.String(),
				amount:         sdk.NewCoins(sdk.NewInt64Coin(denom2, 10)),
			},
			wantErr: sdkerrors.Wrapf(types.ErrInsufficientBalance, "sender balance is less than amount to send"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			simApp := simapp.Setup()
			ctx := simApp.NewContext()
			require.NoError(t, simApp.OnomyApp().BankKeeper.MintCoins(ctx, types.ModuleName, tt.args.accountBalance))

			senderAddr, err := sdk.AccAddressFromBech32(tt.args.sender)
			require.NoError(t, err)
			require.NoError(t, simApp.OnomyApp().BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, tt.args.accountBalance))

			err = simApp.OnomyApp().DaoKeeper.FundTreasuryProposal(ctx, &types.FundTreasuryProposal{
				Sender: tt.args.sender,
				Amount: tt.args.amount,
			})

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err, err)

			got := simApp.OnomyApp().DaoKeeper.Treasury(ctx)
			require.Equal(t, tt.wantTreasuryBalance, got)
		})
	}
}

func TestKeeper_ExchangeWithTreasuryProposal(t *testing.T) {
	const (
		denom1 = "denom1"
		denom2 = "denom2"
		denom3 = "denom3"
		denom4 = "denom4"
	)

	account := simapp.GenAccount()

	type args struct {
		treasuryBalance sdk.Coins
		accountBalance  sdk.Coins
		sender          string
		coinsPairs      []types.CoinsExchangePair
	}

	tests := []struct {
		name                string
		args                args
		wantAccountBalance  sdk.Coins
		wantTreasuryBalance sdk.Coins
		wantErr             error
	}{
		{
			name: "positive_exchange_full",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10000)),
				accountBalance:  sdk.NewCoins(sdk.NewInt64Coin(denom2, 10)),
				sender:          account.String(),
				coinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 8),
						CoinBid: sdk.NewInt64Coin(denom2, 10),
					},
				},
			},
			wantTreasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom2, 10), sdk.NewInt64Coin(denom1, 9992)),
			wantAccountBalance:  sdk.NewCoins(sdk.NewInt64Coin(denom1, 8)),
		},
		{
			name: "positive_exchange_multiple_pairs",
			args: args{
				treasuryBalance: sdk.NewCoins(
					sdk.NewInt64Coin(denom1, 50000),
					sdk.NewInt64Coin(denom2, 50000),
					sdk.NewInt64Coin(denom3, 50000),
					sdk.NewInt64Coin(denom4, 50000),
				),
				accountBalance: sdk.NewCoins(
					sdk.NewInt64Coin(denom1, 100),
					sdk.NewInt64Coin(denom2, 100),
					sdk.NewInt64Coin(denom3, 100),
				),
				sender: account.String(),
				coinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 8),
						CoinBid: sdk.NewInt64Coin(denom2, 10),
					},
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 4),
						CoinBid: sdk.NewInt64Coin(denom2, 2),
					},
					{
						CoinAsk: sdk.NewInt64Coin(denom3, 5),
						CoinBid: sdk.NewInt64Coin(denom1, 5),
					},
				},
			},
			wantTreasuryBalance: sdk.NewCoins(
				// 50000 - 8 - 4 + 5
				sdk.NewInt64Coin(denom1, 49993),
				// 50000 + 10 + 2
				sdk.NewInt64Coin(denom2, 50012),
				// 50000 - 5
				sdk.NewInt64Coin(denom3, 49995),
				sdk.NewInt64Coin(denom4, 50000),
			),
			wantAccountBalance: sdk.NewCoins(
				// 100 + 8 + 4 - 5
				sdk.NewInt64Coin(denom1, 107),
				// 100 - 10 - 2
				sdk.NewInt64Coin(denom2, 88),
				// 100 + 5
				sdk.NewInt64Coin(denom3, 105),
			),
		},
		{
			name: "negative_insufficient_sender_balance",
			args: args{
				treasuryBalance: sdk.NewCoins(
					sdk.NewInt64Coin(denom1, 50),
					sdk.NewInt64Coin(denom2, 50),
				),
				accountBalance: sdk.NewCoins(
					sdk.NewInt64Coin(denom1, 4),
					sdk.NewInt64Coin(denom2, 12),
				),
				sender: account.String(),
				coinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 8),
						CoinBid: sdk.NewInt64Coin(denom2, 10),
					},
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 4),
						CoinBid: sdk.NewInt64Coin(denom2, 2),
					},
					{
						CoinAsk: sdk.NewInt64Coin(denom3, 5),
						CoinBid: sdk.NewInt64Coin(denom1, 5),
					},
				},
			},
			wantErr: sdkerrors.Wrapf(types.ErrInsufficientBalance, "sender balance is less than bid coins amount"),
		},
		{
			name: "negative_insufficient_treasury_balance",
			args: args{
				treasuryBalance: sdk.NewCoins(
					sdk.NewInt64Coin(denom1, 50),
					sdk.NewInt64Coin(denom3, 2),
				),
				accountBalance: sdk.NewCoins(
					sdk.NewInt64Coin(denom1, 100),
					sdk.NewInt64Coin(denom2, 100),
				),
				sender: account.String(),
				coinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 8),
						CoinBid: sdk.NewInt64Coin(denom2, 10),
					},
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 4),
						CoinBid: sdk.NewInt64Coin(denom2, 2),
					},
					{
						CoinAsk: sdk.NewInt64Coin(denom3, 5),
						CoinBid: sdk.NewInt64Coin(denom1, 5),
					},
				},
			},
			wantErr: sdkerrors.Wrapf(types.ErrInsufficientBalance, "treasury balance is less than ask coins amount"),
		},
		{
			name: "negative_prohibited_proposal_amount",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 100)),
				accountBalance:  sdk.NewCoins(sdk.NewInt64Coin(denom2, 10)),
				sender:          account.String(),
				coinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 8),
						CoinBid: sdk.NewInt64Coin(denom2, 10),
					},
				},
			},
			wantErr: sdkerrors.Wrapf(types.ErrProhibitedCoinsAmount, "requested denom1:8 amount is more than max allowed denom1:5 "),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			simApp := simapp.Setup()
			ctx := simApp.NewContext()

			simApp.OnomyApp().DaoKeeper.SetParams(ctx, types.DefaultParams())
			require.NoError(t, simApp.OnomyApp().BankKeeper.MintCoins(ctx, types.ModuleName, tt.args.treasuryBalance))
			require.NoError(t, simApp.OnomyApp().BankKeeper.MintCoins(ctx, types.ModuleName, tt.args.accountBalance))
			senderAddr, err := sdk.AccAddressFromBech32(tt.args.sender)
			require.NoError(t, err)
			require.NoError(t, simApp.OnomyApp().BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, tt.args.accountBalance))

			err = simApp.OnomyApp().DaoKeeper.ExchangeWithTreasuryProposal(ctx, &types.ExchangeWithTreasuryProposal{
				Sender:     tt.args.sender,
				CoinsPairs: tt.args.coinsPairs,
			})

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err, err)

			got := simApp.OnomyApp().DaoKeeper.Treasury(ctx)
			require.Equal(t, tt.wantTreasuryBalance, got)

			senderBalance := simApp.OnomyApp().BankKeeper.GetAllBalances(ctx, senderAddr)
			require.Equal(t, tt.wantAccountBalance, senderBalance)
		})
	}
}

func TestKeeper_FundAccountProposal(t *testing.T) {
	const (
		denom1 = "denom1"
		denom2 = "denom2"
	)

	account := simapp.GenAccount()

	type args struct {
		treasuryBalance sdk.Coins
		recipient       string
		amount          sdk.Coins
	}

	tests := []struct {
		name               string
		args               args
		wantAccountBalance sdk.Coins
		wantErr            error
	}{
		{
			name: "positive_one_coin_full",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10000)),
				recipient:       account.String(),
				amount:          sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
			},
			wantAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
		},
		{
			name: "positive_one_coin_partial",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10000)),
				recipient:       account.String(),
				amount:          sdk.NewCoins(sdk.NewInt64Coin(denom1, 5)),
			},
			wantAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 5)),
		},
		{
			name: "positive_two_coins_partial",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10000), sdk.NewInt64Coin(denom2, 80000)),
				recipient:       account.String(),
				amount:          sdk.NewCoins(sdk.NewInt64Coin(denom1, 5), sdk.NewInt64Coin(denom2, 8)),
			},
			wantAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 5), sdk.NewInt64Coin(denom2, 8)),
		},
		{
			name: "negative_insufficient_balance",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
				recipient:       account.String(),
				amount:          sdk.NewCoins(sdk.NewInt64Coin(denom1, 11)),
			},
			wantErr: sdkerrors.Wrapf(types.ErrInsufficientBalance, "treasury balance is less than amount to send"),
		},
		{
			name: "negative_not_existing_token",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
				recipient:       account.String(),
				amount:          sdk.NewCoins(sdk.NewInt64Coin(denom2, 10)),
			},
			wantErr: sdkerrors.Wrapf(types.ErrInsufficientBalance, "treasury balance is less than amount to send"),
		},
		{
			name: "negative_prohibited_proposal_amount",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 50)),
				recipient:       account.String(),
				amount:          sdk.NewCoins(sdk.NewInt64Coin(denom1, 5)),
			},
			wantErr: sdkerrors.Wrapf(types.ErrProhibitedCoinsAmount, "requested denom1:5 amount is more than max allowed denom1:2 "),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			simApp := simapp.Setup()
			ctx := simApp.NewContext()
			simApp.OnomyApp().DaoKeeper.SetParams(ctx, types.DefaultParams())

			require.NoError(t, simApp.OnomyApp().BankKeeper.MintCoins(ctx, types.ModuleName, tt.args.treasuryBalance))
			err := simApp.OnomyApp().DaoKeeper.FundAccountProposal(ctx, &types.FundAccountProposal{
				Recipient: tt.args.recipient,
				Amount:    tt.args.amount,
			})

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err, err)

			recipientAddr, err := sdk.AccAddressFromBech32(tt.args.recipient)
			require.NoError(t, err)
			got := simApp.OnomyApp().BankKeeper.GetAllBalances(ctx, recipientAddr)
			require.Equal(t, tt.wantAccountBalance, got)
		})
	}
}

func TestKeeper_ProposalsFullCycle(t *testing.T) {
	var (
		acc               = simapp.GenAccount()
		tenPercents       = sdk.NewDec(1).QuoInt64(10)
		oneBondCoin       = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction))
		twoBondCoins      = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(2, sdk.DefaultPowerReduction))
		tenBondCoins      = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
		thousandBondCoins = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction))
		commission        = stakingtypes.NewCommissionRates(tenPercents, tenPercents, tenPercents)
	)

	type args struct {
		treasuryBalance sdk.Coins
		proposal        func(proposer sdk.AccAddress) govtypes.Content
	}

	tests := []struct {
		name                string
		args                args
		wantTreasuryBalance sdk.Coins
	}{
		{
			name: "positive_fund_treasury",
			args: args{
				proposal: func(proposer sdk.AccAddress) govtypes.Content {
					return types.NewFundTreasuryProposal(proposer, "title", "desc", sdk.NewCoins(oneBondCoin))
				},
			},
			// the expectation includes the DefaultPoolRate logic because of the end-blocker rebalancer,
			// also we add 1 because of the redelegation truncation
			wantTreasuryBalance: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, oneBondCoin.
				Amount.ToDec().Mul(types.DefaultPoolRate).TruncateInt().AddRaw(1))),
		},
		{
			name: "positive_exchange_with_treasury",
			args: args{
				treasuryBalance: sdk.NewCoins(thousandBondCoins),
				proposal: func(proposer sdk.AccAddress) govtypes.Content {
					return types.NewExchangeWithTreasuryProposal(proposer, "title", "desc", []types.CoinsExchangePair{
						{
							CoinAsk: oneBondCoin,
							CoinBid: twoBondCoins,
						},
					})
				},
			},
			// the expectation includes the DefaultPoolRate logic because of the end-blocker rebalancer,
			// also we add 1 because of the redelegation truncation
			wantTreasuryBalance: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, thousandBondCoins.Sub(oneBondCoin).Add(twoBondCoins).
				Amount.ToDec().Mul(types.DefaultPoolRate).TruncateInt())),
		},
		{
			name: "positive_fund_account",
			args: args{
				treasuryBalance: sdk.NewCoins(thousandBondCoins),
				proposal: func(_ sdk.AccAddress) govtypes.Content {
					// we can any acc for that proposal
					return types.NewFundAccountProposal(acc, "title", "desc", sdk.NewCoins(oneBondCoin))
				},
			},
			// the expectation includes the DefaultPoolRate logic because of the end-blocker rebalancer,
			// also we add 1 because of the redelegation truncation
			wantTreasuryBalance: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, thousandBondCoins.Sub(oneBondCoin).
				Amount.ToDec().Mul(types.DefaultPoolRate).TruncateInt().AddRaw(1))),
		},
		{
			name: "positive_change_pool_rate_param",
			args: args{
				treasuryBalance: sdk.NewCoins(thousandBondCoins),
				proposal: func(_ sdk.AccAddress) govtypes.Content {
					return proposal.NewParameterChangeProposal("title", "desc", []proposal.ParamChange{
						{
							Subspace: types.ModuleName,
							Key:      string(types.KeyPoolRate),
							Value:    `"0.5"`,
						},
					})
				},
			},
			// new pool rate must change the treasury pool
			wantTreasuryBalance: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, thousandBondCoins.Amount.QuoRaw(2).AddRaw(1))),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const (
				proposer = "proposer"
				voter    = "voter"
			)

			vals := map[string]simapp.ValReq{
				// low voting poser
				proposer: {
					SelfBondCoin: oneBondCoin,
					Commission:   commission,
					Balance:      sdk.NewCoins(tenBondCoins),
				},
				// high voting poser
				voter: {
					SelfBondCoin: tenBondCoins,
					Commission:   commission,
					Balance:      sdk.NewCoins(tenBondCoins.Add(tenBondCoins)),
				},
			}

			options := make([]simapp.Option, 0)
			if tt.args.treasuryBalance != nil {
				// treasury genesis
				treasuryOverrideOpt := simapp.WithGenesisOverride(
					func(m map[string]json.RawMessage) map[string]json.RawMessage {
						daoGenesis := types.DefaultGenesis()
						daoGenesis.TreasuryBalance = tt.args.treasuryBalance
						daoGenesisString, err := json.Marshal(daoGenesis)
						require.NoError(t, err)
						m[types.ModuleName] = daoGenesisString
						return m
					})
				options = append(options, treasuryOverrideOpt)
			}

			simApp, privs := simapp.SetupWithValidators(t, vals, options...)

			privProposer := privs[proposer]
			privVoter := privs[voter]

			simApp.BeginNextBlock()
			ctx := simApp.CurrentContext()

			addressProposer := sdk.AccAddress(privProposer.PubKey().Address())
			simApp.CreateProposal(t, tt.args.proposal(addressProposer), oneBondCoin, privProposer)

			simApp.EndBlock(ctx)

			proposalID := uint64(1)

			// vote
			simApp.BeginNextBlock()
			ctx = simApp.CurrentContext()
			simApp.VoteProposal(t, proposalID, govtypes.OptionYes, privVoter)
			simApp.EndBlock(ctx)

			// execute the proposal
			simApp.BeginNextBlock()
			ctx = simApp.CurrentContext()
			// +30 days to finish the voting
			ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 24 * 30))
			gov.EndBlocker(ctx, simApp.OnomyApp().GovKeeper)
			// assert the proposal
			proposal, _ := simApp.OnomyApp().GovKeeper.GetProposal(ctx, proposalID)
			require.Equal(t, govtypes.StatusPassed, proposal.Status)

			// execute the proposal
			simApp.BeginNextBlock()
			ctx = simApp.CurrentContext()
			simApp.EndBlock(ctx)

			// assert treasury
			require.Equal(t, tt.wantTreasuryBalance.String(), simApp.OnomyApp().DaoKeeper.Treasury(ctx).String())
		})
	}
}
