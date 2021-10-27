//go:build integration
// +build integration

// Package wnom contains integration test for the wnom to nom swapping use case.
package wnom

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/cmd/onomyd/cmd"
)

const (
	chainName              = "chain-1"
	chainFlag              = "--chain-id=" + chainName
	keyRingFlag            = "--keyring-backend=test"
	jsonOutFlag            = "--output=json"
	chainDenom             = "anom"
	validatorGenesysAmount = "10000000000000" + chainDenom
	validatorStakeAmount   = "1000000000000" + chainDenom
	validator1Name         = "validator1"
	validator1EthAddress   = "0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d"

	onomyGrpcHost = "127.0.0.1"
	onomyGrpcPort = "9090"
)

type onomyChain struct {
	homeFlag  string
	Validator keyring.KeyOutput
}

func newOnomyChain() (*onomyChain, error) {
	// prepare test folder for genesys and data files
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	homeFlag := "--home=" + dir

	// generate genesys
	executeChainCmd("init", chainName, chainFlag, homeFlag)

	// enable swagger and rest API:
	if err := replaceStringInFile(filepath.Join(dir, "config", "app.toml"), "enable = false", "enable = true"); err != nil {
		return nil, err
	}

	if err := replaceStringInFile(filepath.Join(dir, "config", "app.toml"), "swagger = false", "swagger = true"); err != nil {
		return nil, err
	}

	if err := replaceStringInFile(filepath.Join(dir, "config", "genesis.json"), "\"stake\"", "\""+chainDenom+"\""); err != nil {
		return nil, err
	}

	// add new user
	val1KeyString := executeChainCmd("keys add", validator1Name, keyRingFlag, jsonOutFlag, homeFlag)
	var val1KeyOutput keyring.KeyOutput
	if err := json.Unmarshal([]byte(val1KeyString), &val1KeyOutput); err != nil {
		return nil, err
	}

	// add user to genesys
	executeChainCmd("add-genesis-account", val1KeyOutput.Address, validatorGenesysAmount, homeFlag)

	// gentx
	executeChainCmd("gentx", validator1Name, validatorStakeAmount, validator1EthAddress, val1KeyOutput.Address, chainFlag, keyRingFlag, homeFlag)

	// collect gentx
	executeChainCmd("collect-gentxs", homeFlag)

	return &onomyChain{
		homeFlag:  homeFlag,
		Validator: val1KeyOutput,
	}, nil
}

func (oc *onomyChain) start(timeout time.Duration) error {
	go executeChainCmd("start", oc.homeFlag)

	// wait for grpc port
	return awaitForPort(onomyGrpcHost, onomyGrpcPort, timeout)
}

func (oc *onomyChain) getAccountBalance(address string) ([]sdkTypes.Coin, error) {
	balanceString := executeChainCmd("query bank balances", address, oc.homeFlag, jsonOutFlag)
	var balances struct {
		Balances []sdkTypes.Coin `json:"balances"`
	}
	if err := json.Unmarshal([]byte(balanceString), &balances); err != nil {
		return nil, err
	}
	return balances.Balances, nil
}

func executeChainCmd(cmd string, args ...string) string {
	oldArgs := os.Args
	// this call is required because otherwise flags panics, if args are set between flag.Parse calls
	flag.CommandLine = flag.NewFlagSet("command", flag.ExitOnError)
	// we need a value to set Args[0] to, cause flag begins parsing at Args[1]
	args = append(strings.Fields(cmd), args...)
	os.Args = append([]string{"onomyd"}, args...)

	// rests config seal protection
	config := sdkTypes.GetConfig()
	setField(config, "sealed", false)
	setField(config, "sealedch", make(chan struct{}))

	out := captureOutput(func() {
		mainTestRunner()
	})

	os.Args = oldArgs

	return out
}

func mainTestRunner() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		switch e := err.(type) { // nolint:errorlint
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()
	os.Stdout = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader) // nolint: errcheck
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close() // nolint: errcheck
	return <-out
}

func setField(object interface{}, fieldName string, value interface{}) {
	rs := reflect.ValueOf(object).Elem()
	field := rs.FieldByName(fieldName)
	// rf can't be read or set.
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

func replaceStringInFile(filePath, from, to string) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	output := bytes.ReplaceAll(input, []byte(from), []byte(to))
	return ioutil.WriteFile(filePath, output, 0666) // nolint:gomnd
}
