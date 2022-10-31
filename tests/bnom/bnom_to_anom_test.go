//go:build integration
// +build integration

package bnom

import (
	"context"
	"fmt"
	"testing"
	"time"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/onomyprotocol/onomy/testutil/integration"
	"github.com/onomyprotocol/onomy/testutil/retry"
)

func TestIntegrationBnomToAnom(t *testing.T) { // nolint:gocyclo, cyclop
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	defer ctx.Done()
	const (
		bootstrappingTimeout    = time.Minute
		onomyDestinationAddress = "onomy1txg674n2km4ft6jfdccs94xtg8vl2kyksw3scl"
		// https://erc20faucet.com/
		fauTokeAddress = "0xBA62BCfcAaFc6622853cca2BE6Ac7d845BC0f2Dc"
	)

	bnomTestsBaseContainer, err := newBnomTestsBaseContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	defer bnomTestsBaseContainer.terminate(ctx, t)

	// run ethereum node
	err = retry.WithTimeout(func() error {
		err := bnomTestsBaseContainer.runEthNode(ctx)
		if err != nil {
			t.Logf("run eth failed: %s, will be retried in %d", err.Error(), retry.DefaultRetryTimeout)
		}
		return err
	}, bootstrappingTimeout)
	if err != nil {
		t.Fatal(err)
	}

	// run onomy chain
	onomyChain, err := integration.NewOnomyChain()
	if err != nil {
		t.Fatal(err)
	}

	if err := onomyChain.Start(bootstrappingTimeout); err != nil {
		t.Fatal(err)
	}

	t.Logf("onomy chain is up, validator: %+v", onomyChain.Validator)

	// deploy gravity
	err = retry.WithTimeout(func() error {
		err := bnomTestsBaseContainer.deployGravity(ctx)
		if err != nil {
			t.Logf("deployGravity failed: %s, will be retried in %d", err.Error(), retry.DefaultRetryTimeout)
		}
		return err
	}, bootstrappingTimeout)

	if err != nil {
		t.Fatal(err)
	}
	t.Log("gravity contract is deployed")

	// start orchestrator
	if err := bnomTestsBaseContainer.startOrchestrator(ctx, onomyChain.Validator.Mnemonic); err != nil {
		t.Fatal(err)
	}
	t.Log("orchestrator is started")

	// send bNOM tokens to onomy
	erc20Amount := int64(10)
	if err := bnomTestsBaseContainer.sendToCosmos(ctx, integration.BnomERC20Address, erc20Amount, onomyDestinationAddress); err != nil {
		t.Fatal(err)
	}
	if err := bnomTestsBaseContainer.sendToCosmos(ctx, fauTokeAddress, erc20Amount, onomyDestinationAddress); err != nil {
		t.Fatal(err)
	}
	t.Log("ERC20 tokens are sent to the gravity contract")

	// waif for the bNOM on the validator balance
	err = retry.WithTimeout(func() error {
		balance, err := onomyChain.GetAccountBalance(onomyDestinationAddress)
		if err != nil {
			return err
		}

		checks := 0
		for _, coin := range balance {
			if coin.Denom == integration.AnomDenom {
				assert.Equal(t, coin.Amount, sdkTypes.NewIntWithDecimal(erc20Amount, 18))
				checks++
			}
			if coin.Denom == fmt.Sprintf("gravity%s", fauTokeAddress) {
				assert.Equal(t, coin.Amount, sdkTypes.NewIntWithDecimal(erc20Amount, 18))
				checks++
			}
		}
		if checks == 2 {
			t.Log(fmt.Sprintf("%q reveived test tokens", onomyDestinationAddress))
			return nil
		}

		err = fmt.Errorf("the %q hasn't received the %s tokens yet, balance: %+v", onomyDestinationAddress, integration.BnomERC20Address, balance)
		t.Logf("%v, will be retried in %d", err, retry.DefaultRetryTimeout)

		return err
	}, bootstrappingTimeout)

	if err != nil {
		t.Fatal(err)
	}
}
