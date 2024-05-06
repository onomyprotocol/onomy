package ics

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"

	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	appConsumer "github.com/cosmos/interchain-security/v3/app/consumer"
	"github.com/cosmos/interchain-security/v3/tests/integration"
	icstestingutils "github.com/cosmos/interchain-security/v3/testutil/ibc_testing"

	onomyApp "github.com/onomyprotocol/onomy-rebuild/app"
)

func TestCCVTestSuite(t *testing.T) {
	// Pass in concrete app types that implement the interfaces defined in https://github.com/cosmos/interchain-security/testutil/integration/interfaces.go
	// IMPORTANT: the concrete app types passed in as type parameters here must match the
	// concrete app types returned by the relevant app initers.
	ccvSuite := integration.NewCCVTestSuite[*onomyApp.OnomyApp, *appConsumer.App](
		// Pass in ibctesting.AppIniters for onomy (provider) and consumer.
		OnomyAppIniter, icstestingutils.ConsumerAppIniter, []string{})

	// Run tests
	suite.Run(t, ccvSuite)
}

// OnomyAppIniter implements ibctesting.AppIniter for the onomy app
func OnomyAppIniter() (ibctesting.TestingApp, map[string]json.RawMessage) {
	encoding := onomyApp.RegisterEncodingConfig()
	app := onomyApp.NewOnomyApp(
		log.NewNopLogger(),
		tmdb.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		onomyApp.DefaultNodeHome,
		encoding,
		onomyApp.EmptyAppOptions{})

	testApp := ibctesting.TestingApp(app)
	return testApp, onomyApp.NewDefaultGenesisState(encoding)
}
