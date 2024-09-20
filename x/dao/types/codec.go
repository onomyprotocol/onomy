package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterCodec registers the legacy amino codec.
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&FundTreasuryProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeFundTreasuryProposal), nil)
	cdc.RegisterConcrete(&ExchangeWithTreasuryProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeExchangeWithTreasuryProposal), nil)
	cdc.RegisterConcrete(&FundAccountProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeFundAccountProposal), nil)
}

// RegisterInterfaces registers the cdctypes interface.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&FundTreasuryProposal{},
		&ExchangeWithTreasuryProposal{},
		&FundAccountProposal{},
	)
}

var (
	// Amino holds the LegacyAmino codec.
	Amino = codec.NewLegacyAmino() //nolint:gochecknoglobals // cosmos sdk style
	// ModuleCdc holds the default proto codec.
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry()) //nolint:gochecknoglobals // cosmos sdk style
)
