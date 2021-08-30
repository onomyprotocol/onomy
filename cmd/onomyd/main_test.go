package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"unsafe"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInitToCollectGentxsFlow(t *testing.T) {
	const (
		chainName            = "chain-1"
		chainFlag            = "--chain-id=" + chainName
		keyRingFlag          = "--keyring-backend=test"
		jsonOutFlag          = "--output=json"
		stakeAmount          = "1000000000000stake"
		validator1Name       = "validator1"
		validator1EthAddress = "0x033030FEeBd93E3178487c35A9c8cA80874353C9"
	)

	// prepare test folder for genesys and data files
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	homeFlag := "--home=" + dir
	defer os.RemoveAll(dir) // nolint: errcheck

	// generate genesys
	executeCmd(t, "init:", []string{"init", chainName, chainFlag, homeFlag})

	// add new user
	val1KeyString := executeCmd(t, "add key:", []string{"keys", "add", validator1Name, keyRingFlag, jsonOutFlag, homeFlag})
	t.Log(val1KeyString)
	var val1KeyOutput keyring.KeyOutput
	if err := json.Unmarshal([]byte(val1KeyString), &val1KeyOutput); err != nil {
		t.Fatal(err)
	}

	// add user to genesys
	executeCmd(t, "add validator1 to genesys:", []string{"add-genesis-account", val1KeyOutput.Address, stakeAmount, homeFlag})

	// gentx
	executeCmd(t, "gentx:", []string{"gentx", validator1Name, stakeAmount, validator1EthAddress, val1KeyOutput.Address, chainFlag, keyRingFlag, homeFlag})

	// collect gentx
	collectGentxsOut := executeCmd(t, "collect-gentxs:", []string{"collect-gentxs", homeFlag})
	t.Log(collectGentxsOut)

	// eth keys
	ethKey := executeCmd(t, "gentx:", []string{"eth_keys", "add", homeFlag, "--output=json"})
	t.Log(ethKey)
}

func executeCmd(t *testing.T, name string, args []string) string {
	t.Helper()
	oldArgs := os.Args
	// this call is required because otherwise flags panics, if args are set between flag.Parse calls
	flag.CommandLine = flag.NewFlagSet(name, flag.ExitOnError)
	// we need a value to set Args[0] to, cause flag begins parsing at Args[1]
	os.Args = append([]string{"onomyd"}, args...)
	t.Log(strings.Join(os.Args, " "))

	// rests config seal protection
	config := sdk.GetConfig()
	setField(config, "sealed", false)
	setField(config, "sealedch", make(chan struct{}))

	out := captureOutput(func() {
		main()
	})

	os.Args = oldArgs

	return out
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
