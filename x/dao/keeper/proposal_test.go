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
