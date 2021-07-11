package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/onomyprotocol/onomy/x/market/types"
)

func (k msgServer) SendCreatePair(goCtx context.Context, msg *types.MsgSendCreatePair) (*types.MsgSendCreatePairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pairIndex := types.OrderBookIndex(msg.Port, msg.ChannelID, msg.SourceDenom, msg.TargetDenom)
	_, found := k.GetOrderBook(ctx, pairIndex)
	if found {
		return &types.MsgSendCreatePairResponse{}, errors.New("the pair already exist")
	}

	// Construct the packet
	var packet types.CreatePairPacketData

	packet.SourceDenom = msg.SourceDenom
	packet.TargetDenom = msg.TargetDenom

	// Transmit the packet
	err := k.TransmitCreatePairPacket(
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

	return &types.MsgSendCreatePairResponse{}, nil
}
