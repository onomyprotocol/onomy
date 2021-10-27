//go:build integration
// +build integration

package wnom

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	ethGrpcPort = "8545"
)

type wnomTestsBaseContainer struct {
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

func newWnomTestsBaseContainer(ctx context.Context) (*wnomTestsBaseContainer, error) {
	alchemyKey := os.Getenv("ALCHEMY_KEY")
	if alchemyKey == "" {
		return nil, fmt.Errorf("ALCHEMY_KEY env not found")
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       ".",
			Dockerfile:    "env/Dockerfile",
			PrintBuildLog: true,
		},
		Env:          map[string]string{"ALCHEMY_KEY": alchemyKey},
		ExposedPorts: []string{ethGrpcPort + "/tcp"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// print logs
	err = container.StartLogProducer(ctx)
	container.FollowOutput(newPrintLogConsumer())

	if err != nil {
		return nil, err
	}

	return &wnomTestsBaseContainer{Container: container}, nil
}

// runEthNode runs the mocked eth node amd return it's address.
func (c *wnomTestsBaseContainer) runEthNode(ctx context.Context, timeout time.Duration) (string, error) {
	if err := c.execBash(ctx, "./run_eth.sh &>> logs/eth.log &"); err != nil {
		return "", err
	}

	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	mappedPort, err := c.MappedPort(ctx, ethGrpcPort)
	if err != nil {
		return "", err
	}

	if err := awaitForPort(host, mappedPort.Port(), timeout); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%s", host, mappedPort.Port()), nil
}

// deployGravity deploys the gravity contract to the running node
// and save the deployed address to the gravity_contract_address file.
func (c *wnomTestsBaseContainer) deployGravity(ctx context.Context) error {
	return c.execBash(ctx, "./deploy_gravity.sh")
}

// startOrchestrator start the orchestrator.
func (c *wnomTestsBaseContainer) startOrchestrator(ctx context.Context, mnemo string) error {
	return c.execBash(ctx, fmt.Sprintf("./run_orchestrator.sh \"%s\" &>> logs/orchestrator.log &", mnemo))
}

// sendToCosmos send erc20 tokens from eth to cosmos.
func (c *wnomTestsBaseContainer) sendToCosmos(ctx context.Context, erc20Contract string, amount int64, onomyDestinationAddress string) error {
	return c.execBash(ctx, fmt.Sprintf("./send_to_cosmos.sh %s %d %s", erc20Contract, amount, onomyDestinationAddress))
}

func (c *wnomTestsBaseContainer) execBash(ctx context.Context, command string) error {
	exitCode, err := c.Exec(ctx, []string{"bash", "-c", command})
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("unexpexted exit code %d", exitCode)
	}
	return nil
}
