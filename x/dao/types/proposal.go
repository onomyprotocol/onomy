package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeFundTreasuryProposal defines the type for a FundTreasuryProposal.
	ProposalTypeFundTreasuryProposal = "FundTreasuryProposal"
	// ProposalTypeExchangeWithTreasuryProposal defines the type for a ExchangeWithTreasuryProposal.
	ProposalTypeExchangeWithTreasuryProposal = "ExchangeWithTreasuryProposal"
)

// Assert FundTreasuryProposal implements govtypes.Content at compile-time.
var (
	_ govtypes.Content = &FundTreasuryProposal{}
	_ govtypes.Content = &ExchangeWithTreasuryProposal{}
)

func init() { // nolint:gochecknoinits // cosmos sdk style
	govtypes.RegisterProposalType(ProposalTypeFundTreasuryProposal)
	govtypes.RegisterProposalTypeCodec(&FundTreasuryProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeFundTreasuryProposal))

	govtypes.RegisterProposalType(ProposalTypeExchangeWithTreasuryProposal)
	govtypes.RegisterProposalTypeCodec(&ExchangeWithTreasuryProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeExchangeWithTreasuryProposal))
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if !m.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// GetProposer returns the proposer from the proposal struct.
func (m *FundTreasuryProposal) GetProposer() string { return m.Sender }

// String implements the Stringer interface.
func (m FundTreasuryProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Fund treasury Proposal:
  Sender: %s
  Title: %s
  Description: %s
  Amount: %s
`, m.Sender, m.Title, m.Description, m.Amount))
	return b.String()
}

// NewExchangeWithTreasuryProposal creates a new fund treasury proposal.
func NewExchangeWithTreasuryProposal(sender sdk.AccAddress, title, description string, amountAsk, amountBid sdk.Coins) *ExchangeWithTreasuryProposal {
	return &ExchangeWithTreasuryProposal{sender.String(), title, description, amountAsk, amountBid}
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if !m.AmountAsk.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.AmountAsk.String())
	}

	if !m.AmountAsk.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.AmountAsk.String())
	}

	if !m.AmountBid.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.AmountBid.String())
	}

	if !m.AmountBid.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.AmountBid.String())
	}

	return nil
}

// GetProposer returns the proposer from the proposal struct.
func (m *ExchangeWithTreasuryProposal) GetProposer() string { return m.Sender }

// String implements the Stringer interface.
func (m ExchangeWithTreasuryProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Exchange with treasury Proposal:
  Sender: %s
  Title: %s
  Description: %s
  AmountAsk: %s
  AmountBid: %s	
`, m.Sender, m.Title, m.Description, m.AmountAsk, m.AmountBid))
	return b.String()
}
