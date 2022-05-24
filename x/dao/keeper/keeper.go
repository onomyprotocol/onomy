// Package keeper contains dao module keeper.
package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

type (
	// Keeper is a dao keeper struct.
	Keeper struct {
		cdc                codec.BinaryCodec
		storeKey           sdk.StoreKey
		memKey             sdk.StoreKey
		ps                 types.ParamSubspace
		bankKeeper         types.BankKeeper
		accountKeeper      types.AccountKeeper
		distributionKeeper types.DistributionKeeper
		govKeeper          types.GovKeeper
		stakingKeeper      types.StakingKeeper
	}
)

// NewKeeper creates new dao keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps types.ParamSubspace,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	distributionKeeper types.DistributionKeeper,
	govKeeper types.GovKeeper,
	stakingKeeper types.StakingKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// ensure dao module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		ps:                 ps,
		bankKeeper:         bankKeeper,
		accountKeeper:      accountKeeper,
		distributionKeeper: distributionKeeper,
		govKeeper:          govKeeper,
		stakingKeeper:      stakingKeeper,
	}
}

// Logger returns keeper logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
