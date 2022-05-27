package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// FundTreasuryProposal submits the FundTreasuryProposal.
func (k Keeper) FundTreasuryProposal(ctx sdk.Context, request *types.FundTreasuryProposal) error {
	senderAddr, err := sdk.AccAddressFromBech32(request.Sender)
	if err != nil {
		return err
	}

	senderBalance := k.bankKeeper.GetAllBalances(ctx, senderAddr)
	amountToSend := request.Amount
	if _, isNegative := senderBalance.SafeSub(amountToSend); isNegative {
		return sdkerrors.Wrapf(types.ErrInsufficientBalance, "sender balance is less than amount to send")
	}

	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, amountToSend)
}

// ExchangeWithTreasuryProposal submits the ExchangeWithTreasuryProposal.
func (k Keeper) ExchangeWithTreasuryProposal(ctx sdk.Context, request *types.ExchangeWithTreasuryProposal) error {
	senderAddr, err := sdk.AccAddressFromBech32(request.Sender)
	if err != nil {
		return err
	}

	coinsAsk := sdk.NewCoins()
	coinsBid := sdk.NewCoins()
	for i := range request.CoinsPairs {
		coinsAsk = coinsAsk.Add(request.CoinsPairs[i].CoinAsk)
		coinsBid = coinsBid.Add(request.CoinsPairs[i].CoinBid)
	}

	senderBalance := k.bankKeeper.GetAllBalances(ctx, senderAddr)
	if _, isNegative := senderBalance.SafeSub(coinsBid); isNegative {
		return sdkerrors.Wrapf(types.ErrInsufficientBalance, "sender balance is less than bid coins amount")
	}

	treasuryBalance := k.bankKeeper.GetAllBalances(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName))
	if _, isNegative := treasuryBalance.SafeSub(coinsAsk); isNegative {
		return sdkerrors.Wrapf(types.ErrInsufficientBalance, "treasury balance is less than ask coins amount")
	}

	if err := k.validateMaxProposalRate(ctx, treasuryBalance, coinsAsk); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coinsBid); err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coinsAsk)
}

// FundAccountProposal submits the FundAccountProposal.
func (k Keeper) FundAccountProposal(ctx sdk.Context, request *types.FundAccountProposal) error {
	recipientAddr, err := sdk.AccAddressFromBech32(request.Recipient)
	if err != nil {
		return err
	}

	treasuryBalance := k.bankKeeper.GetAllBalances(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName))
	amountToSend := request.Amount
	if _, isNegative := treasuryBalance.SafeSub(amountToSend); isNegative {
		return sdkerrors.Wrapf(types.ErrInsufficientBalance, "treasury balance is less than amount to send")
	}
	if err := k.validateMaxProposalRate(ctx, treasuryBalance, amountToSend); err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, amountToSend)
}

func (k Keeper) validateMaxProposalRate(ctx sdk.Context, treasuryBalance, coinsAsk sdk.Coins) error {
	maxProposalRate := k.MaxProposalRate(ctx)
	for _, tcoin := range treasuryBalance {
		askAmount := coinsAsk.AmountOf(tcoin.Denom)
		allowedAmount := tcoin.Amount.ToDec().Mul(maxProposalRate).TruncateInt()
		if allowedAmount.LT(askAmount) {
			return sdkerrors.Wrapf(types.ErrProhibitedCoinsAmount, "requested %s:%s amount is more than max allowed %s:%s ",
				tcoin.Denom, askAmount, tcoin.Denom, allowedAmount.String())
		}
	}
	return nil
}
