package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendCreateOrder{}

func NewMsgSendCreateOrder(
	sender string,
	port string,
	channelID string,
	timeoutTimestamp uint64,
	amountDenom string,
	amount int32,
	sourceCoin string,
	targetCoin string,
	exchRateDenom string,
	exchRate string,
) *MsgSendCreateOrder {
	return &MsgSendCreateOrder{
		Sender:           sender,
		Port:             port,
		ChannelID:        channelID,
		TimeoutTimestamp: timeoutTimestamp,
		AmountDenom:      amountDenom,
		Amount:           amount,
		SourceCoin:       sourceCoin,
		TargetCoin:       targetCoin,
		ExchRateDenom:    exchRateDenom,
		ExchRate:         exchRate,
	}
}

func (msg *MsgSendCreateOrder) Route() string {
	return RouterKey
}

func (msg *MsgSendCreateOrder) Type() string {
	return "SendCreateOrder"
}

func (msg *MsgSendCreateOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgSendCreateOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendCreateOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
