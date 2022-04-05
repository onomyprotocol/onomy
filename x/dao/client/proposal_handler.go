// Package client contains dao client implementation.
package client

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	"github.com/onomyprotocol/onomy/x/dao/client/cli"
)

var (
	// FundTreasuryProposalHandler is the cli handler used for the gov cli integration.
	FundTreasuryProposalHandler = govclient.NewProposalHandler(cli.CmdFundTreasuryProposal, emptyRestHandler) // nolint:gochecknoglobals // cosmos-sdk style
	// ExchangeWithTreasuryProposalProposalHandler is the cli handler used for the gov cli integration.
	ExchangeWithTreasuryProposalProposalHandler = govclient.NewProposalHandler(cli.CmdExchangeWithTreasuryProposal, emptyRestHandler) // nolint:gochecknoglobals // cosmos-sdk style
)

func emptyRestHandler(client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "unsupported-dao-routes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Legacy REST Routes are not supported for DAO proposals")
		},
	}
}
