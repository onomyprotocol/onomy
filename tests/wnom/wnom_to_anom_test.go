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
	err = retryWithTimeout(func() error {
		err := wnomTestsBaseContainer.runEthNode(ctx)
		if err != nil {
			t.Logf("run eth failed: %s, will be retried in %d", err.Error(), defaultRetryTimeout)
		}
		return err
	}, bootstrappingTimeout)
	if err != nil {
		t.Fatal(err)
	}

	// run onomy chain
	onomyChain, err := newOnomyChain()
	if err != nil {
		t.Fatal(err)
	}

	if err := onomyChain.start(bootstrappingTimeout); err != nil {
		t.Fatal(err)
	}

	t.Logf("onomy chain is up, validator: %+v", onomyChain.Validator)

	// deploy gravity
	err = retryWithTimeout(func() error {
		err := wnomTestsBaseContainer.deployGravity(ctx)
		if err != nil {
			t.Logf("deployGravity failed: %s, will be retried in %d", err.Error(), defaultRetryTimeout)
		}
		return err
	}, bootstrappingTimeout)

	if err != nil {
		t.Fatal(err)
	}
	t.Log("gravity contract deployed")

	// start orchestrator
	if err := wnomTestsBaseContainer.startOrchestrator(ctx, onomyChain.Validator.Mnemonic); err != nil {
		t.Fatal(err)
	}
	t.Log("orchestrator is started")

	// send wNOM tokens to onomy
	erc20Amount := int64(10)
	if err := wnomTestsBaseContainer.sendToCosmos(ctx, wnomERC20Address, erc20Amount, onomyDestinationAddress); err != nil {
		t.Fatal(err)
	}
	if err := wnomTestsBaseContainer.sendToCosmos(ctx, fauTokeAddress, erc20Amount, onomyDestinationAddress); err != nil {
		t.Fatal(err)
	}
	t.Log("ERC20 tokens are sent to the gravity contract")

	// waif for the wNOM on the validator balance
	err = retryWithTimeout(func() error {
		balance, err := onomyChain.getAccountBalance(onomyDestinationAddress)
		if err != nil {
			return err
		}

		checks := 0
		for _, coin := range balance {
			if coin.Denom == anomDenom {
				assert.Equal(t, coin.Amount, sdkTypes.NewIntWithDecimal(erc20Amount, 18))
				checks++
			}
			if coin.Denom == fmt.Sprintf("gravity%s", fauTokeAddress) {
				assert.Equal(t, coin.Amount, sdkTypes.NewIntWithDecimal(erc20Amount, 18))
				checks++
			}
		}
		if checks == 2 {
			return nil
		}

		err = fmt.Errorf("the node hasn't received the %s tokens yet, balance: %+v", wnomERC20Address, balance)
		t.Logf("%v, will be retried in %d", err, defaultRetryTimeout)

		return err
	}, bootstrappingTimeout)

	if err != nil {
		t.Fatal(err)
	}
}
