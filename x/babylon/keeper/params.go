package keeper

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetParams sets the module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetParams gets the module's parameters.
func (k Keeper) GetParams(clientCtx sdk.Context) (params types.Params) {
	store := clientCtx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) GetFinalityInflationRate(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).FinalityInflationRate
}

func (k Keeper) GetExpectedBlocksPerYear(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).BlocksPerYear
}

func (k Keeper) GetMaxSudoGas(ctx sdk.Context) storetypes.Gas {
	return storetypes.Gas(k.GetParams(ctx).MaxGasBeginBlocker)
}
