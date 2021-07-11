package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/onomyprotocol/onomy/x/market/types"
)

// TransmitCreatePairPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitCreatePairPacket(
	ctx sdk.Context,
	packetData types.CreatePairPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) error {

	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(
			channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", sourcePort, sourceChannel,
		)
	}

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: "+err.Error())
	}

	packet := channeltypes.NewPacket(
		packetBytes,
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp,
	)

	if err := k.channelKeeper.SendPacket(ctx, channelCap, packet); err != nil {
		return err
	}

	return nil
}

// OnRecvCreatePairPacket processes packet reception
func (k Keeper) OnRecvCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData) (packetAck types.CreatePairPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// Check if the buy order book exists
	pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.SourceDenom, data.TargetDenom)
	_, found := k.GetOrderBook(ctx, pairIndex)
	if found {
		return packetAck, errors.New("the pair already exist")
	}

	// Set the buy order book
	book := types.NewOrderBook(data.SourceDenom, data.TargetDenom)
	book.Index = pairIndex
	k.SetOrderBook(ctx, book)

	return packetAck, nil
}

// OnAcknowledgementCreatePairPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.CreatePairPacketAck
		err := packetAck.Unmarshal(dispatchedAck.Result)
		if err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.SourceDenom, data.TargetDenom)
		book := types.NewOrderBook(data.SourceDenom, data.TargetDenom)
		book.Index = pairIndex
		k.SetOrderBook(ctx, book)

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutCreatePairPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData) error {

	// TODO: packet timeout logic

	return nil
}
