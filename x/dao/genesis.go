package dao

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/keeper"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

// InitGenesis initializes the dao module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitGenesis(ctx, genState)
}

// ExportGenesis returns the dao module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return k.ExportGenesis(ctx)
}
