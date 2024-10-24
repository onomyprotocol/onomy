// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v2_1_1 //nolint:revive,stylecheck // app version

import (
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/onomyprotocol/onomy/app/keepers"
)

// Name is migration name.
const Name = "v2.1.1"

// UpgradeHandler is an x/upgrade handler.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {

	return func(c context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return vm, nil
	}
}
