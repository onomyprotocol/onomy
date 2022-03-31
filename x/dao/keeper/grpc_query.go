package keeper

import (
	"context"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

var _ types.QueryServer = Keeper{}

// Treasury returns the treasury balance.
func (k Keeper) Treasury(ctx context.Context, request *types.QueryTreasuryRequest) (*types.QueryTreasuryResponse, error) {
	return &types.QueryTreasuryResponse{}, nil
}
