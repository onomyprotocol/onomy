package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var (
	// DefaultWithdrawRewardPeriod is default value for the DefaultWithdrawRewardPeriod param.
	DefaultWithdrawRewardPeriod = int64(51840) //nolint:gomnd,gochecknoglobals // cosmos-sdk style
	// DefaultStakingTokenPoolRate is default value for the DefaultStakingTokenPoolRate param.
	DefaultStakingTokenPoolRate = sdk.NewDec(1).Quo(sdk.NewDec(20)) //nolint:gomnd,gochecknoglobals // cosmos-sdk style
	// DefaultStakingTokenMaxProposalRate is default value for the DefaultStakingTokenMaxProposalRate param.
	DefaultStakingTokenMaxProposalRate = sdk.NewDec(1).Quo(sdk.NewDec(20)) //nolint:gomnd,gochecknoglobals // cosmos-sdk style
	// DefaultStakingMaxCommissionRate is default value for the DefaultStakingMaxCommissionRate param.
	DefaultStakingMaxCommissionRate = sdk.NewDec(1).Quo(sdk.NewDec(5)) //nolint:gomnd,gochecknoglobals // cosmos-sdk style
)

// Parameter store keys.
var (
	// KeyWithdrawRewardPeriod is byte key for KeyWithdrawRewardPeriod param.
	KeyWithdrawRewardPeriod = []byte("WithdrawRewardPeriod") //nolint:gochecknoglobals // cosmos-sdk style
	// KeyStakingTokenPoolRate is byte key for KeyStakingTokenPoolRate param.
	KeyStakingTokenPoolRate = []byte("StakingTokenPoolRate") //nolint:gochecknoglobals // cosmos-sdk style
	// KeyStakingTokenMaxProposalRate is byte key for KeyStakingTokenMaxProposalRate param.
	KeyStakingTokenMaxProposalRate = []byte("StakingTokenMaxProposalRate") //nolint:gochecknoglobals // cosmos-sdk style
	// KeyStakingMaxCommissionRate is byte key for KeyStakingMaxCommissionRate param.
	KeyStakingMaxCommissionRate = []byte("StakingMaxCommissionRate") //nolint:gochecknoglobals // cosmos-sdk style
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance.
func NewParams(
	withdrawRewardPeriod int64,
	stakingTokenPoolRate,
	stakingTokenMaxProposalRate,
	stakingMaxCommissionRate sdk.Dec,
) Params {
	return Params{
		WithdrawRewardPeriod:        withdrawRewardPeriod,
		StakingTokenPoolRate:        stakingTokenPoolRate,
		StakingTokenMaxProposalRate: stakingTokenMaxProposalRate,
		StakingMaxCommissionRate:    stakingMaxCommissionRate,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultWithdrawRewardPeriod, DefaultStakingTokenPoolRate, DefaultStakingTokenMaxProposalRate,
		DefaultStakingMaxCommissionRate,
	)
}

// ParamSetPairs get the params.ParamSet.
func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyWithdrawRewardPeriod, &m.WithdrawRewardPeriod, validateWithdrawRewardPeriod),
		paramtypes.NewParamSetPair(KeyStakingTokenPoolRate, &m.StakingTokenPoolRate, validateStakingTokenPoolRate),
		paramtypes.NewParamSetPair(KeyStakingTokenMaxProposalRate, &m.StakingTokenMaxProposalRate, validateStakingTokenMaxProposalRate),
		paramtypes.NewParamSetPair(KeyStakingMaxCommissionRate, &m.StakingMaxCommissionRate, validateStakingMaxCommissionRate),
	}
}

// Validate validates the set of params.
func (m Params) Validate() error {
	if err := validateWithdrawRewardPeriod(m.WithdrawRewardPeriod); err != nil {
		return err
	}
	if err := validateStakingTokenPoolRate(m.StakingTokenPoolRate); err != nil {
		return err
	}
	if err := validateStakingTokenMaxProposalRate(m.StakingTokenMaxProposalRate); err != nil {
		return err
	}

	return validateStakingMaxCommissionRate(m.StakingMaxCommissionRate)
}

// String implements the Stringer interface.
func (m Params) String() string {
	out, _ := yaml.Marshal(m) //nolint:errcheck // error is not expected here
	return string(out)
}

func validateWithdrawRewardPeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("withdraw reward period must be positive: %d", v)
	}

	return nil
}

func validateStakingTokenPoolRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("staking token pool rate cannot be negative or nil: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("staking token pool rate too large: %s", v)
	}

	return nil
}

func validateStakingTokenMaxProposalRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("staking token max proposal rate cannot be negative or nil: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("staking token max proposal rate too large: %s", v)
	}

	return nil
}

func validateStakingMaxCommissionRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("staking max commission rate cannot be negative or nil: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("staking max commission rate too large: %s", v)
	}

	return nil
}
