// Package types is onomy cosmos sdk types.
package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registries the codec in the registry.
func RegisterCodec(cdc *codec.LegacyAmino) { // nolint:staticcheck
	// this line is used by starport scaffolding # 2
}

// RegisterInterfaces registries the codec interface in the registry.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
