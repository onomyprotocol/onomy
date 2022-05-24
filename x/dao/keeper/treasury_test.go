package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

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
		want sdk.Coins
	}{
		{
			name: "get_from_genesis",
			args: args{
				treasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 1), sdk.NewInt64Coin(denom2, 2)),
			},
			want: sdk.NewCoins(sdk.NewInt64Coin(denom1, 1), sdk.NewInt64Coin(denom2, 2)),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			simApp := simapp.Setup()
			ctx := simApp.OnomyApp().BaseApp.NewContext(false, tmproto.Header{})

			err := simApp.OnomyApp().DaoKeeper.InitGenesis(ctx, types.GenesisState{
				Params:          types.DefaultParams(),
				TreasuryBalance: tt.args.treasuryBalance,
			})
			require.NoError(t, err)

			got := simApp.OnomyApp().DaoKeeper.Treasury(ctx)
			require.Equal(t, tt.want, got)
		})
	}
}
