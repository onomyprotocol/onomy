package onomy_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	db "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	onomy "github.com/onomyprotocol/onomy-rebuild/app"
	onomyhelpers "github.com/onomyprotocol/onomy-rebuild/app/helpers"
)

type EmptyAppOptions struct{}

func (ao EmptyAppOptions) Get(_ string) interface{} {
	return nil
}

func TestOnomyApp_BlockedModuleAccountAddrs(t *testing.T) {
	encConfig := onomy.RegisterEncodingConfig()
	app := onomy.NewOnomyApp(
		log.NewNopLogger(),
		db.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		onomy.DefaultNodeHome,
		encConfig,
		EmptyAppOptions{},
	)

	moduleAccountAddresses := app.ModuleAccountAddrs()
	blockedAddrs := app.BlockedModuleAccountAddrs(moduleAccountAddresses)

	require.NotContains(t, blockedAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
}

func TestOnomyApp_Export(t *testing.T) {
	app := onomyhelpers.Setup(t)
	_, err := app.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
