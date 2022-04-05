package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/onomyprotocol/onomy/testutil/network"
	"github.com/onomyprotocol/onomy/x/dao/client/cli"
)

const (
	normalToken = "node0token"
	stakeToken  = "stake"
)

func TestCLI_FundTreasuryProposal(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		err  error
		code uint32
	}{
		{
			name: "positive",
			args: fmt.Sprintf("5000000000000000000%s --title=Title --description=Description --deposit=1000000000000000000%s", normalToken, stakeToken),
		},
		{
			name: "negative_insufficient_balance",
			args: fmt.Sprintf("5000000000000000000%s --title=Title --description=Description --deposit=1000000000000000000%s", "newcoin", stakeToken),
			code: govtypes.ErrInvalidProposalContent.ABCICode(),
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			testNetwork := network.New(t)
			defer testNetwork.Cleanup()
			args := strings.Split(tt.args, " ")
			args = append(args, testNetwork.TxValidator1Args()...)
			ctx := testNetwork.Validator1Ctx()
			// on the onomy chain those flag will be added ny the gov module
			cmd := cli.CmdFundTreasuryProposal()
			flags.AddTxFlagsToCmd(cmd)
			out, err := clitestutil.ExecTestCLICmd(ctx, cmd, args)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, testNetwork.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tt.code, resp.Code, fmt.Sprintf("%+v", resp))
			}
		})
	}
}
