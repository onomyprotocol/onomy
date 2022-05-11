package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// VoteAbstain votes abstain on all the proposals from the DAO account.
func (k Keeper) VoteAbstain(ctx sdk.Context) (err error) {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	k.govKeeper.IterateProposals(ctx, func(proposal govtypes.Proposal) bool {
		if proposal.Status != govtypes.StatusVotingPeriod {
			return false
		}
		_, found := k.govKeeper.GetVote(ctx, proposal.ProposalId, daoAddr)
		// the dao should vote now
		if !found {
			err = k.govKeeper.AddVote(ctx, proposal.ProposalId, daoAddr, govtypes.NewNonSplitVoteOption(govtypes.OptionAbstain))
			if err != nil {
				return true
			}
			k.Logger(ctx).Info(fmt.Sprintf("voted abstain on proposal[%d]: %s", proposal.ProposalId, proposal.GetTitle()))
		}
		return false
	})

	return err
}
