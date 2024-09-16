package dao

import (
	"context"
	"runtime/debug"
	// "time"

	"cosmossdk.io/math"
	// "github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/onomyprotocol/onomy/x/dao/keeper"
	// "github.com/onomyprotocol/onomy/x/dao/types"
)

// BeginBlocker does any custom logic for the DAO upon `BeginBlocker`
func BeginBlocker(ctx context.Context, k keeper.Keeper) {
}

// EndBlocker calls the dao re-balancing every block.
func EndBlocker(ctx context.Context, k keeper.Keeper) {
	// defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if k.GetDaoDelegationSupply(ctx).GT(math.LegacyZeroDec()) {
		if err := k.VoteAbstain(ctx); err != nil {
			k.Logger(ctx).Error("dao EndBlocker error: %v", err)
			debug.PrintStack()
		}

		if err := k.WithdrawReward(ctx); err != nil {
			k.Logger(ctx).Error("dao EndBlocker error: %v", err)
			debug.PrintStack()
		}

		if err := k.UndelegateAllValidators(ctx); err != nil {
			k.Logger(ctx).Error("dao EndBlocker error: %v", err)
			debug.PrintStack()
		}
	}

	if err := k.InflateDao(ctx); err != nil {
		k.Logger(ctx).Error("dao EndBlocker error: %v", err)
		debug.PrintStack()
	}
}
