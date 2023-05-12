package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// VoteAbstain votes abstain on all the proposals from the DAO account.
func (k Keeper) InflateDao(ctx sdk.Context) (err error) {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	daoBalance := k.bankKeeper.GetBalance(ctx, daoAddr, "nom")
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
