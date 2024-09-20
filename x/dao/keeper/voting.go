package keeper

import (
	"context"
	"fmt"

	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

// VoteAbstain votes abstain on all the proposals from the DAO account.
func (k Keeper) VoteAbstain(ctx context.Context) (err error) {
	daoAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	k.govKeeper.IterateProposals(ctx, func(proposal v1.Proposal) bool {
		if proposal.Status != v1.StatusVotingPeriod {
			return false
		}
		_, err := k.govKeeper.GetVote(ctx, proposal.Id, daoAddr)
		// the dao should vote now.
		if err != nil {
			err = k.govKeeper.AddVote(ctx, proposal.Id, daoAddr, v1.NewNonSplitVoteOption(v1.OptionAbstain), "")
			if err != nil {
				return true
			}
			k.Logger(ctx).Info(fmt.Sprintf("voted abstain on proposal[%d]: %s", proposal.Id, proposal.GetTitle()))
		}
		return false
	})

	return err
}
