// Package v1_0_1 is contains chain upgrade of the corresponding version.
package v1_0_1 //nolint:revive,stylecheck // app version

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Name is migration name.
const Name = "v1.0.1"

// UpgradeHandler is an x/upgrade handler.
func UpgradeHandler(_ context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return vm, nil
}

// func(ctx context.Context, plan Plan, fromVM module.VersionMap) (module.VersionMap, error).
