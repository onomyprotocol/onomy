package cli_test

import (
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/onomyprotocol/onomy/testutil/network"
	"github.com/onomyprotocol/onomy/x/dao/client/cli"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestCLI_CmdShowParams(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		prep func(testNetwork *network.TestNetwork)
		err  error
		want types.QueryParamsResponse
	}{
		{
			name: "positive_default",
			want: types.QueryParamsResponse{
				Params: types.DefaultParams(),
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			testNetwork := network.New(t)
			defer testNetwork.Cleanup()

			if tt.prep != nil {
				tt.prep(testNetwork)
			}

			tt.args = append(tt.args, fmt.Sprintf("--%s=json", tmcli.OutputFlag))
			ctx := testNetwork.Validator1Ctx()

			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowParams(), tt.args)
			if tt.err != nil {
				stat, ok := status.FromError(tt.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tt.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryParamsResponse
				require.NoError(t, testNetwork.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp)
				require.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestCLI_ShowTreasury(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		prep func(testNetwork *network.TestNetwork)
		err  error
		want types.QueryTreasuryResponse
	}{
		{
			name: "positive_empty",
			want: types.QueryTreasuryResponse{
				TreasuryBalance: []sdk.Coin{},
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			testNetwork := network.New(t)
			defer testNetwork.Cleanup()

			if tt.prep != nil {
				tt.prep(testNetwork)
			}

			tt.args = append(tt.args, fmt.Sprintf("--%s=json", tmcli.OutputFlag))
			ctx := testNetwork.Validator1Ctx()

			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowTreasury(), tt.args)
			if tt.err != nil {
				stat, ok := status.FromError(tt.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tt.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryTreasuryResponse
				require.NoError(t, testNetwork.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp)
				require.Equal(t, tt.want, resp)
			}
		})
	}
}
