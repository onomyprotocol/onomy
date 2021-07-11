package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/onomyprotocol/onomy/x/market/types"
)

var _ = strconv.Itoa(0)

func CmdCancelOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancelOrder [port] [channel] [amountDenom] [exchRateDenom] [orderID]",
		Short: "Cancel an order",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsPort := string(args[0])
			argsChannel := string(args[1])
			argsAmountDenom := string(args[2])
			argsExchRateDenom := string(args[3])
			argsOrderID, _ := strconv.ParseInt(args[4], 10, 64)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelOrder(clientCtx.GetFromAddress().String(), string(argsPort), string(argsChannel), string(argsAmountDenom), string(argsExchRateDenom), int32(argsOrderID))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
