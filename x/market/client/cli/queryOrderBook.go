package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/onomyprotocol/onomy/x/market/types"
	"github.com/spf13/cobra"
)

func CmdListOrderBook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-orderBook",
		Short: "list all orderBook",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllOrderBookRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.OrderBookAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowOrderBook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-orderBook [index]",
		Short: "shows a orderBook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetOrderBookRequest{
				Index: args[0],
			}

			res, err := queryClient.OrderBook(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
