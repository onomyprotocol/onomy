package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(params Params, treasuryBalance sdk.Coins) *GenesisState {
	return &GenesisState{
		Params:          params,
		TreasuryBalance: treasuryBalance,
	}
}

// DefaultGenesis returns the default dao genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (m GenesisState) Validate() error {
	if err := m.Params.Validate(); err != nil {
		return err
	}
	return m.TreasuryBalance.Validate()
}
