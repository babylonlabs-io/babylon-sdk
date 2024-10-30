package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// MintBlockRewards mints new tokens and sends them to the staking contract for distribution.
// Authorization of the actor and recipient should be handled before entering this method.
// The amount is computed based on the inflation rate, the blocks per year, and the total staking token supply,
// in the bonded denom
func (k Keeper) MintBlockRewards(pCtx sdk.Context, recipient sdk.AccAddress, inflationRate sdkmath.LegacyDec, blocksPerYear int64) (sdkmath.Int, error) {
	bondDenom, err := k.Staking.BondDenom(pCtx)
	totalSupply, err := k.Staking.StakingTokenSupply(pCtx)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	blockRewards := inflationRate.MulInt(totalSupply).QuoInt64(blocksPerYear).MulInt64(1e6).TruncateInt()
	amt := sdk.NewCoin(bondDenom, blockRewards)

	if amt.Amount.IsNil() || amt.Amount.IsZero() || amt.Amount.IsNegative() {
		return sdkmath.ZeroInt(), errors.ErrInvalidRequest.Wrap("amount")
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
