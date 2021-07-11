package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/onomyprotocol/onomy/x/market/types"
)

func (k msgServer) SendCreateOrder(goCtx context.Context, msg *types.MsgSendCreateOrder) (*types.MsgSendCreateOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: logic before transmitting the packet

	// Construct the packet
	var packet types.CreateOrderPacketData

	packet.AmountDenom = msg.AmountDenom
	packet.Amount = msg.Amount
	packet.SourceCoin = msg.SourceCoin
	packet.TargetCoin = msg.TargetCoin
	packet.ExchRateDenom = msg.ExchRateDenom
	packet.ExchRate = msg.ExchRate

	// Transmit the packet
	err := k.TransmitCreateOrderPacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendCreateOrderResponse{}, nil
}
