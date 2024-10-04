package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	// ProposalTypeFundTreasuryProposal defines the type for a FundTreasuryProposal.
	ProposalTypeFundTreasuryProposal = "FundTreasuryProposal"
	// ProposalTypeExchangeWithTreasuryProposal defines the type for a ExchangeWithTreasuryProposal.
	ProposalTypeExchangeWithTreasuryProposal = "ExchangeWithTreasuryProposal"
	// ProposalTypeFundAccountProposal defines the type for a FundAccountProposal.
	ProposalTypeFundAccountProposal = "FundAccountProposal"
)

var (
	_ govtypes.Content = &FundTreasuryProposal{}
	_ govtypes.Content = &ExchangeWithTreasuryProposal{}
	_ govtypes.Content = &FundAccountProposal{}
)

func init() { //nolint:gochecknoinits // cosmos sdk style
	govtypes.RegisterProposalType(ProposalTypeFundTreasuryProposal)
	govtypes.RegisterProposalType(ProposalTypeExchangeWithTreasuryProposal)
	govtypes.RegisterProposalType(ProposalTypeFundAccountProposal)
}

// NewFundTreasuryProposal creates a new fund treasury proposal.
func NewFundTreasuryProposal(sender sdk.AccAddress, title, description string, amount sdk.Coins) *FundTreasuryProposal {
	return &FundTreasuryProposal{sender.String(), title, description, amount}
}

// GetTitle returns the title of a fund treasury proposal.
func (m *FundTreasuryProposal) GetTitle() string { return m.Title }

// GetDescription returns the description of a fund treasury proposal.
func (m *FundTreasuryProposal) GetDescription() string { return m.Description }

// ProposalRoute returns the routing key of a fund treasury proposal.
func (m *FundTreasuryProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the fund treasury proposal.
func (m *FundTreasuryProposal) ProposalType() string { return ProposalTypeFundTreasuryProposal }

// ValidateBasic runs basic stateless validity checks.
func (m *FundTreasuryProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(sender); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if !m.Amount.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsAllPositive() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// GetProposer returns the proposer from the proposal struct.
func (m *FundTreasuryProposal) GetProposer() string { return m.Sender }

// String implements the Stringer interface.
func (m FundTreasuryProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Fund treasury proposal:
  Sender: %s
  Title: %s
  Description: %s
  Amount: %s
`, m.Sender, m.Title, m.Description, m.Amount))
	return b.String()
}

// NewExchangeWithTreasuryProposal creates a new fund treasury proposal.
func NewExchangeWithTreasuryProposal(sender sdk.AccAddress, title, description string, coinsPairs []CoinsExchangePair) *ExchangeWithTreasuryProposal {
	return &ExchangeWithTreasuryProposal{sender.String(), title, description, coinsPairs}
}

// GetTitle returns the title of a fund treasury proposal.
func (m *ExchangeWithTreasuryProposal) GetTitle() string { return m.Title }

// GetDescription returns the description of a fund treasury proposal.
func (m *ExchangeWithTreasuryProposal) GetDescription() string { return m.Description }

// ProposalRoute returns the routing key of a fund treasury proposal.
func (m *ExchangeWithTreasuryProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the fund treasury proposal.
func (m *ExchangeWithTreasuryProposal) ProposalType() string {
	return ProposalTypeExchangeWithTreasuryProposal
}

// ValidateBasic runs basic stateless validity checks.
func (m *ExchangeWithTreasuryProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(sender); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if len(m.CoinsPairs) == 0 {
		return errors.Wrapf(ErrInvalidCoinsPair, "coins pairs can't be empty")
	}

	for i := range m.CoinsPairs {
		if err := m.CoinsPairs[i].ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// GetProposer returns the proposer from the proposal struct.
func (m *ExchangeWithTreasuryProposal) GetProposer() string { return m.Sender }

// String implements the Stringer interface.
func (m ExchangeWithTreasuryProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Exchange with treasury proposal:
  Sender: %s
  Title: %s
  Description: %s
  Coin Pairs: %s
`, m.Sender, m.Title, m.Description, m.CoinsPairs))
	return b.String()
}

// ValidateBasic validates CoinsExchangePair basic options.
func (m *CoinsExchangePair) ValidateBasic() error {
	if m == nil {
		return errors.Wrapf(ErrInvalidCoinsPair, "coins pairs can't be nil")
	}

	if !m.CoinAsk.IsValid() || m.CoinAsk.IsZero() {
		return errors.Wrapf(ErrInvalidCoinsPair, "invalid coin ask")
	}

	if !m.CoinBid.IsValid() || m.CoinBid.IsZero() {
		return errors.Wrapf(ErrInvalidCoinsPair, "invalid coin bid")
	}

	return nil
}

// NewFundAccountProposal creates a new fund account proposal.
func NewFundAccountProposal(recipient sdk.AccAddress, title, description string, amount sdk.Coins) *FundAccountProposal {
	return &FundAccountProposal{recipient.String(), title, description, amount}
}

// GetTitle returns the title of a fund account proposal.
func (m *FundAccountProposal) GetTitle() string { return m.Title }

// GetDescription returns the description of a fund account proposal.
func (m *FundAccountProposal) GetDescription() string { return m.Description }

// ProposalRoute returns the routing key of a fund account proposal.
func (m *FundAccountProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the fund account proposal.
func (m *FundAccountProposal) ProposalType() string { return ProposalTypeFundAccountProposal }

// ValidateBasic runs basic stateless validity checks.
func (m *FundAccountProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}
	sender, err := sdk.AccAddressFromBech32(m.Recipient)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(sender); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if !m.Amount.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsAllPositive() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// String implements the Stringer interface.
func (m FundAccountProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Fund account proposal:
  Recipient: %s
  Title: %s
  Description: %s
  Amount: %s
`, m.Recipient, m.Title, m.Description, m.Amount))
	return b.String()
}
