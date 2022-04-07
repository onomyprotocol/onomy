package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

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

	tests := []struct { //nolint:dupl // test template
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
			app := simapp.Setup(false)
			ctx := app.BaseApp.NewContext(false, tmproto.Header{})
			wctx := sdk.WrapSDKContext(ctx)
			require.NoError(t, app.BankKeeper.MintCoins(ctx, types.ModuleName, tt.args.accountBalance))

			senderAddr, err := sdk.AccAddressFromBech32(tt.args.sender)
			require.NoError(t, err)
			require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, tt.args.accountBalance))

			err = app.DaoKeeper.FundTreasuryProposal(ctx, &types.FundTreasuryProposal{
				Sender: tt.args.sender,
				Amount: tt.args.amount,
			})

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			require.NoError(t, err)

			got, err := app.DaoKeeper.Treasury(wctx, &types.QueryTreasuryRequest{})
			require.NoError(t, err)

			require.Equal(t, tt.wantTreasuryBalance, got.TreasuryBalance)
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
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 8)),
				accountBalance:  sdk.NewCoins(sdk.NewInt64Coin(denom2, 10)),
				sender:          account.String(),
				coinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 8),
						CoinBid: sdk.NewInt64Coin(denom2, 10),
					},
				},
			},
			wantTreasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom2, 10)),
			wantAccountBalance:  sdk.NewCoins(sdk.NewInt64Coin(denom1, 8)),
		},
		{
			name: "positive_exchange_multiple_pairs",
			args: args{
				treasuryBalance: sdk.NewCoins(
					sdk.NewInt64Coin(denom1, 50),
					sdk.NewInt64Coin(denom2, 50),
					sdk.NewInt64Coin(denom3, 50),
					sdk.NewInt64Coin(denom4, 50),
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
				// 50 - 8 - 4 + 5
				sdk.NewInt64Coin(denom1, 43),
				// 50 + 10 + 2
				sdk.NewInt64Coin(denom2, 62),
				// 50 - 5
				sdk.NewInt64Coin(denom3, 45),
				sdk.NewInt64Coin(denom4, 50),
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.BaseApp.NewContext(false, tmproto.Header{})
			wctx := sdk.WrapSDKContext(ctx)

			require.NoError(t, app.BankKeeper.MintCoins(ctx, types.ModuleName, tt.args.treasuryBalance))
			require.NoError(t, app.BankKeeper.MintCoins(ctx, types.ModuleName, tt.args.accountBalance))
			senderAddr, err := sdk.AccAddressFromBech32(tt.args.sender)
			require.NoError(t, err)
			require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, tt.args.accountBalance))

			err = app.DaoKeeper.ExchangeWithTreasuryProposal(ctx, &types.ExchangeWithTreasuryProposal{
				Sender:     tt.args.sender,
				CoinsPairs: tt.args.coinsPairs,
			})

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			got, err := app.DaoKeeper.Treasury(wctx, &types.QueryTreasuryRequest{})
			require.NoError(t, err)
			require.Equal(t, tt.wantTreasuryBalance, got.TreasuryBalance)

			senderBalance := app.BankKeeper.GetAllBalances(ctx, senderAddr)
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

	tests := []struct { // nolint:dupl // test template
		name               string
		args               args
		wantAccountBalance sdk.Coins
		wantErr            error
	}{
		{
			name: "positive_one_coin_full",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
				recipient:       account.String(),
				amount:          sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
			},
			wantAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
		},
		{
			name: "positive_one_coin_partial",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
				recipient:       account.String(),
				amount:          sdk.NewCoins(sdk.NewInt64Coin(denom1, 5)),
			},
			wantAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 5)),
		},
		{
			name: "positive_two_coins_partial",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 10), sdk.NewInt64Coin(denom2, 8)),
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.BaseApp.NewContext(false, tmproto.Header{})
			require.NoError(t, app.BankKeeper.MintCoins(ctx, types.ModuleName, tt.args.treasuryBalance))
			err := app.DaoKeeper.FundAccountProposal(ctx, &types.FundAccountProposal{
				Recipient: tt.args.recipient,
				Amount:    tt.args.amount,
			})

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			require.NoError(t, err)

			recipientAddr, err := sdk.AccAddressFromBech32(tt.args.recipient)
			require.NoError(t, err)
			got := app.BankKeeper.GetAllBalances(ctx, recipientAddr)
			require.Equal(t, tt.wantAccountBalance, got)
		})
	}
}
