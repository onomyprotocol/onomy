package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCancelOrder{}

func NewMsgCancelOrder(creator string, port string, channel string, amountDenom string, exchRateDenom string, orderID int32) *MsgCancelOrder {
	return &MsgCancelOrder{
		Creator:       creator,
		Port:          port,
		Channel:       channel,
		AmountDenom:   amountDenom,
		ExchRateDenom: exchRateDenom,
		OrderID:       orderID,
	}
}

func (msg *MsgCancelOrder) Route() string {
	return RouterKey
}

func (msg *MsgCancelOrder) Type() string {
	return "CancelOrder"
}

func (msg *MsgCancelOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
