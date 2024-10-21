// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v2_0_1 //nolint:revive,stylecheck // app version

import (
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/onomyprotocol/onomy/app/keepers"
)

// Name is migration name.
const Name = "v2.0.1"

// UpgradeHandler is an x/upgrade handler.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {

	return func(c context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		keepers.ConsensusParamsKeeper.ParamsStore.Set(ctx, ctx.ConsensusParams())

		// onomy1vwr8z00ty7mqnk4dtchr9mn9j96nuh6wrlww93
		addrTreasury := sdk.MustAccAddressFromBech32("onomy1vwr8z00ty7mqnk4dtchr9mn9j96nuh6wrlww93")
		treasuryTokens := keepers.BankKeeper.GetAllBalances(ctx, addrTreasury)
		err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addrTreasury, govtypes.ModuleName, treasuryTokens)
		if err != nil {
			return vm, err
		}

		return vm, nil
	}
}
