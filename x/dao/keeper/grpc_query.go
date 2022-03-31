package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

var _ types.QueryServer = Keeper{}

// Treasury returns the treasury balance.
func (k Keeper) Treasury(c context.Context, _ *types.QueryTreasuryRequest) (*types.QueryTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	daoAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	daoBalance := k.bankKeeper.GetAllBalances(ctx, daoAddress)
	return &types.QueryTreasuryResponse{
		TreasuryBalance: daoBalance,
	}, nil
}
