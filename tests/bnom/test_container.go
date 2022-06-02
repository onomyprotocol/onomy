//go:build integration
// +build integration

package bnom

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

const (
	alchemyKey = "ALCHEMY_KEY"
)

type bnomTestsBaseContainer struct {
	testcontainers.Container
}

type printLogConsumer struct{}

func newPrintLogConsumer() *printLogConsumer {
	return &printLogConsumer{}
}

func (c *printLogConsumer) Accept(logline testcontainers.Log) {
	logText := strings.TrimSpace(string(logline.Content))
	if logText != "" {
		log.Println(logText)
	}
}

func newBnomTestsBaseContainer(ctx context.Context) (*bnomTestsBaseContainer, error) {
	alchemyKeyValue := os.Getenv(alchemyKey)
	if alchemyKeyValue == "" {
		return nil, fmt.Errorf("%s env not found", alchemyKey)
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       ".",
			Dockerfile:    "env/Dockerfile",
			PrintBuildLog: true,
		},
		Env:         map[string]string{alchemyKey: alchemyKeyValue},
		AutoRemove:  true,
		NetworkMode: "host",
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// print logs
	container.FollowOutput(newPrintLogConsumer())
	err = container.StartLogProducer(ctx)
	if err != nil {
		return nil, err
	}

	return &bnomTestsBaseContainer{Container: container}, nil
}

// runEthNode runs the mocked eth node amd return it's address.
func (c *bnomTestsBaseContainer) runEthNode(ctx context.Context) error {
	return c.execBash(ctx, "./run_eth.sh &>> logs/eth.log &")
}

// deployGravity deploys the gravity contract to the running node
// and save the deployed address to the gravity_contract_address file.
func (c *bnomTestsBaseContainer) deployGravity(ctx context.Context) error {
	return c.execBash(ctx, "./deploy_gravity.sh &>> logs/orchestrator.log")
}

// startOrchestrator start the orchestrator.
func (c *bnomTestsBaseContainer) startOrchestrator(ctx context.Context, mnemo string) error {
	return c.execBash(ctx, fmt.Sprintf("./run_orchestrator.sh \"%s\" &>> logs/orchestrator.log &", mnemo))
}

// sendToCosmos send erc20 tokens from eth to cosmos.
func (c *bnomTestsBaseContainer) sendToCosmos(ctx context.Context, erc20Contract string, amount int64, onomyDestinationAddress string) error {
	return c.execBash(ctx, fmt.Sprintf("./send_to_cosmos.sh %s %d %s", erc20Contract, amount, onomyDestinationAddress))
}

func (c *bnomTestsBaseContainer) execBash(ctx context.Context, command string) error {
	exitCode, err := c.Exec(ctx, []string{"bash", "-c", command})
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("unexpexted exit code %d, %w", exitCode, err)
	}
	return nil
}

func (c *bnomTestsBaseContainer) terminate(ctx context.Context, t *testing.T) {
	t.Helper()

	t.Logf("container logs:")
	logs, err := c.logs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(logs)
	c.Container.Terminate(ctx) // nolint:errcheck
}

func (c *bnomTestsBaseContainer) logs(ctx context.Context) (string, error) {
	readCloser, err := c.Container.Logs(ctx)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(readCloser)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
