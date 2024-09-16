package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
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
	GetAllBalances(context.Context, sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(context.Context, sdk.AccAddress, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(context.Context, string, sdk.AccAddress, sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderPool, recipientPool string, amt sdk.Coins) error
	MintCoins(context.Context, string, sdk.Coins) error
}

// DistributionKeeper expected distribution keeper.
type DistributionKeeper interface {
	HasDelegatorStartingInfo(context.Context, sdk.ValAddress, sdk.AccAddress) bool
	WithdrawDelegationRewards(context.Context, sdk.AccAddress, sdk.ValAddress) (sdk.Coins, error)
}

// GovKeeper expected gov keeper.
type GovKeeper interface {
	AddVote(context.Context, uint64, sdk.AccAddress, govtypes.WeightedVoteOptions) error
	GetVote(context.Context, uint64, sdk.AccAddress) (govtypes.Vote, bool)
	IterateProposals(context.Context, func(proposal govtypes.Proposal) bool)
}

// MintKeeper expected mint keeper.
type MintKeeper interface {
	GetMinter(ctx context.Context) (minter minttypes.Minter)
	GetParams(ctx context.Context) (params minttypes.Params)
}

// StakingKeeper expected staking keeper.
type StakingKeeper interface {
	BondDenom(context.Context) string
	Delegate(context.Context, sdk.AccAddress, math.Int, stakingtypes.BondStatus, stakingtypes.Validator, bool) (math.LegacyDec, error)
	GetDelegation(context.Context, sdk.AccAddress, sdk.ValAddress) (stakingtypes.Delegation, bool)
	GetAllValidators(context.Context) []stakingtypes.Validator
	UnbondAndUndelegateCoins(context.Context, sdk.AccAddress, sdk.ValAddress, math.LegacyDec) (math.Int, error)
}
