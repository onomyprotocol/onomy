package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// RegisterCodec registers the legacy amino codec.
func RegisterCodec(cdc *codec.LegacyAmino) {
}

// RegisterInterfaces registers the cdctypes interface.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
}

var (
	// Amino holds the LegacyAmino codec.
	Amino = codec.NewLegacyAmino() // nolint:gochecknoglobals // cosmos sdk style
	// ModuleCdc holds the default proto codec.
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry()) // nolint:gochecknoglobals // cosmos sdk style
)
