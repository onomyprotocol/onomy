package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// InflateDao inflates treasury by APR from minter.
func (k Keeper) InflateDao(ctx context.Context) (err error) {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	daoBalance := k.bankKeeper.GetBalance(ctx, daoAddr, "anom")
	minter := k.mintKeeper.GetMinter(ctx)
	params := k.mintKeeper.GetParams(ctx)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, daoBalance.Amount)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, mintedCoins)
	if err != nil {
		panic(err)
	}

	return err
}
