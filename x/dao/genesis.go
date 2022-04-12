package dao

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/keeper"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

// InitGenesis initializes the dao module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := k.InitGenesis(ctx, genState); err != nil {
		panic(fmt.Errorf("can't init genesis for %q module, genState: %+v, err: %w", types.ModuleName, genState, err))
	}
}

// ExportGenesis returns the dao module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return k.ExportGenesis(ctx)
}
