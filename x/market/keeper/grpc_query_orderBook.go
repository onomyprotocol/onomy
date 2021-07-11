package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/onomyprotocol/onomy/x/market/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OrderBookAll(c context.Context, req *types.QueryAllOrderBookRequest) (*types.QueryAllOrderBookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var orderBooks []*types.OrderBook
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	orderBookStore := prefix.NewStore(store, types.KeyPrefix(types.OrderBookKey))

	pageRes, err := query.Paginate(orderBookStore, req.Pagination, func(key []byte, value []byte) error {
		var orderBook types.OrderBook
		if err := k.cdc.UnmarshalBinaryBare(value, &orderBook); err != nil {
			return err
		}

		orderBooks = append(orderBooks, &orderBook)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOrderBookResponse{OrderBook: orderBooks, Pagination: pageRes}, nil
}

func (k Keeper) OrderBook(c context.Context, req *types.QueryGetOrderBookRequest) (*types.QueryGetOrderBookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetOrderBook(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetOrderBookResponse{OrderBook: &val}, nil
}
