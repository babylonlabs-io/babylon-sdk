package keeper

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetParams sets the x/incentive module parameters.
func (k Keeper) SetParams(ctx context.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&p)
	return store.Set(types.ParamsKey, bz)
}

// GetParams returns the current x/incentive module parameters.
func (k Keeper) GetParams(ctx context.Context) (p types.Params) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return p
	}
	k.cdc.MustUnmarshal(bz, &p)
	return p
}

func (k Keeper) GetMaxSudoGas(ctx sdk.Context) storetypes.Gas {
	return storetypes.Gas(k.GetParams(ctx).MaxGasBeginBlocker)
}
