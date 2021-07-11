package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	channelutils "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/client/utils"
	"github.com/onomyprotocol/onomy/x/market/types"
)

var _ = strconv.Itoa(0)

func CmdSendCreateOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-createOrder [src-port] [src-channel] [amountDenom] [amount] [sourceCoin] [targetCoin] [exchRateDenom] [exchRate]",
		Short: "Send a createOrder over IBC",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsAmountDenom := string(args[2])
			argsAmount, _ := strconv.ParseInt(args[3], 10, 64)
			argsSourceCoin := string(args[4])
			argsTargetCoin := string(args[5])
			argsExchRateDenom := string(args[6])
			argsExchRate := string(args[7])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()
			srcPort := args[0]
			srcChannel := args[1]

			// Get the relative timeout timestamp
			timeoutTimestamp, err := cmd.Flags().GetUint64(flagPacketTimeoutTimestamp)
			if err != nil {
				return err
			}
			consensusState, _, _, err := channelutils.QueryLatestConsensusState(clientCtx, srcPort, srcChannel)
			if err != nil {
				return err
			}
			if timeoutTimestamp != 0 {
				timeoutTimestamp = consensusState.GetTimestamp() + timeoutTimestamp
			}

			msg := types.NewMsgSendCreateOrder(sender, srcPort, srcChannel, timeoutTimestamp, string(argsAmountDenom), int32(argsAmount), string(argsSourceCoin), string(argsTargetCoin), string(argsExchRateDenom), string(argsExchRate))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64(flagPacketTimeoutTimestamp, DefaultRelativePacketTimeoutTimestamp, "Packet timeout timestamp in nanoseconds. Default is 10 minutes.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
