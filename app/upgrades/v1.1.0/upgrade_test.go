package v1_1_0_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"

	v1_1_0 "github.com/onomyprotocol/onomy/app/upgrades/v1.1.0"
	"github.com/onomyprotocol/onomy/testutil/simapp"
)

func TestSetupUpgradeHandler(t *testing.T) {
	tenPercents := sdk.NewDec(1).Quo(sdk.NewDec(10))
	commission := stakingtypes.NewCommissionRates(tenPercents, tenPercents, tenPercents)

	// generate the validator with the eth gravity orchestrator
	vals := map[string]simapp.ValReq{
		"v1": {
			SelfBondCoin: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)),
			Commission:   commission,
			Balance:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(2000, sdk.DefaultPowerReduction))),
		},
		"v2": {
			SelfBondCoin: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)),
			Commission:   commission,
			Balance:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(2000, sdk.DefaultPowerReduction))),
		},
	}

	simApp, _ := simapp.SetupWithValidators(t, vals)
	ctx := simApp.NewContext()

	handler := v1_1_0.PrepareUpgradeHandler(simApp.OnomyApp().ArcEthGravityKeeper, simApp.OnomyApp().ArcBnbGravityKeeper)
	_, err := handler(ctx, upgradetypes.Plan{}, module.VersionMap{})
	require.NoError(t, err)

	params := simApp.OnomyApp().ArcEthGravityKeeper.GetParams(ctx)

	require.Equal(t, uint64(100_000), params.SignedValsetsWindow)
	require.Equal(t, uint64(100_000), params.SignedBatchesWindow)
	require.Equal(t, uint64(100_000), params.SignedLogicCallsWindow)
	require.Equal(t, uint64(100_000), params.UnbondSlashingValsetsWindow)
}
