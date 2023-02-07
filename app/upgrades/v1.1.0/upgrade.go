// Package v1_1_0 is contains chain upgrade of the corresponding version.
package v1_1_0 //nolint:revive,stylecheck // app version

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	arcbnbkeeper "github.com/onomyprotocol/arc/module/bnb/x/gravity/keeper"
	arcbnbtypes "github.com/onomyprotocol/arc/module/bnb/x/gravity/types"
	arcethkeeper "github.com/onomyprotocol/cosmos-gravity-bridge/module/x/gravity/keeper"
)

// Name is migration name.
const Name = "v1.1.0"

// PrepareUpgradeHandler prepares the v1.1.0 upgrade.
func PrepareUpgradeHandler(arcEthGravityKeeper arcethkeeper.Keeper, arcBnbGravityKeeper arcbnbkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		delegateKeys := arcEthGravityKeeper.GetDelegateKeys(ctx)

		// use the same keys as for the gravity for the arc bnb
		for _, key := range delegateKeys {
			val, err := sdk.ValAddressFromBech32(key.Validator)
			if err != nil {
				return nil, err
			}
			orch, err := sdk.AccAddressFromBech32(key.Orchestrator)
			if err != nil {
				return nil, err
			}
			addr, err := arcbnbtypes.NewEthAddress(key.EthAddress)
			if err != nil {
				return nil, err
			}
			arcBnbGravityKeeper.SetOrchestratorValidator(ctx, val, orch)
			arcBnbGravityKeeper.SetEthAddressForValidator(ctx, val, *addr)
		}

		params := arcEthGravityKeeper.GetParams(ctx)
		// set all windows to a huge number in order to give the validators more time to set up the orchestrators
		// those params will be later changed via gov
		params.SignedValsetsWindow = 100_000
		params.SignedBatchesWindow = 100_000
		params.SignedLogicCallsWindow = 100_000
		params.UnbondSlashingValsetsWindow = 100_000

		arcEthGravityKeeper.SetParams(ctx, params)

		return vm, nil
	}
}
