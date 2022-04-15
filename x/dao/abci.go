package dao

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/keeper"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

// EndBlocker calls the dao re-balancing every block.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if err := k.ReBalanceDelegation(ctx); err != nil {
		k.Logger(ctx).Error("re-balance delegation error handled in %s EndBlocker , %v", types.ModuleName, err)
	}
}
