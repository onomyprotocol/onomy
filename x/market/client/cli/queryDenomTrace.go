package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/onomyprotocol/onomy/x/market/types"
	"github.com/spf13/cobra"
)

func CmdListDenomTrace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-denomTrace",
		Short: "list all denomTrace",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDenomTraceRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DenomTraceAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowDenomTrace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-denomTrace [index]",
		Short: "shows a denomTrace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetDenomTraceRequest{
				Index: args[0],
			}

			res, err := queryClient.DenomTrace(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
