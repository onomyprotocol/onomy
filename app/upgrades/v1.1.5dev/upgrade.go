// Package v1_1_5dev is contains chain upgrade of the corresponding version.
package v1_1_5dev //nolint:revive,stylecheck // app version

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// Name is migration name.
const Name = "v1.1.5dev"

// UpgradeHandler is an x/upgrade handler.
func UpgradeHandler(_ sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return vm, nil
}
