package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	// ErrInsufficientBalance - the balance is insufficient for the operation.
	ErrInsufficientBalance = sdkerrors.Register(ModuleName, 1, "insufficient balance")
	// ErrInvalidCoinsPair - the coins pair is invalid.
	ErrInvalidCoinsPair = sdkerrors.Register(ModuleName, 2, "invalid coins pair") // nolint:gomnd // error number
)
