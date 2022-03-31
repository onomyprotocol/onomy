package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestKeeper_Treasury(t *testing.T) {
	const (
		denom1 = "denom1"
		denom2 = "denom2"
	)

	type args struct {
		treasuryBalance sdk.Coins
	}

	tests := []struct {
		name string
		args args
		want *types.QueryTreasuryResponse
	}{
		{
			name: "get_from_genesis",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 1), sdk.NewInt64Coin(denom2, 2)),
			},
			want: &types.QueryTreasuryResponse{
				TreasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 1), sdk.NewInt64Coin(denom2, 2)),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app, ctx := simapp.Setup()
			wctx := sdk.WrapSDKContext(ctx)
			app.DaoKeeper.InitGenesis(ctx, types.GenesisState{
				TreasuryBalance: tt.args.treasuryBalance,
			})

			got, err := app.DaoKeeper.Treasury(wctx, &types.QueryTreasuryRequest{})
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
