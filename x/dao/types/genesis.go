package types

// DefaultGenesis returns the default dao genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (m GenesisState) Validate() error {
	return m.TreasuryBalance.Validate()
}
