package v2_1_0 //nolint:revive,stylecheck // app version

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onomyprotocol/onomy/app/keepers"
)

// UpgradeHandler is an x/upgrade handler.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {

	return func(c context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		collectionsConsensus(ctx, keepers)
		return vm, nil
	}
}

func collectionsConsensus(ctx sdk.Context, keepers *keepers.AppKeepers) {
	keepers.ConsensusParamsKeeper.ParamsStore.Set(ctx, ctx.ConsensusParams())
}
