// Package testutil contains testutils.
package testutil

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
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
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/app"
	onomydCmd "github.com/onomyprotocol/onomy/cmd/onomyd/cmd"
)

const (
	// AnomDenom is anom name .
	AnomDenom = "anom"
	// WnomERC20Address is wnom eth address .
	WnomERC20Address = "0xe7c0fd1f0A3f600C1799CD8d335D31efBE90592C"

	// TestChainName is default test chain name.
	TestChainName = "chain-1"
	// TestChainFlag is default chain flag.
	TestChainFlag = "--chain-id=" + TestChainName
	// TestChainKeyRingFlag is default keyring flag.
	TestChainKeyRingFlag = "--keyring-backend=test"
	jsonOutFlag          = "--output=json"
	// TestChainDenom is default chain denom.
	TestChainDenom = AnomDenom
	// TestChainValidatorGenesysAmount is default validator genesys amount.
	TestChainValidatorGenesysAmount = "10000000000000" + TestChainDenom
	// TestChainValidatorStakeAmount is default validator genesys stake amount.
	TestChainValidatorStakeAmount = "1000000000000" + TestChainDenom
	// TestChainValidator1Name is default validator name.
	TestChainValidator1Name = "validator1"
	// TestChainValidator1EthAddress is default validator eth pub key.
	TestChainValidator1EthAddress = "0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d"

	// OnomyGrpcHost is default host.
	OnomyGrpcHost = "127.0.0.1"
	// OnomyGrpcPort is default port.
	OnomyGrpcPort = "9090"
)

// OnomyChain is test struct for the chain running.
type OnomyChain struct {
	homeFlag  string
	Validator keyring.KeyOutput
}

// NewOnomyChain creates a new OnomyChain.
func NewOnomyChain() (*OnomyChain, error) {
	// prepare test folder for genesys and data files
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	homeFlag := "--home=" + dir

	// generate genesys
	ExecuteChainCmd("init", TestChainName, TestChainFlag, homeFlag)

	// enable swagger and rest API:
	if err := replaceStringInFile(filepath.Join(dir, "config", "app.toml"), "enable = false", "enable = true"); err != nil {
		return nil, err
	}

	if err := replaceStringInFile(filepath.Join(dir, "config", "app.toml"), "swagger = false", "swagger = true"); err != nil {
		return nil, err
	}

	if err := replaceStringInFile(filepath.Join(dir, "config", "genesis.json"), "\"stake\"", "\""+TestChainDenom+"\""); err != nil {
		return nil, err
	}

	// set up swap parameters
	if err := replaceGenesysSettings(filepath.Join(dir, "config", "genesis.json"), "app_state.gravity.params.erc20_to_denom_permanent_swap",
		json.RawMessage(fmt.Sprintf(`{"erc20": "%s", "denom": "%s"}`, WnomERC20Address, AnomDenom))); err != nil {
		return nil, err
	}

	if err := replaceStringInFile(filepath.Join(dir, "config", "config.toml"), "log_level = \"info\"", "log_level = \"error\""); err != nil {
		return nil, err
	}

	// add new user
	val1KeyString := ExecuteChainCmd("keys add", TestChainValidator1Name, TestChainKeyRingFlag, jsonOutFlag, homeFlag)
	var val1KeyOutput keyring.KeyOutput
	if err := json.Unmarshal([]byte(val1KeyString), &val1KeyOutput); err != nil {
		return nil, err
	}

	// add user to genesys
	ExecuteChainCmd("add-genesis-account", val1KeyOutput.Address, TestChainValidatorGenesysAmount, homeFlag)

	// gentx
	ExecuteChainCmd("gentx", TestChainValidator1Name, TestChainValidatorStakeAmount, TestChainValidator1EthAddress, val1KeyOutput.Address, TestChainFlag, TestChainKeyRingFlag, homeFlag)

	// collect gentx
	ExecuteChainCmd("collect-gentxs", homeFlag)

	return &OnomyChain{
		homeFlag:  homeFlag,
		Validator: val1KeyOutput,
	}, nil
}

// Start start the OnomyChain.
func (oc *OnomyChain) Start(timeout time.Duration) error {
	go ExecuteChainCmd("start", oc.homeFlag)

	// wait for grpc port
	return AwaitForPort(OnomyGrpcHost, OnomyGrpcPort, timeout)
}

// Stop stops the OnomyChain.
func (oc *OnomyChain) Stop() {
	ExecuteChainCmd("stop", oc.homeFlag)
}

// GetAccountBalance return the 'address' balance.
func (oc *OnomyChain) GetAccountBalance(address string) ([]sdkTypes.Coin, error) {
	balanceString := ExecuteChainCmd("query bank balances", address, oc.homeFlag, jsonOutFlag)
	var balances struct {
		Balances []sdkTypes.Coin `json:"balances"`
	}
	if err := json.Unmarshal([]byte(balanceString), &balances); err != nil {
		return nil, err
	}
	return balances.Balances, nil
}

// ExecuteChainCmd executes any cmd on the onomyd cli.
func ExecuteChainCmd(cmd string, args ...string) string {
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
	rootCmd, _ := onomydCmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
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

func replaceGenesysSettings(filePath, settingPath string, newValue json.RawMessage) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var fileRawJSON map[string]json.RawMessage
	if err := json.Unmarshal(input, &fileRawJSON); err != nil {
		return err
	}

	if err := replaceJSONInJSONmap(fileRawJSON, strings.Split(settingPath, "."), newValue); err != nil {
		return err
	}

	output, err := json.Marshal(fileRawJSON)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, output, 0666) // nolint:gomnd
}

func replaceJSONInJSONmap(object map[string]json.RawMessage, settingPath []string, newValue json.RawMessage) error {
	if len(settingPath) == 0 {
		return nil
	}
	for key := range object {
		if key == settingPath[0] && len(settingPath) == 1 {
			object[key] = newValue
			return nil
		}

		var nextRawJSON map[string]json.RawMessage
		if err := json.Unmarshal(object[key], &nextRawJSON); err != nil {
			// not object
			continue
		}

		if err := replaceJSONInJSONmap(nextRawJSON, settingPath[1:], newValue); err != nil {
			return err
		}

		nextRawJSONBytes, err := json.Marshal(nextRawJSON)
		if err != nil {
			return err
		}
		object[key] = nextRawJSONBytes
	}
	return nil
}
