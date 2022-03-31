// Package cli contains test cli utils.
package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

// QueryArgs returns cli query default args used for tests.
func QueryArgs() []string {
	return []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
}

// ExecTestCLICmd builds the client context, mocks the output and executes the command.
func ExecTestCLICmd(clientCtx client.Context, cmd *cobra.Command, extraArgs []string) (testutil.BufferWriter, error) {
	return clitestutil.ExecTestCLICmd(clientCtx, cmd, extraArgs)
}
