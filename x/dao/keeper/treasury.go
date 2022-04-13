package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// Treasury returns the treasury balance.
func (k Keeper) Treasury(ctx sdk.Context) sdk.Coins {
	daoAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	return k.bankKeeper.GetAllBalances(ctx, daoAddress)
}
