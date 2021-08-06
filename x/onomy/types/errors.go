package types

// DONTCOVER.

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	errSampleCode = 1100
)

// x/onomy module sentinel errors.
var (
	ErrSample = sdkerrors.Register(ModuleName, errSampleCode, "sample error")
	// this line is used by starport scaffolding # ibc/errors.
)
