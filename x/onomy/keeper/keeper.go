// Package keeper is onomy cosmos sdk keeper.
package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/onomyprotocol/onomy/x/onomy/types"
)

type (
	// Keeper is a struct that holds module store actions.
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
		// this line is used by starport scaffolding # ibc/keeper/attribute

	}
)

// NewKeeper creates an instance of Keeper.
func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	// this line is used by starport scaffolding # ibc/keeper/parameter

) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		// this line is used by starport scaffolding # ibc/keeper/return

	}
}

// Logger logger wrapped with the keeper.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
