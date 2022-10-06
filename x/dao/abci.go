package dao

import (
	"runtime/debug"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/keeper"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

// EndBlocker calls the dao re-balancing every block.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if err := endBlocker(ctx, k); err != nil {
		k.Logger(ctx).Error("dao EndBlocker error: %v", err)
		debug.PrintStack()
	}
}

func endBlocker(ctx sdk.Context, k keeper.Keeper) error {
	if ctx.BlockHeight()%k.WithdrawRewardPeriod(ctx) == 0 {
		if err := k.WithdrawReward(ctx); err != nil {
			return err
		}
	}

	if err := k.ReBalanceDelegation(ctx); err != nil {
		return err
	}

	return k.VoteAbstain(ctx)
}
