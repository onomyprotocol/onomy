package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	testkeeper "github.com/onomyprotocol/onomy/testutil/keeper"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestTreasury(t *testing.T) {
	keeper, ctx := testkeeper.DaoKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	response, err := keeper.Treasury(wctx, &types.QueryTreasuryRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryTreasuryResponse{}, response)
}
