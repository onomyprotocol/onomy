package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

var _ types.QueryServer = Keeper{}

// Params return dao module current params values.
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// Treasury returns the treasury balance.
func (k Keeper) Treasury(c context.Context, _ *types.QueryTreasuryRequest) (*types.QueryTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	daoAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	daoBalance := k.bankKeeper.GetAllBalances(ctx, daoAddress)
	return &types.QueryTreasuryResponse{
		TreasuryBalance: daoBalance,
	}, nil
}
