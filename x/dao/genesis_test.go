package dao_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/onomyprotocol/onomy/testutil/keeper"
	"github.com/onomyprotocol/onomy/testutil/nullify"
	"github.com/onomyprotocol/onomy/x/dao"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := keepertest.DaoKeeper(t)
	dao.InitGenesis(ctx, *k, genesisState)
	got := dao.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
