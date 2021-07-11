package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendCreatePair{}

func NewMsgSendCreatePair(
	sender string,
	port string,
	channelID string,
	timeoutTimestamp uint64,
	sourceDenom string,
	targetDenom string,
) *MsgSendCreatePair {
	return &MsgSendCreatePair{
		Sender:           sender,
		Port:             port,
		ChannelID:        channelID,
		TimeoutTimestamp: timeoutTimestamp,
		SourceDenom:      sourceDenom,
		TargetDenom:      targetDenom,
	}
}

func (msg *MsgSendCreatePair) Route() string {
	return RouterKey
}

func (msg *MsgSendCreatePair) Type() string {
	return "SendCreatePair"
}

func (msg *MsgSendCreatePair) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgSendCreatePair) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendCreatePair) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
