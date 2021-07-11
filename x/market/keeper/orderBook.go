package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onomyprotocol/onomy/x/market/types"
)

// SetOrderBook set a specific orderBook in the store from its index
func (k Keeper) SetOrderBook(ctx sdk.Context, orderBook types.OrderBook) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrderBookKey))
	b := k.cdc.MustMarshalBinaryBare(&orderBook)
	store.Set(types.KeyPrefix(orderBook.Index), b)
}

// GetOrderBook returns a orderBook from its index
func (k Keeper) GetOrderBook(ctx sdk.Context, index string) (val types.OrderBook, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrderBookKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// DeleteOrderBook removes a orderBook from the store
func (k Keeper) RemoveOrderBook(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrderBookKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllOrderBook returns all orderBook
func (k Keeper) GetAllOrderBook(ctx sdk.Context) (list []types.OrderBook) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrderBookKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrderBook
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
