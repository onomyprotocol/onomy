//go:build integration
// +build integration

package wnom

import (
	"context"
	"fmt"
	"testing"
	"time"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/onomyprotocol/onomy/testutil"
)

func TestIntegrationWnomToAnom(t *testing.T) { // nolint:gocyclo, cyclop
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	defer ctx.Done()
	const (
		bootstrappingTimeout    = time.Minute
		onomyDestinationAddress = "onomy1txg674n2km4ft6jfdccs94xtg8vl2kyksw3scl"
		fauTokeAddress          = "0xFab46E002BbF0b4509813474841E0716E6730136"
	)

	wnomTestsBaseContainer, err := newWnomTestsBaseContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	defer wnomTestsBaseContainer.terminate(ctx, t)

	// run ethereum node
	err = testutil.RetryWithTimeout(func() error {
		err := wnomTestsBaseContainer.runEthNode(ctx)
		if err != nil {
			t.Logf("run eth failed: %s, will be retried in %d", err.Error(), testutil.DefaultRetryTimeout)
		}
		return err
	}, bootstrappingTimeout)
	if err != nil {
		t.Fatal(err)
	}

	// run onomy chain
	onomyChain, err := testutil.NewOnomyChain()
	if err != nil {
		t.Fatal(err)
	}

	if err := onomyChain.Start(bootstrappingTimeout); err != nil {
		t.Fatal(err)
	}

	t.Logf("onomy chain is up, validator: %+v", onomyChain.Validator)

	// deploy gravity
	err = testutil.RetryWithTimeout(func() error {
		err := wnomTestsBaseContainer.deployGravity(ctx)
		if err != nil {
			t.Logf("deployGravity failed: %s, will be retried in %d", err.Error(), testutil.DefaultRetryTimeout)
		}
		return err
	}, bootstrappingTimeout)

	if err != nil {
		t.Fatal(err)
	}
	t.Log("gravity contract is deployed")

	// start orchestrator
	if err := wnomTestsBaseContainer.startOrchestrator(ctx, onomyChain.Validator.Mnemonic); err != nil {
		t.Fatal(err)
	}
	t.Log("orchestrator is started")

	// send wNOM tokens to onomy
	erc20Amount := int64(10)
	if err := wnomTestsBaseContainer.sendToCosmos(ctx, testutil.WnomERC20Address, erc20Amount, onomyDestinationAddress); err != nil {
		t.Fatal(err)
	}
	if err := wnomTestsBaseContainer.sendToCosmos(ctx, fauTokeAddress, erc20Amount, onomyDestinationAddress); err != nil {
		t.Fatal(err)
	}
	t.Log("ERC20 tokens are sent to the gravity contract")

	// waif for the wNOM on the validator balance
	err = testutil.RetryWithTimeout(func() error {
		balance, err := onomyChain.GetAccountBalance(onomyDestinationAddress)
		if err != nil {
			return err
		}

		checks := 0
		for _, coin := range balance {
			if coin.Denom == testutil.AnomDenom {
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

		err = fmt.Errorf("the %q hasn't received the %s tokens yet, balance: %+v", onomyDestinationAddress, testutil.WnomERC20Address, balance)
		t.Logf("%v, will be retried in %d", err, testutil.DefaultRetryTimeout)

		return err
	}, bootstrappingTimeout)

	if err != nil {
		t.Fatal(err)
	}
}
