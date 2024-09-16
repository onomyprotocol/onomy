// Package client contains dao client implementation.
package dao

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/onomyprotocol/onomy/x/dao/client/cli"
	"github.com/onomyprotocol/onomy/x/dao/keeper"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

var (
	// FundTreasuryProposalHandler is the cli handler used for the gov cli integration.
	FundTreasuryProposalHandler = govclient.NewProposalHandler(cli.CmdFundTreasuryProposal) // nolint:gochecknoglobals // cosmos-sdk style
	// ExchangeWithTreasuryProposalProposalHandler is the cli handler used for the gov cli integration.
	ExchangeWithTreasuryProposalProposalHandler = govclient.NewProposalHandler(cli.CmdExchangeWithTreasuryProposal) // nolint:gochecknoglobals // cosmos-sdk style
)

// NewProposalHandler defines the dao proposal handler.
func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.FundTreasuryProposal:
			return k.FundTreasuryProposal(ctx, c)

		case *types.ExchangeWithTreasuryProposal:
			return k.ExchangeWithTreasuryProposal(ctx, c)

		case *types.FundAccountProposal:
			return k.FundAccountProposal(ctx, c)

		default:
			return errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ibc proposal content type: %T", c)
		}
	}
}
