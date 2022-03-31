package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// InitGenesis sets dao module information from genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	balance := genState.TreasuryBalance
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, balance); err != nil {
		panic(fmt.Errorf("can't mint contins from %q module in InitGenesis", types.ModuleName))
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	daoAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	daoBalance := k.bankKeeper.GetAllBalances(ctx, daoAddress)
	return &types.GenesisState{
		TreasuryBalance: daoBalance,
	}
}
