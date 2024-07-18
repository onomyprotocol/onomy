// Package v1_1_5 is contains chain upgrade of the corresponding version.
package v1_1_5 //nolint:revive,stylecheck // app version

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	daotypes "github.com/onomyprotocol/onomy/x/dao/types"
)

// Name is migration name.
const Name = "v1.1.5"

var burntAddrs = []string{
	"onomy1302z68lkj2qa4c9qxh387pr973k8ulvj6880vu",
	"onomy1ttm9jz0zlr9gtwf3j33a57jsk6lx0yeejzy6ek",
	"onomy1x6mlfjektdjgum6hfgfgcz57dl4s9dss59tmjy",
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak *authkeeper.AccountKeeper,
	bk *bankkeeper.BaseKeeper,
	sk *stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		for _, addr := range burntAddrs {
			// get target account
			account := ak.GetAccount(ctx, sdk.MustAccAddressFromBech32(addr))

			// check that it's a vesting account type
			continousVestAccount, ok := account.(*vesting.ContinuousVestingAccount)
			if ok {
				// overwrite vest account to a normal base account
				ak.SetAccount(ctx, continousVestAccount.BaseAccount)
			}

			periodicVestAccount, ok := account.(*vesting.PeriodicVestingAccount)
			if ok {
				// overwrite vest account to a normal base account
				ak.SetAccount(ctx, periodicVestAccount.BaseAccount)
			}

			delayedVestAccount, ok := account.(*vesting.DelayedVestingAccount)
			if ok {
				// overwrite vest account to a normal base account
				ak.SetAccount(ctx, delayedVestAccount.BaseAccount)
			}

			permanentLockedAccount, ok := account.(*vesting.PermanentLockedAccount)
			if ok {
				// overwrite vest account to a normal base account
				ak.SetAccount(ctx, permanentLockedAccount.BaseAccount)
			}

			// unbond all delegations from account
			err := forceUnbondTokens(ctx, addr, bk, sk)
			if err != nil {
				ctx.Logger().Error("Error force unbonding delegations")
				return nil, err
			}

			// finish all current unbonding entries
			err = forceFinishUnbonding(ctx, addr, bk, sk)
			if err != nil {
				ctx.Logger().Error("Error force finishing unbonding delegations")
				return nil, err
			}

			// send to dao module account
			// vesting account should be able to send coins normaly after
			// we converted it back to a base account
			bal := bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(addr))
			err = bk.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(addr), daotypes.ModuleName, bal)
			if err != nil {
				ctx.Logger().Error("Error reallocating funds")
				return nil, err
			}
		}
		ctx.Logger().Info("Finished reallocating funds")

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

func forceFinishUnbonding(ctx sdk.Context, delAddr string, bk *bankkeeper.BaseKeeper, sk *stakingkeeper.Keeper) error {
	ubdQueue := sk.GetAllUnbondingDelegations(ctx, sdk.MustAccAddressFromBech32(delAddr))
	bondDenom := sk.BondDenom(ctx)

	for _, ubd := range ubdQueue {
		for _, entry := range ubd.Entries {
			err := bk.UndelegateCoinsFromModuleToAccount(ctx, stakingtypes.NotBondedPoolName, sdk.MustAccAddressFromBech32(delAddr), sdk.NewCoins(sdk.NewCoin(bondDenom, entry.Balance)))
			if err != nil {
				return err
			}
		}

		// empty out all entries
		ubd.Entries = []stakingtypes.UnbondingDelegationEntry{}
		sk.SetUnbondingDelegation(ctx, ubd)
	}

	return nil
}

func forceUnbondTokens(ctx sdk.Context, delAddr string, bk *bankkeeper.BaseKeeper, sk *stakingkeeper.Keeper) error {
	delAccAddr := sdk.MustAccAddressFromBech32(delAddr)
	dels := sk.GetDelegatorDelegations(ctx, delAccAddr, 100)

	for _, del := range dels {
		valAddr := del.GetValidatorAddr()

		validator, found := sk.GetValidator(ctx, valAddr)
		if !found {
			return stakingtypes.ErrNoValidatorFound
		}

		returnAmount, err := sk.Unbond(ctx, delAccAddr, valAddr, del.GetShares())
		if err != nil {
			return err
		}

		coins := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), returnAmount))

		// transfer the validator tokens to the not bonded pool
		if validator.IsBonded() {
			// doing stakingKeeper.bondedTokensToNotBonded
			err = bk.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, coins)
			if err != nil {
				return err
			}
		}

		err = bk.UndelegateCoinsFromModuleToAccount(ctx, stakingtypes.NotBondedPoolName, delAccAddr, coins)
		if err != nil {
			return err
		}
	}

	return nil
}
