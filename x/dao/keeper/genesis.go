package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// InitGenesis sets dao module information from genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) error {
	k.SetParams(ctx, genState.Params)

	balance := genState.TreasuryBalance
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, balance)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:          k.getParams(ctx),
		TreasuryBalance: k.treasury(ctx),
	}
}
