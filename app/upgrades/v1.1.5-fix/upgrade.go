// Package v1_1_5 is contains chain upgrade of the corresponding version.
package v1_1_5_fix //nolint:revive,stylecheck // app version

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcproviderkeeper "github.com/cosmos/interchain-security/x/ccv/provider/keeper"
)

// Name is migration name.
const Name = "v1.1.5-fix"

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak *authkeeper.AccountKeeper,
	sk *stakingkeeper.Keeper,
	pk *ibcproviderkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
