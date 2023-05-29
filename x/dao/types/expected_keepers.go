package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
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
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(sdk.Context, sdk.AccAddress, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(sdk.Context, string, sdk.AccAddress, sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	MintCoins(sdk.Context, string, sdk.Coins) error
}

// DistributionKeeper expected distribution keeper.
type DistributionKeeper interface {
	HasDelegatorStartingInfo(sdk.Context, sdk.ValAddress, sdk.AccAddress) bool
	WithdrawDelegationRewards(sdk.Context, sdk.AccAddress, sdk.ValAddress) (sdk.Coins, error)
}

// GovKeeper expected gov keeper.
type GovKeeper interface {
	AddVote(sdk.Context, uint64, sdk.AccAddress, govtypes.WeightedVoteOptions) error
	GetVote(sdk.Context, uint64, sdk.AccAddress) (govtypes.Vote, bool)
	IterateProposals(sdk.Context, func(proposal govtypes.Proposal) bool)
}

// MintKeeper expected mint keeper.
type MintKeeper interface {
	GetMinter(ctx sdk.Context) (minter minttypes.Minter)
	GetParams(ctx sdk.Context) (params minttypes.Params)
}

// StakingKeeper expected staking keeper.
type StakingKeeper interface {
	BondDenom(sdk.Context) string
	Delegate(sdk.Context, sdk.AccAddress, sdk.Int, stakingtypes.BondStatus, stakingtypes.Validator, bool) (sdk.Dec, error)
	GetDelegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) (stakingtypes.Delegation, bool)
	GetAllValidators(sdk.Context) []stakingtypes.Validator
	UnbondAndUndelegateCoins(sdk.Context, sdk.AccAddress, sdk.ValAddress, sdk.Dec) (sdk.Int, error)
}
