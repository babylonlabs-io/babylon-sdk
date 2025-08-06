package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

// SetParams sets the module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
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

func (k Keeper) GetMaxSudoGasBeginBlocker(ctx sdk.Context) storetypes.Gas {
	return storetypes.Gas(k.GetParams(ctx).MaxGasBeginBlocker)
}

func (k Keeper) GetMaxSudoGasEndBlocker(ctx sdk.Context) storetypes.Gas {
	return storetypes.Gas(k.GetParams(ctx).MaxGasEndBlocker)
}
