// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v2_1_0 //nolint:revive,stylecheck // app version

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/onomyprotocol/onomy/app/keepers"
)

// Name is migration name.
const Name = "v2.1.0"

// UpgradeHandler is an x/upgrade handler.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {

	return func(c context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		baseAppLegacySS := keepers.ParamsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, &keepers.ConsensusParamsKeeper.ParamsStore)

		// onomy1vwr8z00ty7mqnk4dtchr9mn9j96nuh6wrlww93
		addrTreasury := sdk.MustAccAddressFromBech32("onomy1vwr8z00ty7mqnk4dtchr9mn9j96nuh6wrlww93")
		treasuryTokens := keepers.BankKeeper.GetAllBalances(ctx, addrTreasury)
		err := keepers.DistrKeeper.FundCommunityPool(ctx, treasuryTokens, addrTreasury)
		if err != nil {
			return vm, err
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
