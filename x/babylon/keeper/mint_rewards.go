package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// MintBlockRewards mints new "virtual" bonding tokens and sends them to the staking contract for distribution.
// The amount minted is removed from the SupplyOffset (so that it will become negative), when supported.
// Authorization of the actor/recipient should be handled before entering this method.
// The amount is computed based on the finality inflation rate and the total staking token supply, in the bonded denom
func (k Keeper) MintBlockRewards(pCtx sdk.Context, recipient sdk.AccAddress, amt sdk.Coin) (sdkmath.Int, error) {
	if amt.Amount.IsNil() || amt.Amount.IsZero() || amt.Amount.IsNegative() {
		return sdkmath.ZeroInt(), errors.ErrInvalidRequest.Wrap("amount")
	}

	// Ensure staking constraints
	bondDenom, err := k.Staking.BondDenom(pCtx)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	if amt.Denom != bondDenom {
		return sdkmath.ZeroInt(), errors.ErrInvalidRequest.Wrapf("invalid coin denomination: got %s, expected %s", amt.Denom, bondDenom)
	}
	// FIXME? Remove this constraint for flexibility
	params := k.GetParams(pCtx)
	if recipient.String() != params.BtcStakingContractAddress {
		return sdkmath.ZeroInt(), errors.ErrUnauthorized.Wrapf("invalid recipient: got %s, expected staking contract (%s)",
			recipient, params.BtcStakingContractAddress)
	}

	// TODO?: Ensure Babylon constraints

	cacheCtx, done := pCtx.CacheContext() // work in a cached store as Osmosis (safety net?)

	// Mint rewards tokens
	coins := sdk.NewCoins(amt)
	err = k.bank.MintCoins(cacheCtx, types.ModuleName, coins)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	// FIXME: Confirm we want this supply offset enabled for rewards, i.e.
	// as virtual coins that do not count to the total supply
	//k.bank.AddSupplyOffset(cacheCtx, bondDenom, amt.Amount.Neg())

	err = k.bank.SendCoinsFromModuleToAccount(cacheCtx, types.ModuleName, recipient, coins)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	done()
	return amt.Amount, err
}
