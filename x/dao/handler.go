package dao

// import (
// 	"fmt"

// 	sdkerrors "cosmossdk.io/errors"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	// "github.com/onomyprotocol/onomy/x/dao/keeper"
// 	"github.com/onomyprotocol/onomy/x/dao/types"
// ).

// // NewHandler ...
// func NewHandler() sdk.Handler {
// 	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
// 		switch msg := msg.(type) { // nolint:gocritic //the module doesn't support messages handling
// 		default:
// 			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
// 			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
// 		}
// 	}
// }.
