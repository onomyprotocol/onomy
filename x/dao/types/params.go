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
	// DefaultPoolRate is default value for the DefaultPoolRate param.
	DefaultPoolRate = sdk.NewDec(1).Quo(sdk.NewDec(20)) //nolint:gomnd,gochecknoglobals // cosmos-sdk style
	// DefaultMaxProposalRate is default value for the DefaultMaxProposalRate param.
	DefaultMaxProposalRate = sdk.NewDec(1).Quo(sdk.NewDec(20)) //nolint:gomnd,gochecknoglobals // cosmos-sdk style
	// DefaultMaxValCommission is default value for the DefaultMaxValCommission param.
	DefaultMaxValCommission = sdk.NewDec(1).Quo(sdk.NewDec(5)) //nolint:gomnd,gochecknoglobals // cosmos-sdk style
)

// Parameter store keys.
var (
	// KeyWithdrawRewardPeriod is byte key for KeyWithdrawRewardPeriod param.
	KeyWithdrawRewardPeriod = []byte("WithdrawRewardPeriod") //nolint:gochecknoglobals // cosmos-sdk style
	// KeyPoolRate is byte key for KeyPoolRate param.
	KeyPoolRate = []byte("PoolRate") //nolint:gochecknoglobals // cosmos-sdk style
	// KeyMaxProposalRate is byte key for KeyMaxProposalRate param.
	KeyMaxProposalRate = []byte("MaxProposalRate") //nolint:gochecknoglobals // cosmos-sdk style
	// KeyMaxValCommission is byte key for KeyMaxValCommission param.
	KeyMaxValCommission = []byte("MaxValCommission") //nolint:gochecknoglobals // cosmos-sdk style
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance.
func NewParams(
	withdrawRewardPeriod int64,
	poolRate,
	maxProposalRate,
	maxValCommission sdk.Dec,
) Params {
	return Params{
		WithdrawRewardPeriod: withdrawRewardPeriod,
		PoolRate:             poolRate,
		MaxProposalRate:      maxProposalRate,
		MaxValCommission:     maxValCommission,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultWithdrawRewardPeriod, DefaultPoolRate, DefaultMaxProposalRate,
		DefaultMaxValCommission,
	)
}

// ParamSetPairs get the params.ParamSet.
func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyWithdrawRewardPeriod, &m.WithdrawRewardPeriod, validateWithdrawRewardPeriod),
		paramtypes.NewParamSetPair(KeyPoolRate, &m.PoolRate, validatePoolRate),
		paramtypes.NewParamSetPair(KeyMaxProposalRate, &m.MaxProposalRate, validateMaxProposalRate),
		paramtypes.NewParamSetPair(KeyMaxValCommission, &m.MaxValCommission, validateMaxValCommission),
	}
}

// Validate validates the set of params.
func (m Params) Validate() error {
	if err := validateWithdrawRewardPeriod(m.WithdrawRewardPeriod); err != nil {
		return err
	}
	if err := validatePoolRate(m.PoolRate); err != nil {
		return err
	}
	if err := validateMaxProposalRate(m.MaxProposalRate); err != nil {
		return err
	}

	return validateMaxValCommission(m.MaxValCommission)
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

func validatePoolRate(i interface{}) error {
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

func validateMaxProposalRate(i interface{}) error {
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

func validateMaxValCommission(i interface{}) error {
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
