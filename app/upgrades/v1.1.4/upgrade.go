// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v1_1_4 //nolint:revive,stylecheck // app version

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Name is migration name.
const Name = "v1.1.4"

// UpgradeHandler is an x/upgrade handler.
func UpgradeHandler(_ sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return vm, nil
}
