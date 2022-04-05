package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

var _ types.QueryServer = Keeper{}

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
	return nil
}
