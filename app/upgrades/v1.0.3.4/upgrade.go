// Package v1_0_3_4 is contains chain upgrade of the corresponding version.
package v1_0_3_4 //nolint:revive,stylecheck // app version

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Name is migration name.
const Name = "v1.0.3.4"

// UpgradeHandler is an x/upgrade handler.
func UpgradeHandler(_ context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return vm, nil
}
