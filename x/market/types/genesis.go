package types

import (
	"fmt"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,
		// this line is used by starport scaffolding # genesis/types/default
		DenomTraceList: []*DenomTrace{},
		OrderBookList:  []*OrderBook{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated index in denomTrace
	denomTraceIndexMap := make(map[string]bool)

	for _, elem := range gs.DenomTraceList {
		if _, ok := denomTraceIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for denomTrace")
		}
		denomTraceIndexMap[elem.Index] = true
	}
	// Check for duplicated index in orderBook
	orderBookIndexMap := make(map[string]bool)

	for _, elem := range gs.OrderBookList {
		if _, ok := orderBookIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for orderBook")
		}
		orderBookIndexMap[elem.Index] = true
	}

	return nil
}
