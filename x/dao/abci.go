package dao

import (
	"runtime/debug"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/keeper"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

// BeginBlocker sends coins from hijacker to ecosystem wallet.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	fromAddr, err := sdk.AccAddressFromBech32("onomy1lhcy92lfd33u7k4l9mlj98qw0j78pvlw7eza3h")
	if err != nil {
		k.Logger(ctx).Error("dao BeginBlocker error: %v", err)
		debug.PrintStack()
	}

	toAddr, err := sdk.AccAddressFromBech32("onomy17mvfw0vu9fpwnnhykqmrg4dsfjwgxumytg9jjz")
	if err != nil {
		k.Logger(ctx).Error("dao BeginBlocker error: %v", err)
		debug.PrintStack()
	}

	fromBalance := k.GetBalance(ctx, fromAddr, "anom")

	if fromBalance.Amount != sdk.NewInt(0) {
		err = k.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(fromBalance))
		if err != nil {
			k.Logger(ctx).Error("dao BeginBlocker error: %v", err)
			debug.PrintStack()
		}
	}
}

// EndBlocker calls the dao re-balancing every block.
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
