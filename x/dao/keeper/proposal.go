package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

var _ types.QueryServer = Keeper{}

// FundTreasuryProposal submits the FundTreasuryProposal.
func (k Keeper) FundTreasuryProposal(ctx sdk.Context, request *types.FundTreasuryProposal) error {
	return nil
}

// ExchangeWithTreasuryProposal submits the ExchangeWithTreasuryProposal.
func (k Keeper) ExchangeWithTreasuryProposal(ctx sdk.Context, request *types.ExchangeWithTreasuryProposal) error {
	return nil
}
