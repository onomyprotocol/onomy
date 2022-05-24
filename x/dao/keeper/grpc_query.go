package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// QueryServer is keep wrapper which provides query capabilities.
type QueryServer struct {
	keeper Keeper
}

// NewQueryServer creates a new instance of QueryServer.
func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{
		keeper: keeper,
	}
}

var _ types.QueryServer = QueryServer{}

// Params return dao module current params values.
func (q QueryServer) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: q.keeper.GetParams(ctx)}, nil
}

// Treasury returns the treasury balance.
func (q QueryServer) Treasury(c context.Context, _ *types.QueryTreasuryRequest) (*types.QueryTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryTreasuryResponse{
		TreasuryBalance: q.keeper.Treasury(ctx),
	}, nil
}
