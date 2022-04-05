package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// ErrInsufficientBalance - the user balance is insufficient for the operation.
var ErrInsufficientBalance = sdkerrors.Register(ModuleName, 1, "insufficient balance")
