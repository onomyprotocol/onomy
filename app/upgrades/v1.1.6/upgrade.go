package v1_1_6 //nolint:revive,stylecheck // app version

import (
	"context"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// UpgradeHandler is an x/upgrade handler.
func UpgradeHandler(_ context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return vm, nil
}
