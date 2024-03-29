// Package v1_0_1 is contains chain upgrade of the corresponding version.
package v1_0_1 //nolint:revive,stylecheck // app version

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// Name is migration name.
const Name = "v1.0.1"

// UpgradeHandler is an x/upgrade handler.
func UpgradeHandler(_ sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return vm, nil
}
