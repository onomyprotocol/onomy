package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ParamSubspace defines the expected Subspace interface.
type ParamSubspace interface {
	HasKeyTable() bool
	WithKeyTable(paramtypes.KeyTable) paramtypes.Subspace
	Get(sdk.Context, []byte, interface{})
	GetParamSet(sdk.Context, paramtypes.ParamSet)
	SetParamSet(sdk.Context, paramtypes.ParamSet)
}

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAddress(string) sdk.AccAddress
}

// BankKeeper defines the contract needed to be fulfilled for banking and supply dependencies.
type BankKeeper interface {
	GetAllBalances(sdk.Context, sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(sdk.Context, sdk.AccAddress, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(sdk.Context, string, sdk.AccAddress, sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	MintCoins(sdk.Context, string, sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(sdk.Context, string, sdk.AccAddress, sdk.Coins) error
}

// DistributionKeeper expected distribution keeper.
type DistributionKeeper interface {
	WithdrawDelegationRewards(sdk.Context, sdk.AccAddress, sdk.ValAddress) (sdk.Coins, error)
}

// StakingKeeper expected staking keeper.
type StakingKeeper interface {
	BondDenom(sdk.Context) string
	Delegate(sdk.Context, sdk.AccAddress, sdk.Int, stakingtypes.BondStatus, stakingtypes.Validator, bool) (sdk.Dec, error)
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI
	GetAllValidators(sdk.Context) []stakingtypes.Validator
	GetAllDelegatorDelegations(sdk.Context, sdk.AccAddress) []stakingtypes.Delegation
	Unbond(sdk.Context, sdk.AccAddress, sdk.ValAddress, sdk.Dec) (sdk.Int, error)
}