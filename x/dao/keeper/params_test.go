package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestKeeper_GetAndSetParams(t *testing.T) {
	type args struct {
		params types.Params
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive",
			args: args{
				params: types.DefaultParams(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.BaseApp.NewContext(false, tmproto.Header{})

			app.DaoKeeper.SetParams(ctx, tt.args.params)
			got := app.DaoKeeper.GetParams(ctx)

			require.Equal(t, tt.args.params, got)
		})
	}
}
