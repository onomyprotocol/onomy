package types

import sdkerrors "cosmossdk.io/errors"

var (
	// ErrInsufficientBalance - the balance is insufficient for the operation.
	ErrInsufficientBalance = sdkerrors.Register(ModuleName, 1, "insufficient balance")
	// ErrInvalidCoinsPair - the coins pair is invalid.
	ErrInvalidCoinsPair = sdkerrors.Register(ModuleName, 2, "invalid coins pair") //nolint:gomnd // error number
	// ErrProhibitedCoinsAmount - the requested amount is prohibited.
	ErrProhibitedCoinsAmount = sdkerrors.Register(ModuleName, 3, "prohibited coins amount") //nolint:gomnd // error number
)
