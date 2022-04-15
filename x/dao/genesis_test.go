package dao_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestInitGenesis(t *testing.T) {
	const (
		denom1 = "denom1"
		denom2 = "denom2"
	)
	type args struct {
		genState types.GenesisState
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "init_and_check_bank",
			args: args{
				genState: types.GenesisState{
					Params:          types.DefaultParams(),
					TreasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 1), sdk.NewInt64Coin(denom2, 2)),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		simApp := simapp.Setup()
		ctx := simApp.NewContext()
		t.Run(tt.name, func(t *testing.T) {
			dao.InitGenesis(ctx, simApp.OnomyApp().DaoKeeper, tt.args.genState)
			exportedModuleBalance := simApp.OnomyApp().BankKeeper.GetAllBalances(ctx, simApp.OnomyApp().AccountKeeper.GetModuleAddress(types.ModuleName))
			require.Equal(t, tt.args.genState.TreasuryBalance, exportedModuleBalance)
			require.Equal(t, tt.args.genState.Params, simApp.OnomyApp().DaoKeeper.GetParams(ctx))
		})
	}
}

func TestInitAndExportGenesis(t *testing.T) {
	const (
		denom1 = "denom1"
		denom2 = "denom2"
	)
	type args struct {
		genState types.GenesisState
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "import_same_as_export",
			args: args{
				genState: types.GenesisState{
					TreasuryBalance: sdk.NewCoins(sdk.NewInt64Coin(denom1, 1), sdk.NewInt64Coin(denom2, 2)),
					Params:          types.DefaultParams(),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		simApp := simapp.Setup()
		ctx := simApp.NewContext()
		t.Run(tt.name, func(t *testing.T) {
			dao.InitGenesis(ctx, simApp.OnomyApp().DaoKeeper, tt.args.genState)
			exportedGenesis := dao.ExportGenesis(ctx, simApp.OnomyApp().DaoKeeper)
			require.Equal(t, &tt.args.genState, exportedGenesis)
		})
	}
}
