package keeper

import (
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetBSNContracts stores the BSNContracts object in a single storage key
func (k Keeper) SetBSNContracts(ctx sdk.Context, contracts *types.BSNContracts) error {
	if err := contracts.ValidateBasic(); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(contracts)
	if err != nil {
		return err
	}
	store.Set(types.BSNContractsKey, bz)
	return nil
}

// GetBSNContracts retrieves the BSNContracts object from storage
func (k Keeper) GetBSNContracts(ctx sdk.Context) *types.BSNContracts {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BSNContractsKey)
	if bz == nil {
		return nil
	}
	var contracts types.BSNContracts
	k.cdc.MustUnmarshal(bz, &contracts)
	return &contracts
}
