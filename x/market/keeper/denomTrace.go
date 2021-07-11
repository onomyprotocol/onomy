package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onomyprotocol/onomy/x/market/types"
)

// SetDenomTrace set a specific denomTrace in the store from its index
func (k Keeper) SetDenomTrace(ctx sdk.Context, denomTrace types.DenomTrace) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomTraceKey))
	b := k.cdc.MustMarshalBinaryBare(&denomTrace)
	store.Set(types.KeyPrefix(denomTrace.Index), b)
}

// GetDenomTrace returns a denomTrace from its index
func (k Keeper) GetDenomTrace(ctx sdk.Context, index string) (val types.DenomTrace, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomTraceKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// DeleteDenomTrace removes a denomTrace from the store
func (k Keeper) RemoveDenomTrace(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomTraceKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllDenomTrace returns all denomTrace
func (k Keeper) GetAllDenomTrace(ctx sdk.Context) (list []types.DenomTrace) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomTraceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DenomTrace
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
