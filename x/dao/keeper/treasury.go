package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// Treasury returns the treasury balance.
func (k Keeper) Treasury(ctx sdk.Context) sdk.Coins {
	return k.treasury(ctx)
}

func (k Keeper) treasuryBondDenomAmount(ctx sdk.Context) sdk.Int {
	denom := k.stakingKeeper.BondDenom(ctx)
	return k.treasury(ctx).AmountOf(denom)
}

// treasury returns the treasury balance.
func (k Keeper) treasury(ctx sdk.Context) sdk.Coins {
	daoAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	return k.bankKeeper.GetAllBalances(ctx, daoAddress)
}

// GetBalance returns account balance.
func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, addr, denom)
}

// SendCoins transfers coins from one account to another.
func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}
