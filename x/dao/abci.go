package dao

import (
	"runtime/debug"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/keeper"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

//nolint:all EndBlocker calls the dao re-balancing every block.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if k.GetDaoDelegationSupply(ctx).GT(sdk.NewDec(0)) {
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
